package router

import (
	"gin-demo/internal/handler"
	"gin-demo/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.Logger()) // 使用新的日志中间件

	v1 := r.Group("/api/v1")
	{
		v1.GET("/users", handler.GetUsers)
	}

	return r
}
