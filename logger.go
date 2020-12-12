package hkmh

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(level zapcore.Level) (*zap.SugaredLogger, error) {
	config := zap.NewProductionConfig()
	config.Encoding = "console"
	config.Level = zap.NewAtomicLevelAt(level)
	logger, err := config.Build()
	if err != nil {
		return nil, err
	}
	defer logger.Sync()
	return logger.Sugar(), nil
}
