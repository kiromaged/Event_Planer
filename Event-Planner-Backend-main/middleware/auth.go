package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"event_planner_backend/config"
	"event_planner_backend/utils"
)

// AuthMiddleware validates JWT tokens and sets user context.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.JSONError(c, http.StatusUnauthorized, "authorization header required")
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.JSONError(c, http.StatusUnauthorized, "invalid authorization header format")
			return
		}

		tokenString := parts[1]
		secret := config.MustGetEnv("JWT_SECRET", "dev_secret_change_me")
		claims, err := utils.ParseAndValidateJWT(secret, tokenString)
		if err != nil {
			utils.JSONError(c, http.StatusUnauthorized, "invalid or expired token")
			return
		}

		// Extract user ID from claims
		sub, ok := claims["sub"]
		if !ok {
			utils.JSONError(c, http.StatusUnauthorized, "invalid token claims")
			return
		}

		var userID uint
		switch v := sub.(type) {
		case float64:
			userID = uint(v)
		case uint:
			userID = v
		case int:
			userID = uint(v)
		case string:
			id, err := strconv.ParseUint(v, 10, 32)
			if err != nil {
				utils.JSONError(c, http.StatusUnauthorized, "invalid user ID in token")
				return
			}
			userID = uint(id)
		default:
			utils.JSONError(c, http.StatusUnauthorized, "invalid user ID type in token")
			return
		}

		// Set user ID in context
		c.Set("userID", userID)
		c.Set("email", claims["email"])

		c.Next()
	}
}

// GetUserID extracts the user ID from the context (set by AuthMiddleware).
func GetUserID(c *gin.Context) (uint, bool) {
	userID, exists := c.Get("userID")
	if !exists {
		return 0, false
	}
	id, ok := userID.(uint)
	return id, ok
}

