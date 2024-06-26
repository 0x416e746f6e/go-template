package logutils

import (
	"errors"
	"fmt"
	"strings"

	"github.com/{{ cookiecutter.github_org }}/{{ cookiecutter.project_slug }}/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	ErrLoggerFailedToBuild = errors.New("failed to build the logger")
	ErrLoggerInvalidLevel  = errors.New("invalid log-level")
	ErrLoggerInvalidMode   = errors.New("invalid log-mode")
)

func NewLogger(cfg *config.Log) (
	*zap.Logger, error,
) {
	var config zap.Config
	switch strings.ToLower(cfg.Mode) {
	case "dev":
		config = zap.NewDevelopmentConfig()
	case "prod":
		config = zap.NewProductionConfig()
	default:
		return nil, fmt.Errorf("%w: %s",
			ErrLoggerInvalidMode, cfg.Mode,
		)
	}
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logLevel, err := zap.ParseAtomicLevel(cfg.Level)
	if err != nil {
		return nil, fmt.Errorf("%w: %s: %w",
			ErrLoggerInvalidLevel, cfg.Level, err,
		)
	}
	config.Level = logLevel

	l, err := config.Build()
	if err != nil {
		return nil, fmt.Errorf("%w: %w",
			ErrLoggerFailedToBuild, err,
		)
	}

	return l, nil
}
