package handler

import (
	"auth_service/internal/dto/request"
	"auth_service/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/gogaruda/apperror"
	"github.com/gogaruda/pkg/response"
	"github.com/gogaruda/valigo"
)

type ProfileHandler struct {
	service service.ProfileService
	valid   *valigo.Valigo
}

func NewProfileHandler(s service.ProfileService, v *valigo.Valigo) *ProfileHandler {
	return &ProfileHandler{service: s, valid: v}
}

func (h *ProfileHandler) Setting(c *gin.Context) {
	var req request.ProfileRequest
	userID, _ := c.Get("user_id")
	req.UserID = userID.(string)
	if !h.valid.ValigoJSON(c, &req) {
		return
	}

	if err := h.service.Setting(req); err != nil {
		apperror.HandleHTTPError(c, err)
		return
	}

	response.Created(c, nil, "query ok")
}

func (h *ProfileHandler) Me(c *gin.Context) {
	userID, _ := c.Get("user_id")
	profile, err := h.service.Me(c.Request.Context(), userID.(string))
	if err != nil {
		apperror.HandleHTTPError(c, err)
		return
	}

	response.OK(c, profile, "query ok", nil)
}
