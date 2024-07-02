package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Options struct {
	level        zapcore.Level
	mutableLevel zap.AtomicLevel
}
