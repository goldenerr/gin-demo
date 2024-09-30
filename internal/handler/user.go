package handler

import (
	"gin-demo/internal/database"
	"gin-demo/internal/logger"
	"gin-demo/internal/tracing"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type User struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

func GetUsers(c *gin.Context) {
	ctx := c.Request.Context()
	requestID := tracing.FromContext(ctx)
	logger.Info("GetUsers handler called", zap.String("requestID", requestID))

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		logger.Error("Failed to parse page parameter", zap.Error(err), zap.String("requestID", requestID))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page parameter"})
		return
	}

	size, err := strconv.Atoi(c.DefaultQuery("size", "10"))
	if err != nil {
		logger.Error("Failed to parse size parameter", zap.Error(err), zap.String("requestID", requestID))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid size parameter"})
		return
	}

	logger.Info("Fetching users",
		zap.Int("page", page),
		zap.Int("size", size),
		zap.String("requestID", requestID))

	var users []User
	result := database.DB.WithContext(ctx).Offset((page - 1) * size).Limit(size).Find(&users)
	if result.Error != nil {
		logger.Error("Failed to fetch users from database",
			zap.Error(result.Error),
			zap.String("requestID", requestID))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	logger.Info("Users fetched successfully",
		zap.Int("count", len(users)),
		zap.String("requestID", requestID))

	response := gin.H{
		"users": users,
		"page":  page,
		"size":  size,
	}

	logger.Info("Sending response",
		zap.Any("response", response),
		zap.String("requestID", requestID))

	c.JSON(http.StatusOK, response)

	logger.Info("GetUsers handler completed", zap.String("requestID", requestID))
}
