package logger

import (
	"fmt"

	"go.uber.org/zap"
)

const (
	LoggerLevelINFO  = "INFO"
	LoggerLevelDEBUG = "DEBUG"
)

var Log *zap.Logger = zap.NewNop()

// LoggerOption описывает функциональную опцию для конфигурации логгера.
type LoggerOption func(*zap.Config)

// Init инициализирует синглтон логера с необходимыми опциями
func Init(options ...LoggerOption) error {
	cfg := zap.NewProductionConfig()

	for _, opt := range options {
		opt(&cfg)
	}

	zl, err := cfg.Build()
	if err != nil {
		return fmt.Errorf("logger-Init-cfg.Build-err: %w", err)
	}

	Log = zl
	return nil
}

// WithLevel задает уровень логирования.
func WithLevel(level string) LoggerOption {
	return func(cfg *zap.Config) {
		lvl, err := zap.ParseAtomicLevel(level)
		if err == nil {
			cfg.Level = lvl
		}
	}
}
