package app

import (
	"gin-demo/config"
	"gin-demo/internal/database"
	"gin-demo/internal/logger"
	"gin-demo/internal/router"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type App struct {
	Router *gin.Engine
}

func NewApp() (*App, error) {
	logger.Info("Initializing application")

	// 初始化数据库
	logger.Info("Initializing database")
	if err := database.Init(); err != nil {
		logger.Error("Failed to initialize database", zap.Error(err))
		return nil, err
	}
	logger.Info("Database initialized successfully")

	// 设置 Gin 的模式
	gin.SetMode(gin.ReleaseMode)
	logger.Info("Gin mode set to release")

	// 初始化路由
	logger.Info("Setting up router")
	r := router.SetupRouter()
	logger.Info("Router setup completed")

	logger.Info("Application initialization completed")

	return &App{
		Router: r,
	}, nil
}

func (a *App) Run() error {
	port := config.Get().Server.Port
	logger.Info("Starting server", zap.String("port", port))
	err := a.Router.Run(":" + port)
	if err != nil {
		logger.Error("Server failed to start", zap.Error(err))
	}
	return err
}
