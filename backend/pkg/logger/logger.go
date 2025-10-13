package logger

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

type Config struct {
	File string `env:"LOG_FILE"`
}

var logger *logrus.Logger

func Init(c *Config) error {
	log := logrus.New()

	file, err := os.OpenFile(c.File, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("failed to load logger file: %w", err)
	}

	log.SetOutput(file)

	logger = log

	return nil
}

func Debug(args ...any) {
	logger.Debug(args...)
}

func Info(args ...any) {
	logger.Info(args...)
}

func Error(args ...any) {
	logger.Error(args...)
}

func Fatal(args ...any) {
	logger.Fatal(args...)
}
