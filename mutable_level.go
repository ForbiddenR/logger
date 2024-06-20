package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	defaultLoggerLevel zap.AtomicLevel
)

func ChangeLoggerLevel(level zapcore.Level) {
	defaultLoggerLevel.SetLevel(level)
}
