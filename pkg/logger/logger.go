package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

// InitLogger initializes the global logger with a specific log level
func InitLogger(level string) {
	var logLevel zapcore.Level
	switch level {
	case "debug":
		logLevel = zap.DebugLevel
	case "info":
		logLevel = zap.InfoLevel
	default:
		logLevel = zap.InfoLevel
	}

	cfg := zap.Config{
		Encoding:         "json", // Use "console" for human-readable logs
		Level:            zap.NewAtomicLevelAt(logLevel),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig:    zap.NewProductionEncoderConfig(),
	}
	var err error
	Logger, err = cfg.Build()
	if err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}
}

// Info is a shortcut for logging info level messages
func Info(msg string, fields ...zap.Field) {
	Logger.Info(msg, fields...)
}

// Debug is a shortcut for logging debug level messages
func Debug(msg string, fields ...zap.Field) {
	Logger.Debug(msg, fields...)
}

// SyncLogger ensures all buffered log entries are flushed
func SyncLogger() {
	_ = Logger.Sync() // flushes buffer, if any
}
