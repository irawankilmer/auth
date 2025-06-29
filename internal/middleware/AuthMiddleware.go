package middleware

import (
	"auth_service/internal/config"
	"auth_service/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/gogaruda/pkg/response"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"strings"
	"time"
)

func AuthMiddleware() gin.HandlerFunc {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		panic("JWT_SECRET tidak ditemukan di environment variable")
	}

	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "Header Authorized tidak ditemukan")
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			response.Unauthorized(c, "Format Authorization harus: bearer {token}")
			return
		}

		tokenStr := parts[1]
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			response.Unauthorized(c, "Token tidak valid atau sudah kadaluarsa!")
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			response.Unauthorized(c, "Klaim token tidak valid!")
			return
		}

		exp, ok := claims["exp"].(float64)
		if !ok || int64(exp) < time.Now().Unix() {
			response.Unauthorized(c, "Token sudah kadaluarsa!")
			return
		}

		userID, ok := claims["user_id"].(string)
		if !ok {
			response.Unauthorized(c, "Token tidak memiliki identitas pengguna (user_id)")
			return
		}

		tokenVersion, ok := claims["token_version"].(string)
		if !ok {
			response.Unauthorized(c, "Token version tidak valid!")
			return
		}

		var user model.UserModel
		db := config.ConnectDB()
		if err := db.QueryRow("SELECT token_version FROM users WHERE id = ?", userID).Scan(&user.TokenVersion); err != nil {
			response.Unauthorized(c, "User tidak ditemukan!")
			return
		}

		if user.TokenVersion != tokenVersion {
			response.Unauthorized(c, "Token sudah tidak berlaku, silahkan login lagi!")
			return
		}

		rolesInterface, ok := claims["roles"].([]interface{})
		if !ok {
			response.Unauthorized(c, "Format roles tidak valid dalam token")
			return
		}

		roles := make([]string, 0, len(rolesInterface))
		for _, r := range rolesInterface {
			roleStr, ok := r.(string)
			if !ok {
				response.Unauthorized(c, "Role tidak valid, harus string")
				return
			}
			roles = append(roles, roleStr)
		}

		c.Set("user_id", userID)
		c.Set("roles", roles)
		c.Next()
	}
}
