package user

import (
	"net/http"
	"strconv"
	"time"

	emailDTO "konsera-backend/internal/DTO/email"
	dto "konsera-backend/internal/DTO/user"
	helpers "konsera-backend/internal/helpers"
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

// @Summary Login user
// @Description Authenticate a user with email and password, then return a JWT access token
// @Tags Auth
// @Accept json
// @Produce json
// @Param credentials body dto.LoginRequest true "Login credentials"
// @Success 200 {object} dto.LoginResponse "Login successful"
// @Failure 400 {object} map[string]interface{} "Invalid request or credentials"
// @Router /auth/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req dto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.Error(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	result, err := h.service.Login(c.Request.Context(), &req)
	if err != nil {
		helpers.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	helpers.Success(c, http.StatusOK, "Login successful", result)
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
		helpers.Error(c, http.StatusBadRequest, "Invalid request", err.Error())
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
		helpers.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	helpers.Success(c, http.StatusCreated, "User created successfully. Please check your email for OTP verification.", user)
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
		helpers.Error(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	otp, err := h.service.VerifyOTP(
		c.Request.Context(),
		req.ProfileID,
		req.Code,
	)

	if err != nil {
		helpers.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	if !otp {
		helpers.Error(c, http.StatusBadRequest, "Invalid or expired OTP", nil)
		return
	}

	helpers.Success(c, http.StatusOK, "OTP verified successfully. Your account is now active.", nil)
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
		helpers.Error(c, http.StatusBadRequest, "Missing profile_id or code in path parameters", nil)
		return
	}

	code, err := strconv.Atoi(codeStr)
	if err != nil {
		helpers.Error(c, http.StatusBadRequest, "Invalid code format. It should be an integer.", nil)
		return
	}

	otp, err := h.service.VerifyOTP(
		c.Request.Context(),
		profileID,
		code,
	)

	if err != nil {
		helpers.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	if !otp {
		helpers.Error(c, http.StatusBadRequest, "Invalid or expired OTP", nil)
		return
	}

	helpers.Success(c, http.StatusOK, "OTP verified successfully. Your account is now active.", nil)
}

// @Summary Resend OTP
// @Description Resend the OTP code for user account activation
// @Tags Users
// @Accept json
// @Produce json
// @Param otp body dto.ResendOTPRequest true "Resend OTP information"
// @Success 200 {object} map[string]interface{} "OTP resent successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request or failed to resend OTP"
// @Router /auth/resend-otp [post]
func (h *UserHandler) ResendOTP(c *gin.Context) {
	var req dto.ResendOTPRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.Error(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	emailMeta := &emailDTO.RegisterOTPEmailMeta{
		Device:    c.GetHeader("User-Agent"),
		IPAddress: c.ClientIP(),
		Location:  "Indonesia",
		Time:      time.Now().Format("02 Jan 2006, 15:04:05"),
	}

	err := h.service.ResendOTP(
		c.Request.Context(),
		req.ProfileID,
		emailMeta,
	)

	if err != nil {
		helpers.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	helpers.Success(c, http.StatusOK, "OTP resent successfully. Please check your email.", nil)
}
