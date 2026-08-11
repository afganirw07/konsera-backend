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


// @Summary Create a new user
// @Description Create a new user with the provided information
// @Tags Users
// @Accept json
// @Produce json
// @Param user body dto.CreateUserRequest true "User information"
// @Success 201 {object} dto.UserResponse "User created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request or email/phone already exists"
// @Router /auth/register [post]
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
		c.GetHeader("User-Agent"),
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
