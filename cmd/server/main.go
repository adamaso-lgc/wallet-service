package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/adamaso/wallet-service/internal/bootstrap"
)

func main() {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "local"
	}

	log := bootstrap.NewLogger(env)

	cfg := bootstrap.MustLoad(env)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	app, err := bootstrap.New(ctx, cfg, log)
	if err != nil {
		log.Error("failed to initialise app", "error", err)
		os.Exit(1)
	}

	if err := app.Run(ctx); err != nil {
		log.Error("app stopped with error", "error", err)
		os.Exit(1)
	}
}
