package logger

import (
	"go.uber.org/zap/zapcore"
)

type SLSHook interface {
	Fire(e zapcore.Entry) error
	Level() zapcore.Level
}

var slsHook SLSHook