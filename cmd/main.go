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

	// 启动一个 goroutine 来写入测试日志
	go func() {
		for i := 0; i < 10000; i++ {
			logger.Info(fmt.Sprintf("Test log entry %d", i))
		}
	}()

	// 在 main 函数中，application.Run() 之前添加：
	go func() {
		for i := 0; i < 100000; i++ {
			logger.Info(fmt.Sprintf("Test log entry %d: This is a longer log message to help reach the 1MB threshold faster.", i))
		}
	}()

	if err := application.Run(); err != nil {
		fmt.Printf("Failed to run app: %v\n", err)
		logger.Fatal("Failed to run app", zap.Error(err))
	}
}
