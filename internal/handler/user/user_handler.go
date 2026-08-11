package user

import (
	"net/http"
	"time"

	emailDTO "konsera-backend/internal/DTO/email"
	dto "konsera-backend/internal/DTO/user"
	userService "konsera-backend/internal/services/user"
	"strconv"
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

// @Summary Verify OTP
// @Description Verify the OTP code for user account activation
// @Tags Users
// @Accept json
// @Produce json
// @Param otp body dto.VerifyOTPRequest true "OTP verification information"
// @Success 200 {object} map[string]interface{} "OTP verified successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request or OTP verification failed"
// @Router /auth/verify-otp [post]
func (h *UserHandler) VerifyOTP(c *gin.Context) {
	var req dto.VerifyOTPRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request",
			"error":   err.Error(),
		})
		return
	}

	otp, err := h.service.VerifyOTP(
		c.Request.Context(),
		req.ProfileID,
		req.Code,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	if !otp {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid or expired OTP",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OTP verified successfully. Your account is now active.",
	})
}


// @Summary Verify OTP via Path Parameters
// @Description Verify the OTP code for user account activation using path parameters
// @Tags Users
// @Accept json
// @Produce json
// @Param profile_id path string true "Profile ID"
// @Param code path int true "OTP Code"
// @Success 200 {object} map[string]interface{} "OTP verified successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request or OTP verification failed"
// @Router /auth/verify-otp/{profile_id}/{code} [post]
func (h *UserHandler) VerifyOTPParams(c *gin.Context) {
	profileID := c.Param("profile_id")
	codeStr := c.Param("code")

	if profileID == "" || codeStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Missing profile_id or code in path parameters",
		})
		return
	}

	code, err := strconv.Atoi(codeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid code format. It should be an integer.",
		})
		return
	}

	otp, err := h.service.VerifyOTP(
		c.Request.Context(),
		profileID,
		code,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	if !otp {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid or expired OTP",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OTP verified successfully. Your account is now active.",
	})
}
