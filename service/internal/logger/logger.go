package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitLogger() *zap.SugaredLogger {
	config := zap.Config{
		Encoding:         "console", // JSON-формат логов
		Level:            zap.NewAtomicLevelAt(zap.InfoLevel),
		OutputPaths:      []string{"stdout", "../logs/app.log"}, // Выводим и в консоль, и в файл
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:     "timestamp",
			LevelKey:    "level",
			MessageKey:  "msg",
			EncodeTime:  zapcore.ISO8601TimeEncoder, // Читаемая дата-время
			EncodeLevel: zapcore.LowercaseLevelEncoder,
		},
	}

	l, _ := config.Build()

	logger := l.Sugar()

	return logger
}
