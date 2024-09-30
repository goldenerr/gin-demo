package logger

import (
	"fmt"
	"gin-demo/config"
	"os"
	"path/filepath"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

// Init initializes the logger
func Init() error {
	cfg := config.Get().Log

	// 创建日志目录
	logDir := filepath.Dir(cfg.Filename)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("can't create log directory: %w", err)
	}

	// 设置日志级别
	var level zapcore.Level
	switch cfg.Level {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	default:
		level = zap.InfoLevel
	}

	// 计算日志文件大小（以MB为单位）
	maxSize := calculateMaxSize(cfg.MaxSize, cfg.MaxSizeUnit)

	// 创建 lumberjack logger
	lumberjackLogger := &lumberjack.Logger{
		Filename:   getLogFilename(cfg.Filename),
		MaxSize:    maxSize, // 现在这个值是 KB
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		Compress:   cfg.Compress,
		LocalTime:  true, // 使用本地时间（北京时间）
	}

	// 创建 encoder 配置
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// 创建 core
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(lumberjackLogger),
		level,
	)

	// 创建logger
	logger = zap.New(core)
	zap.ReplaceGlobals(logger)

	return nil
}

// calculateMaxSize 根据配置的单位计算日志文件的最大大小（以MB为单位）
func calculateMaxSize(size int, unit string) int {
	switch unit {
	case "K":
		return size / 1024 // 转换为MB
	case "M":
		return size // 直接返回MB值
	case "G":
		return size * 1024 // 转换为MB
	case "T":
		return size * 1024 * 1024 // 转换为MB
	default:
		return size // 默认假设为MB
	}
}

// getLogFilename 根据配置的文件名模式生成实际的日志文件名
func getLogFilename(filenamePattern string) string {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	now := time.Now().In(loc)
	return now.Format(filenamePattern)
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
