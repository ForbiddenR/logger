package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var consuleWS = zapcore.Lock(os.Stdout)

type Logger struct {
	*zap.Logger
	opts *Options
}

func InitLogger() *Logger {
	logger := &Logger{
		opts: &Options{
			level: zap.DebugLevel,
		},
	}

	logger.Logger = zap.New(logger.cores(),
		zap.AddCaller(),
	)

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
