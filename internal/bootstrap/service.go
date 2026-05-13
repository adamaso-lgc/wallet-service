package bootstrap

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	walletv1 "github.com/adamaso/wallet-service/gen/proto/v1"
	"github.com/adamaso/wallet-service/internal/application"
	"github.com/adamaso/wallet-service/internal/config"
	grpcserver "github.com/adamaso/wallet-service/internal/grpc"
	"github.com/adamaso/wallet-service/internal/infrastructure/postgres"
)

type Service struct {
	grpcServer *grpc.Server
	pool       *pgxpool.Pool
	grpcAddr   string
}

func New(ctx context.Context, cfg config.Config) (*Service, error) {
	pool, err := newPool(ctx, cfg)
	if err != nil {
		return nil, err
	}

	repo := postgres.NewWalletRepository(pool)
	store := postgres.NewWalletViewStore(pool)

	app := application.NewApplication(repo, store)
	srv := grpcserver.NewServer(app)

	grpcSrv := grpc.NewServer()
	walletv1.RegisterWalletServiceServer(grpcSrv, srv)
	reflection.Register(grpcSrv)

	return &Service{
		grpcServer: grpcSrv,
		pool:       pool,
		grpcAddr:   cfg.GRPCAddr(),
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

	fmt.Printf("gRPC server listening on %s\n", s.grpcAddr)

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
