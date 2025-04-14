package logger

import (
	"context"
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const NonLevel zapcore.Level = 99

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
	enabler     zap.AtomicLevel
	slsProducer Producer
	opts        *Options
}

func InitLogger(level zapcore.Level, enabler zap.AtomicLevel, producer Producer, cs ...zapcore.Core) *Logger {
	logger := &Logger{
		enabler:     enabler,
		slsProducer: producer,
		opts: &Options{
			level: level,
		},
	}
	logger.Logger = zap.New(logger.cores(cs...))
	return logger
}

func (l *Logger) ChangeProducerLevel(level zapcore.Level) {
	l.enabler.SetLevel(level)
}

func (l *Logger) ChangeLoggerLevel(level zapcore.Level) {
	l.opts.mutableLevel.SetLevel(level)
}

func (l *Logger) Close() {
	l.slsProducer.Close(100000)
}

func (l *Logger) Infos(template string, args ...any) {
	l.Sugar().Infof(template, args...)
}

func (l *Logger) Info(ctx context.Context, event Event) {
	l.output(ctx, event, zap.InfoLevel, l.Logger.Info)
}

func (l *Logger) Error(ctx context.Context, event Event) {
	l.output(ctx, event, zap.ErrorLevel, l.Logger.Error)
}

func (l Logger) output(ctx context.Context, event Event, level zapcore.Level, fn func(string, ...zap.Field)) {
	if l.slsProducer != nil && l.enabler.Enabled(level) {
		l.slsProducer.SendLog(ctx, zap.InfoLevel, event.Short())
	}
	if l.Logger.Level().Enabled(level) {
		fn(event.Long())
	}
}

func (l *Logger) Debugf(ctx context.Context, template string, args ...any) {
	if l.slsProducer != nil && l.enabler.Enabled(zapcore.InfoLevel) {
		l.slsProducer.SendLog(ctx, zapcore.InfoLevel, fmt.Sprintf(template, args...))
	}
	l.Sugar().Debugf(template, args...)
}

func (l *Logger) Infof(ctx context.Context, template string, args ...any) {
	if l.slsProducer != nil && l.enabler.Enabled(zapcore.InfoLevel) {
		l.slsProducer.SendLog(ctx, zapcore.InfoLevel, fmt.Sprintf(template, args...))
	}
	l.Sugar().Debugf(template, args...)
}

func (l *Logger) Warnf(ctx context.Context, template string, args ...any) {
	if l.slsProducer != nil && l.enabler.Enabled(zapcore.WarnLevel) {
		l.slsProducer.SendLog(ctx, zapcore.WarnLevel, fmt.Sprintf(template, args...))
	}
	l.Sugar().Warnf(template, args...)
}

func (l *Logger) Errorf(ctx context.Context, template string, args ...any) {
	if l.slsProducer != nil && l.enabler.Enabled(zapcore.ErrorLevel) {
		l.slsProducer.SendLog(ctx, zapcore.ErrorLevel, fmt.Sprintf(template, args...))
	}
	l.Sugar().Errorf(template, args...)
}

func (l *Logger) cores(cs ...zapcore.Core) zapcore.Core {
	cores := make([]zapcore.Core, 0, 1+len(cs))
	l.opts.mutableLevel = zap.NewAtomicLevelAt(l.opts.level)
	cores = append(cores, zapcore.NewCore(zapcore.NewConsoleEncoder(newConsoleEncoderConfig()), consuleWS, l.opts.mutableLevel))
	if len(cs) > 0 {
		cores = append(cores, cs...)
	}
	return zapcore.NewTee(cores...)
}
