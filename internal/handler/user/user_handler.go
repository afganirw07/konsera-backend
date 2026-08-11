package user

import (
	"net/http"
	"time"

	emailDTO "konsera-backend/internal/DTO/email"
	dto "konsera-backend/internal/DTO/user"
	userService "konsera-backend/internal/services/user"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service *userService.UserService
}

func NewUserHandler(
	service *userService.UserService,
) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request",
			"error":   err.Error(),
		})
		return
	}

	emailMeta := &emailDTO.RegisterOTPEmailMeta{
		Device:    c.GetHeader("User-Agent"),
		IPAddress: c.ClientIP(),
		Location:  "Indonesia",
		Time:      time.Now().Format("02 Jan 2006, 15:04:05"),
	}

	user, err := h.service.CreateUser(
		c.Request.Context(),
		&req,
		emailMeta,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully. Verification email has been sent.",
		"data":    user,
	})
}
