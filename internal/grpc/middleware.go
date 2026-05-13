package grpcserver

import (
	"context"
	"log/slog"
	"time"

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
