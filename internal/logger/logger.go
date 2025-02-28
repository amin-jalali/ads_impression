package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

func InitLogger() {
	config := zap.NewProductionConfig()

	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	var err error
	Log, err = config.Build()
	if err != nil {
		panic("failed to initialize zap logger: " + err.Error())
	}
}

func Sync() {
	err := Log.Sync() // Ensure logs are flushed
	if err != nil {
		Log.Error("logs didn't flushed")
	}
}
