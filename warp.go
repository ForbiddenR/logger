package logger

import "go.uber.org/zap"

type Field = zap.Field

func String(key string, val string) Field {
	return zap.String(key, val)
}

func Int(key string, val int) Field {
	return zap.Int(key, val)
}

func Int64(key string, val int64) Field {
	return zap.Int64(key, val)
}

func Uint64(key string, val uint64) Field {
	return zap.Uint64(key, val)
}