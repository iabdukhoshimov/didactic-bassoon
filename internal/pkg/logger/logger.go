package logger

import (
	"gitlab.com/tour/pkg/logger"

	"go.uber.org/zap"
)

var (
	Log *zap.Logger
)

func SetLogger(cfg *logger.LoggingConfig) *zap.Logger {
	sugar := zap.Must(zap.NewProduction()).Sugar()

	defer sugar.Sync()
	sugar.Infow("Hello from Sugared Logger!")

	_logger := sugar.Desugar()

	_logger.Info("Hello from Logger!")

	Log = _logger

	return Log
}
