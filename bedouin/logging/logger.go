package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var Log *zap.Logger

func init() {
	config := zap.NewProductionConfig()

	// Set log level from environment variable
	level := os.Getenv("LOG_LEVEL")
	if level == "debug" {
		config.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	} else {
		config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	}

	// Configure log output format
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// Build logger
	logger, err := config.Build()
	if err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}

	Log = logger
}

// Utility function to log errors with additional context
func LogError(err error, message string, fields ...zap.Field) {
	if err != nil {
		Log.Error(message, append(fields, zap.Error(err))...)
	}
}
