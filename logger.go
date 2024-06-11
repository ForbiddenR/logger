package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var consuleWS = zapcore.Lock(os.Stdout)

type DefaultLogger = zap.Logger

func InitYXLogger() *DefaultLogger {
	return zap.New(cores(),
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.WarnLevel))
}

func InitDefaultLogger(hooks ...func(zapcore.Entry) error) *DefaultLogger {
	options := []zap.Option{
		// zap.AddCaller(),
	}
	if len(hooks) > 0 {
		options = append(options, zap.Hooks(hooks...))
	}
	return zap.New(cores(),
		options...,
	)
}

func cores() zapcore.Core {
	cores := []zapcore.Core{
		zapcore.NewCore(zapcore.NewConsoleEncoder(newConsoleEncoderConfig()), consuleWS, zap.DebugLevel),
	}
	return zapcore.NewTee(cores...)
}

type Logger struct {
	*zap.Logger
	opts *Options
}

func InitLogger(hooks ...func(zapcore.Entry) error) *Logger {
	logger := &Logger{
		opts: &Options{
			level: zap.DebugLevel,
		},
	}
	if len(hooks) > 0 && hooks[0] != nil {
		logger.Logger = zap.New(logger.cores(), zap.Hooks(hooks...))
	} else {
		logger.Logger = zap.New(logger.cores())
	}

	return logger
}

func (l *Logger) Debug(template string, args ...any) {
	l.Sugar().Debugf(template, args...)
}

func (l *Logger) Warn(template string, args ...any) {
	l.Sugar().Warnf(template, args...)
}

func (l *Logger) Info(template string, args ...any) {
	l.Sugar().Infof(template, args...)
}

func (l *Logger) Error(err error) {
	l.Sugar().Error(err.Error())
}

func (l *Logger) Panic(template string, args ...any) {
	l.Sugar().Panicf(template, args...)
}

func (l *Logger) cores() zapcore.Core {
	cores := make([]zapcore.Core, 0)

	cores = append(cores, zapcore.NewCore(zapcore.NewConsoleEncoder(newConsoleEncoderConfig()), consuleWS, zap.DebugLevel))

	return zapcore.NewTee(cores...)
}
