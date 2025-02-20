package logwriter

import (
	"go.uber.org/zap"
)

func NewZapLogger() *zap.Logger {
	logger := zap.Must(zap.NewProduction())
	defer logger.Sync()
	logger.Info("Initialing zap as a logger")
	return logger
}
