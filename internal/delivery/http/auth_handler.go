package http

import (
	"net/http"

	"be-job-portal/internal/delivery/http/dto"
	"be-job-portal/internal/domain"
	"be-job-portal/pkg/utils"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authUsecase domain.AuthUsecase
}

func NewAuthHandler(us domain.AuthUsecase) *AuthHandler {
	return &AuthHandler{
		authUsecase: us,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var input dto.RegisterRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid input", err.Error())
		return
	}

	err := h.authUsecase.Register(c.Request.Context(), input.Email, input.Password, input.Role)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Registration failed", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "User registered successfully", nil)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var input dto.LoginRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid input", err.Error())
		return
	}

	token, err := h.authUsecase.Login(c.Request.Context(), input.Email, input.Password)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Login failed", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Login successful", gin.H{"token": token})
}

func (h *AuthHandler) GoogleCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Code not found", "Authorization code is missing")
		return
	}

	token, err := h.authUsecase.GoogleLogin(c.Request.Context(), code)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Google login failed", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Google login successful", gin.H{"token": token})
}
