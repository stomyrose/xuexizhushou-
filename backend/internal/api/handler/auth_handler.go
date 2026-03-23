package handler

import (
	"force-learning/internal/model"
	"force-learning/internal/pkg/response"
	"force-learning/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	authResponse, err := h.authService.Register(&req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Created(c, authResponse)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	authResponse, err := h.authService.Login(&req)
	if err != nil {
		response.Unauthorized(c, err.Error())
		return
	}

	response.Success(c, authResponse)
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var req model.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	authResponse, err := h.authService.Refresh(&req)
	if err != nil {
		response.Unauthorized(c, err.Error())
		return
	}

	response.Success(c, authResponse)
}

func (h *AuthHandler) Verify(c *gin.Context) {
	user, err := h.authService.Verify(c.GetString("token"))
	if err != nil {
		response.Unauthorized(c, err.Error())
		return
	}

	response.Success(c, user)
}

func (h *AuthHandler) GetStatus(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		response.Unauthorized(c, "user not found")
		return
	}

	status, err := h.authService.GetStatusByIDString(userID)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, status)
}
