package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/goawwer/yamyard/pkg/logger"
	_ "github.com/lib/pq"
)

type Config struct {
	Name     string `env:"DB_NAME"`
	Host     string `env:"DB_HOST"`
	Port     int    `env:"DB_PORT"`
	User     string `env:"DB_USER"`
	Password string `env:"DB_PASSWORD"`

	URL string `env:"DB_URL"`
}

type Db struct {
	*sql.DB
}

var database *Db

func Init(ctx context.Context, c *Config) error {
	db, err := sql.Open("postgres", c.URL)
	if err != nil {
		return fmt.Errorf("failed to parse arguments: %w", err)
	}

	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	logger.Info("database connection successfully created")

	database = &Db{db}

	return nil
}

func Get() *Db {
	return database
}

func Close() error {
	return database.Close()
}
