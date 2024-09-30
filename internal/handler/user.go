package handler

import (
	"gin-demo/internal/database"
	"gin-demo/internal/logger"
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
	log := logger.FromContext(ctx)

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		log.Error("Failed to parse page parameter", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page parameter"})
		return
	}

	size, err := strconv.Atoi(c.DefaultQuery("size", "10"))
	if err != nil {
		log.Error("Failed to parse size parameter", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid size parameter"})
		return
	}

	log.Info("Fetching users", zap.Int("page", page), zap.Int("size", size))

	var users []User
	result := database.DB.WithContext(ctx).Offset((page - 1) * size).Limit(size).Find(&users)
	if result.Error != nil {
		log.Error("Failed to fetch users from database", zap.Error(result.Error))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	log.Info("Users fetched successfully", zap.Int("count", len(users)))

	response := gin.H{
		"users": users,
		"page":  page,
		"size":  size,
	}

	c.JSON(http.StatusOK, response)

	log.Info("GetUsers handler completed")
}
