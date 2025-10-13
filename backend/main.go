package main

import (
	"context"

	"github.com/goawwer/yamyard/config"
	"github.com/goawwer/yamyard/internal/app"
	"github.com/goawwer/yamyard/pkg/logger"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		panic("failed to load config file")
	}

	logger.Init(&cfg.Logger)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	app.Start(ctx, cfg)
}
