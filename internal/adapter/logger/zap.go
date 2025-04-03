package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Zap *zap.Logger
var err error

func ZapInit() {
	zapconfig := zap.NewProductionConfig()
	zapconfig.EncoderConfig.TimeKey = "timestamp"
	zapconfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	zapconfig.EncoderConfig.StacktraceKey = ""
	Zap, err = zapconfig.Build(zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}
}

func Info(message string, felids ...zap.Field) {
	Zap.Info(message, felids...)
}
func Dubug(message string, felids ...zap.Field) {
	Zap.Debug(message, felids...)
}
func Error(message interface{}, felids ...zap.Field) {
	switch err := message.(type) {
	case error:
		Zap.Error(err.Error(), felids...)
	case string:
		Zap.Error(err, felids...)
	}
}
func Warning(message string, felids ...zap.Field) {
	Zap.Warn(message, felids...)
}
func Panic(message string, felids ...zap.Field) {
	Zap.Panic(message, felids...)
}
