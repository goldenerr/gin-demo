package main

import (
	"fmt"
	"gin-demo/config"
	"gin-demo/internal/app"
	"gin-demo/internal/logger"
	"os"

	"go.uber.org/zap"
)

func main() {
	// 初始化配置
	if err := config.Init(); err != nil {
		fmt.Printf("Failed to initialize config: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Config initialized successfully")

	// 初始化日志
	if err := logger.Init(); err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Logger initialized successfully")

	// 使用 defer 来确保在程序退出时刷新日志
	defer logger.GetLogger().Sync()

	application, err := app.NewApp()
	if err != nil {
		fmt.Printf("Failed to initialize app: %v\n", err)
		logger.Fatal("Failed to initialize app", zap.Error(err))
	}
	fmt.Println("App initialized successfully")

	if err := application.Run(); err != nil {
		fmt.Printf("Failed to run app: %v\n", err)
		logger.Fatal("Failed to run app", zap.Error(err))
	}
}
