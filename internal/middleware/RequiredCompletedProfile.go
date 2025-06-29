package middleware

import (
	"auth_service/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/gogaruda/apperror"
	"github.com/gogaruda/pkg/response"
	"strings"
)

func RequireCompletedProfile(profileService service.ProfileService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDValue, exists := c.Get("user_id")
		if !exists {
			response.Unauthorized(c, "Token tidak memiliki identitas pengguna (user_id)")
			return
		}

		userID, ok := userIDValue.(string)
		if !ok || userID == "" {
			response.Unauthorized(c, "user_id tidak valid")
			return
		}

		profile, err := profileService.Me(c.Request.Context(), userID)
		if err != nil {
			switch {
			case apperror.Is(err, "PROFILE_NOT_FOUND"):
				response.Unauthorized(c, "profil belum lengkap")
				return
			default:
				response.ServerError(c, "gagal memeriksa profil")
				return
			}
		}

		if strings.TrimSpace(profile.Profile.FullName) == "" {
			response.Forbidden(c, "lengkapi nama lengkap terlebih dahulu")
			return
		}

		c.Set("profile", profile)

		c.Next()
	}
}
