package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"sync"
)

var (
	Log  *zap.Logger
	once sync.Once
)

func InitLogger() *zap.Logger {
	once.Do(func() {
		config := zap.NewProductionConfig()

		config.EncoderConfig.TimeKey = "timestamp"
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

		var err error
		Log, err = config.Build()
		if err != nil {
			panic("failed to initialize zap logger: " + err.Error())
		}
	})

	return Log
}

func Sync() {
	if Log != nil {
		err := Log.Sync()
		if err != nil {
			Log.Error("logs didn't flush", zap.Error(err))
		}
	}
}
