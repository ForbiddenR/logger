package logger

import (
	"go.uber.org/zap/zapcore"
)

type SLSHook interface {
	Fire(e zapcore.Entry) error
	Level() zapcore.Level
}

var slsHook SLSHook


func NewSLSHook(hook SLSHook) {
	if hook == nil {
		panic("null hook is illegal")
	}
	slsHook = hook
}