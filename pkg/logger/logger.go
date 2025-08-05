package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.Logger

func InitLogger(level string) error {
	var zapLevel zapcore.Level
	switch level {
	case "debug":
		zapLevel = zapcore.DebugLevel
	case "info":
		zapLevel = zapcore.InfoLevel
	case "warn":
		zapLevel = zapcore.WarnLevel
	case "error":
		zapLevel = zapcore.ErrorLevel
	default:
		zapLevel = zapcore.InfoLevel
	}

	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(zapLevel)
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.StacktraceKey = "stacktrace"

	var err error
	log, err = config.Build()
	if err != nil {
		return err
	}

	zap.ReplaceGlobals(log)
	return nil
}

func GetLogger() *zap.Logger {
	if log == nil {
		if err := InitLogger("info"); err != nil {
			// Если не удалось инициализировать, создаём базовый логгер
			log, _ = zap.NewProduction()
		}
	}
	return log
}

func Sync() {
	if log != nil {
		if err := log.Sync(); err != nil {
			// Логируем ошибку синхронизации, но не паникуем
			log.Error("failed to sync logger", zap.Error(err))
		}
	}
}
