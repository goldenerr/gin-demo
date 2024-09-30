package logger

import (
	"gin-demo/config"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

// Init initializes the logger
func Init() error {
	cfg := config.Get().Log
	var err error

	logConfig := zap.NewProductionConfig()

	// 检查日志文件路径是否为空
	if cfg.Filename == "" {
		cfg.Filename = "./logs/app.log" // 设置默认日志文件路径
	}

	// 确保日志文件目录存在
	err = os.MkdirAll(filepath.Dir(cfg.Filename), 0755)
	if err != nil {
		return err
	}

	logConfig.OutputPaths = []string{cfg.Filename}
	logConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// Set log level
	switch cfg.Level {
	case "debug":
		logConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		logConfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		logConfig.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		logConfig.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	default:
		logConfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	logger, err = logConfig.Build()
	if err != nil {
		return err
	}

	zap.ReplaceGlobals(logger)
	return nil
}

// GetLogger returns the global logger instance
func GetLogger() *zap.Logger {
	return logger
}

// Info logs a message at InfoLevel
func Info(msg string, fields ...zap.Field) {
	logger.Info(msg, fields...)
}

// Error logs a message at ErrorLevel
func Error(msg string, fields ...zap.Field) {
	logger.Error(msg, fields...)
}

// Debug logs a message at DebugLevel
func Debug(msg string, fields ...zap.Field) {
	logger.Debug(msg, fields...)
}

// Warn logs a message at WarnLevel
func Warn(msg string, fields ...zap.Field) {
	logger.Warn(msg, fields...)
}

// Fatal logs a message at FatalLevel and then calls os.Exit(1)
func Fatal(msg string, fields ...zap.Field) {
	logger.Fatal(msg, fields...)
}
