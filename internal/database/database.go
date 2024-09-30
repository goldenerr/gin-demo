package database

import (
	"context"
	"fmt"
	"gin-demo/config"
	"gin-demo/internal/logger"
	"gin-demo/internal/tracing"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

var DB *gorm.DB

type ContextLogger struct {
	ZapLogger *zap.Logger
}

func (l ContextLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	return l
}

func (l ContextLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	l.ZapLogger.Info(fmt.Sprintf(msg, data...), zap.String("requestID", tracing.FromContext(ctx)))
}

func (l ContextLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	l.ZapLogger.Warn(fmt.Sprintf(msg, data...), zap.String("requestID", tracing.FromContext(ctx)))
}

func (l ContextLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	l.ZapLogger.Error(fmt.Sprintf(msg, data...), zap.String("requestID", tracing.FromContext(ctx)))
}

func (l ContextLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()
	fields := []zap.Field{
		zap.String("requestID", tracing.FromContext(ctx)),
		zap.Duration("elapsed", elapsed),
		zap.Int64("rows", rows),
		zap.String("sql", sql),
	}
	if err != nil {
		fields = append(fields, zap.Error(err))
		l.ZapLogger.Error("SQL query failed", fields...)
	} else {
		l.ZapLogger.Info("SQL query", fields...)
	}
}

func Init() error {
	cfg := config.Get().Database

	logger.Info("Initializing database",
		zap.Bool("logSQL", cfg.LogSQL),
		zap.String("logLevel", cfg.LogLevel))

	gormConfig := &gorm.Config{}

	if cfg.LogSQL {
		logger.Info("SQL logging is enabled")
		logLevel := gormlogger.Info
		if cfg.LogLevel == "silent" {
			logLevel = gormlogger.Silent
		} else if cfg.LogLevel == "error" {
			logLevel = gormlogger.Error
		} else if cfg.LogLevel == "warn" {
			logLevel = gormlogger.Warn
		}
		gormConfig.Logger = ContextLogger{ZapLogger: logger.GetLogger()}.LogMode(logLevel)
		logger.Info("GORM logger configured", zap.String("logLevel", cfg.LogLevel))
	} else {
		logger.Info("SQL logging is disabled")
	}

	var err error
	DB, err = gorm.Open(mysql.Open(cfg.DSN), gormConfig)
	if err != nil {
		logger.Error("Failed to connect to database", zap.Error(err))
		return err
	}

	// 测试数据库连接
	sqlDB, err := DB.DB()
	if err != nil {
		logger.Error("Failed to get database connection", zap.Error(err))
		return err
	}
	if err := sqlDB.Ping(); err != nil {
		logger.Error("Failed to ping database", zap.Error(err))
		return err
	}

	logger.Info("Database connection established", zap.String("dsn", cfg.DSN))

	return nil
}

// 删除 getGormLogLevel 函数，因为我们现在直接在 NewGormLogger 中处理日志级别
