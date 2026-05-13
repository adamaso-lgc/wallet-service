package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/adamaso/wallet-service/internal/bootstrap"
	"github.com/adamaso/wallet-service/internal/config"
)

func main() {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "local"
	}

	cfg := config.MustLoad(env)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	app, err := bootstrap.New(ctx, cfg)
	if err != nil {
		log.Fatalf("failed to initialise app: %v", err)
	}

	if err := app.Run(ctx); err != nil {
		log.Fatalf("app stopped with error: %v", err)
	}
}
