package logger

import (
	"context"

	"go.uber.org/zap/zapcore"
)

type Producer interface {
	SendLog(ctx context.Context, level zapcore.Level, message string) error
	Close(timeoutMs int64)
}

type Event interface {
	Short() string
	Long() string
}
