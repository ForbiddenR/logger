package logger

import "go.uber.org/zap/zapcore"

type Options struct {
	appName    string
	appVersion string
	level      zapcore.Level
}
