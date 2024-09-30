package controller

import (
	"gin-demo/internal/middleware"
	"gin-demo/internal/models"
	"gin-demo/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserController struct {
	userService *service.UserService
	logger      *zap.Logger
}

func NewUserController(userService *service.UserService, logger *zap.Logger) *UserController {
	return &UserController{
		userService: userService,
		logger:      logger,
	}
}

func (uc *UserController) GetUser(c *gin.Context) {
	ctx := c.Request.Context()
	requestID := middleware.GetRequestID(ctx)

	id := c.Param("id")
	uc.logger.Info("GetUser request",
		zap.String("request_id", requestID),
		zap.String("id", id))

	user, err := uc.userService.GetUserByID(ctx, id)
	if err != nil {
		uc.logger.Error("Failed to get user",
			zap.String("request_id", requestID),
			zap.Error(err),
		)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	uc.logger.Info("User retrieved successfully",
		zap.String("request_id", requestID),
		zap.String("user_id", id),
		zap.Any("user", user))

	response := gin.H{"user": user}
	uc.logger.Info("Sending response",
		zap.String("request_id", requestID),
		zap.Any("response", response))
	c.JSON(http.StatusOK, response)
}

func (uc *UserController) GetUsers(c *gin.Context) {
	ctx := c.Request.Context()
	requestID := middleware.GetRequestID(ctx)

	query := models.UserQuery{
		Name:  c.Query("name"),
		Email: c.Query("email"),
		Page:  1,
		Size:  10,
	}

	if page, err := strconv.Atoi(c.Query("page")); err == nil && page > 0 {
		query.Page = page
	}

	if size, err := strconv.Atoi(c.Query("size")); err == nil && size > 0 {
		query.Size = size
	}

	uc.logger.Info("GetUsers request",
		zap.String("request_id", requestID),
		zap.Any("query", query))

	users, total, err := uc.userService.GetUsers(ctx, query)
	if err != nil {
		uc.logger.Error("Failed to get users",
			zap.String("request_id", requestID),
			zap.Error(err),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get users"})
		return
	}

	uc.logger.Info("Users retrieved successfully",
		zap.String("request_id", requestID),
		zap.Int("total", int(total)),
		zap.Int("retrieved", len(users)))

	response := gin.H{
		"users": users,
		"total": total,
		"page":  query.Page,
		"size":  query.Size,
	}
	uc.logger.Info("Sending response",
		zap.String("request_id", requestID),
		zap.Any("response", response))
	c.JSON(http.StatusOK, response)
}
