package middleware

import (
	"force-learning/internal/pkg/jwt"
	"force-learning/internal/pkg/response"
	"strings"

	"github.com/gin-gonic/gin"
)

func Auth(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "authorization header is required")
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Unauthorized(c, "invalid authorization header format")
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims, err := jwt.ValidateAccessToken(tokenString, jwtSecret)
		if err != nil {
			response.Unauthorized(c, "invalid or expired token")
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID.String())
		c.Set("token", tokenString)
		c.Next()
	}
}

func Admin() gin.HandlerFunc {
	return func(c *gin.Context) {
		isAdmin := c.GetHeader("X-Admin-Token")
		if isAdmin != "admin-secret-token" {
			response.Forbidden(c, "admin access required")
			c.Abort()
			return
		}
		c.Next()
	}
}

func RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}
