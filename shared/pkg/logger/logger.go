package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

// Initialize initializes the global logger
func Initialize(environment string) error {
	var config zap.Config

	if environment == "production" {
		config = zap.NewProductionConfig()
		config.EncoderConfig.TimeKey = "timestamp"
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	} else {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	logger, err := config.Build()
	if err != nil {
		return err
	}

	Log = logger
	return nil
}

// Info logs an info message
func Info(message string, fields ...zap.Field) {
	if Log != nil {
		Log.Info(message, fields...)
	}
}

// Error logs an error message
func Error(message string, fields ...zap.Field) {
	if Log != nil {
		Log.Error(message, fields...)
	}
}

// Debug logs a debug message
func Debug(message string, fields ...zap.Field) {
	if Log != nil {
		Log.Debug(message, fields...)
	}
}

// Warn logs a warning message
func Warn(message string, fields ...zap.Field) {
	if Log != nil {
		Log.Warn(message, fields...)
	}
}

// Fatal logs a fatal message and exits
func Fatal(message string, fields ...zap.Field) {
	if Log != nil {
		Log.Fatal(message, fields...)
	}
}

// Sync flushes any buffered log entries
func Sync() {
	if Log != nil {
		_ = Log.Sync()
	}
}
