package controller

import (
	"gin-demo/internal/models"
	"gin-demo/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService *service.UserService
}

func NewUserController(userService *service.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

func (c *UserController) GetUsers(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(ctx.DefaultQuery("size", "10"))

	query := models.UserQuery{
		Page: page,
		Size: size,
	}

	users, total, err := c.userService.GetUsers(ctx.Request.Context(), query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"users": users,
		"total": total,
		"page":  page,
		"size":  size,
	})
}

// 其他方法，如 GetUserByID, CreateUser, UpdateUser, DeleteUser 等，也需要类似的修改
