package handler

import (
	"auth_service/internal/dto/request"
	"auth_service/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/gogaruda/apperror"
	"github.com/gogaruda/pkg/response"
	"github.com/gogaruda/valigo"
	"net/http"
)

type AuthHandler struct {
	authService service.AuthService
	valid       *valigo.Valigo
}

func NewAuthHandler(a service.AuthService, v *valigo.Valigo) *AuthHandler {
	return &AuthHandler{authService: a, valid: v}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req request.RegisterRequest
	req.Roles = []string{"tamu"}
	if !h.valid.ValigoJSON(c, &req) {
		return
	}

	if err := h.authService.Register(c.Request.Context(), req); err != nil {
		apperror.HandleHTTPError(c, err)
		return
	}

	response.OK(c, nil, "daftar berhasil", nil)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req request.LoginRequest
	if !h.valid.ValigoJSON(c, &req) {
		return
	}

	token, err := h.authService.Login(req)
	if err != nil {
		apperror.HandleHTTPError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":   http.StatusOK,
		"status": "success",
		"token":  token,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	userID, _ := c.Get("user_id")
	if err := h.authService.Logout(userID.(string)); err != nil {
		apperror.HandleHTTPError(c, err)
		return
	}

	response.OK(c, nil, "berhasil logout", nil)
}
