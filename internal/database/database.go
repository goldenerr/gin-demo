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

type ZapGormWriter struct {
	Logger *zap.Logger
}

func (w ZapGormWriter) Printf(format string, args ...interface{}) {
	w.Logger.Sugar().Infof(format, args...)
}

type TracingLogger struct {
	ZapGormWriter
	LogLevel gormlogger.LogLevel
}

func (l TracingLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	newLogger := l
	newLogger.LogLevel = level
	return &newLogger
}

func (l TracingLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Info {
		l.Logger.Info(fmt.Sprintf(msg, data...), zap.String("requestID", tracing.FromContext(ctx)))
	}
}

func (l TracingLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Warn {
		l.Logger.Warn(fmt.Sprintf(msg, data...), zap.String("requestID", tracing.FromContext(ctx)))
	}
}

func (l TracingLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Error {
		l.Logger.Error(fmt.Sprintf(msg, data...), zap.String("requestID", tracing.FromContext(ctx)))
	}
}

func (l TracingLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= gormlogger.Silent {
		return
	}

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
		l.Logger.Error("trace", fields...)
	} else {
		l.Logger.Info("trace", fields...)
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
		gormConfig.Logger = &TracingLogger{
			ZapGormWriter: ZapGormWriter{Logger: logger.GetLogger()},
			LogLevel:      logLevel,
		}
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
