package grpcserver

import (
	"context"
	"log/slog"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UnaryLoggingInterceptor logs method, duration, and status code for every unary RPC.
func UnaryLoggingInterceptor(log *slog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		start := time.Now()
		resp, err := handler(ctx, req)
		code := status.Code(err)

		attrs := []any{
			slog.String("method", info.FullMethod),
			slog.Duration("duration", time.Since(start)),
			slog.String("code", code.String()),
		}
		if err != nil && code != codes.NotFound {
			attrs = append(attrs, slog.String("error", err.Error()))
		}

		if code == codes.OK || code == codes.NotFound {
			log.InfoContext(ctx, "rpc", attrs...)
		} else {
			log.ErrorContext(ctx, "rpc", attrs...)
		}

		return resp, err
	}
}

// Metrics holds the Prometheus instruments for gRPC request observability.
type Metrics struct {
	requests *prometheus.CounterVec
	duration *prometheus.HistogramVec
}

func NewMetrics(reg prometheus.Registerer) *Metrics {
	factory := promauto.With(reg)
	return &Metrics{
		requests: factory.NewCounterVec(prometheus.CounterOpts{
			Name: "wallet_grpc_requests_total",
			Help: "Total number of gRPC requests by method and status code.",
		}, []string{"method", "code"}),

		duration: factory.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "wallet_grpc_request_duration_seconds",
			Help:    "gRPC request latency in seconds by method.",
			Buckets: prometheus.DefBuckets,
		}, []string{"method"}),
	}
}

// UnaryMetricsInterceptor records request count and latency for every unary RPC.
func UnaryMetricsInterceptor(m *Metrics) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		start := time.Now()
		resp, err := handler(ctx, req)

		m.requests.WithLabelValues(info.FullMethod, status.Code(err).String()).Inc()
		m.duration.WithLabelValues(info.FullMethod).Observe(time.Since(start).Seconds())

		return resp, err
	}
}
