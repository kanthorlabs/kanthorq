package xlogger

import (
	"os"

	"go.uber.org/zap"
)

func New() *zap.Logger {
	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	logger, err := zap.NewProduction(zap.Fields(zap.String("hostname", hostname)))
	if err != nil {
		panic(err)
	}
	return logger
}

func NewNoop() *zap.Logger {
	return zap.NewNop()
}
