package bootstrap

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	walletv1 "github.com/adamaso/wallet-service/gen/proto/v1"
	"github.com/adamaso/wallet-service/internal/application"
	"github.com/adamaso/wallet-service/internal/config"
	grpcserver "github.com/adamaso/wallet-service/internal/grpc"
	"github.com/adamaso/wallet-service/internal/infrastructure/postgres"
	"github.com/adamaso/wallet-service/internal/logger"
)

type Service struct {
	grpcServer   *grpc.Server
	healthServer *health.Server
	pool         *pgxpool.Pool
	grpcAddr     string
	log          *slog.Logger
}

func New(ctx context.Context, cfg config.Config, env string) (*Service, error) {
	log := logger.New(env)

	pool, err := newPool(ctx, cfg)
	if err != nil {
		return nil, err
	}

	repo := postgres.NewWalletRepository(pool)
	store := postgres.NewWalletViewStore(pool)

	app := application.NewApplication(repo, store)
	srv := grpcserver.NewServer(app)

	grpcSrv := grpc.NewServer(
		grpc.UnaryInterceptor(grpcserver.UnaryLoggingInterceptor(log)),
	)
	walletv1.RegisterWalletServiceServer(grpcSrv, srv)

	healthSrv := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcSrv, healthSrv)
	healthSrv.SetServingStatus(walletv1.WalletService_ServiceDesc.ServiceName, grpc_health_v1.HealthCheckResponse_SERVING)

	reflection.Register(grpcSrv)

	return &Service{
		grpcServer:   grpcSrv,
		healthServer: healthSrv,
		pool:         pool,
		grpcAddr:     cfg.GRPCAddr(),
		log:          log,
	}, nil
}

func (s *Service) Run(ctx context.Context) error {
	lis, err := net.Listen("tcp", s.grpcAddr)
	if err != nil {
		return fmt.Errorf("listen %s: %w", s.grpcAddr, err)
	}

	errCh := make(chan error, 1)
	go func() {
		if err := s.grpcServer.Serve(lis); err != nil {
			errCh <- err
		}
	}()

	s.log.Info("gRPC server started", slog.String("addr", s.grpcAddr))

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		s.Shutdown(shutdownCtx)
		return nil
	case err := <-errCh:
		return err
	}
}

func (s *Service) Shutdown(ctx context.Context) {
	s.log.Info("shutting down")
	s.healthServer.Shutdown()

	stopped := make(chan struct{})
	go func() {
		s.grpcServer.GracefulStop()
		close(stopped)
	}()

	select {
	case <-ctx.Done():
		s.grpcServer.Stop() // deadline exceeded — force stop
	case <-stopped:
	}

	s.pool.Close()
	s.log.Info("shutdown complete")
}

func newPool(ctx context.Context, cfg config.Config) (*pgxpool.Pool, error) {
	poolCfg, err := pgxpool.ParseConfig(cfg.Database.URL)
	if err != nil {
		return nil, fmt.Errorf("parse database url: %w", err)
	}
	poolCfg.MaxConns = cfg.Database.MaxConns

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("create db pool: %w", err)
	}

	return pool, nil
}
