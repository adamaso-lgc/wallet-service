package bootstrap

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	walletv1 "github.com/adamaso/wallet-service/gen/proto/v1"
	"github.com/adamaso/wallet-service/internal/application"
	"github.com/adamaso/wallet-service/internal/config"
	grpcserver "github.com/adamaso/wallet-service/internal/grpc"
	"github.com/adamaso/wallet-service/internal/infrastructure/postgres"
)

type Service struct {
	grpcServer    *grpc.Server
	healthServer  *health.Server
	metricsServer *http.Server
	pool          *pgxpool.Pool
	grpcAddr      string
	log           *slog.Logger
}

func New(ctx context.Context, cfg config.Config, log *slog.Logger) (*Service, error) {
	pool, err := newPool(ctx, cfg, log)
	if err != nil {
		return nil, err
	}

	repo := postgres.NewWalletRepository(pool)
	store := postgres.NewWalletViewStore(pool)
	app := application.NewApplication(repo, store)
	srv := grpcserver.NewWalletHandler(app)

	reg := prometheus.NewRegistry()
	reg.MustRegister(collectors.NewGoCollector(), collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	metrics := grpcserver.NewMetrics(reg)

	grpcSrv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpcserver.UnaryLoggingInterceptor(log),
			grpcserver.UnaryMetricsInterceptor(metrics),
		),
	)
	walletv1.RegisterWalletServiceServer(grpcSrv, srv)

	healthSrv := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcSrv, healthSrv)
	healthSrv.SetServingStatus(walletv1.WalletService_ServiceDesc.ServiceName, grpc_health_v1.HealthCheckResponse_SERVING)

	reflection.Register(grpcSrv)

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	metricsSrv := &http.Server{Addr: cfg.MetricsAddr(), Handler: mux}

	return &Service{
		grpcServer:    grpcSrv,
		healthServer:  healthSrv,
		metricsServer: metricsSrv,
		pool:          pool,
		grpcAddr:      cfg.GRPCAddr(),
		log:           log,
	}, nil
}

func (s *Service) Run(ctx context.Context) error {
	lis, err := net.Listen("tcp", s.grpcAddr)
	if err != nil {
		return fmt.Errorf("listen %s: %w", s.grpcAddr, err)
	}

	go func() {
		s.log.Info("metrics server started", slog.String("addr", s.metricsServer.Addr))
		if err := s.metricsServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.log.Error("metrics server error", slog.Any("error", err))
		}
	}()

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

	_ = s.metricsServer.Shutdown(ctx)

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

func newPool(ctx context.Context, cfg config.Config, log *slog.Logger) (*pgxpool.Pool, error) {
	poolCfg, err := pgxpool.ParseConfig(cfg.Database.URL)
	if err != nil {
		return nil, fmt.Errorf("parse database url: %w", err)
	}
	poolCfg.MaxConns = cfg.Database.MaxConns

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("create db pool: %w", err)
	}

	if err := pingWithRetry(ctx, pool, log); err != nil {
		pool.Close()
		return nil, err
	}

	return pool, nil
}

func pingWithRetry(ctx context.Context, pool *pgxpool.Pool, log *slog.Logger) error {
	backoff := time.Second
	for {
		if err := pool.Ping(ctx); err == nil {
			return nil
		}

		log.Warn("database not ready, retrying",
			slog.Duration("backoff", backoff),
		)

		select {
		case <-ctx.Done():
			return fmt.Errorf("gave up waiting for database: %w", ctx.Err())
		case <-time.After(backoff):
		}

		if backoff < 16*time.Second {
			backoff *= 2
		}
	}
}
