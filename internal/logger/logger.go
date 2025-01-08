package logger

import (
	"fmt"

	"go.uber.org/zap"
)

func Initialize(level string) (*zap.Logger, error) {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return nil, fmt.Errorf("Initialize->ParseAtomicLevel: %w",
			err)
	}

	cfg := zap.NewProductionConfig()

	cfg.Level = lvl

	zl, err := cfg.Build()
	if err != nil {
		return nil, fmt.Errorf("Initialize->Build: %w", err)
	}

	return zl, nil
}

func DoInfoLogFromErr(
	method string,
	err error,
	logger *zap.Logger,
) {
	logger.Info(method,
		zap.String("err:", err.Error()),
	)
}

func DoInfoLogFromStr(
	method string,
	txt string,
	logger *zap.Logger,
) {
	logger.Info(method,
		zap.String("err:", txt),
	)
}
