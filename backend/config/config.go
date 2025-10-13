package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/goawwer/yamyard/internal/adataper/database"
	"github.com/goawwer/yamyard/internal/middleware"
	"github.com/goawwer/yamyard/pkg/logger"
	"github.com/goawwer/yamyard/pkg/server"
)

type Config struct {
	Database database.Config
	Logger   logger.Config
	Server   server.Config
	JWT middleware.Config
}

func New() (*Config, error) {
	var cfg Config

	cfg, err := env.ParseAs[Config]()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	return &cfg, nil
}
