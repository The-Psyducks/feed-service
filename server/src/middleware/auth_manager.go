package middleware

import (
	"fmt"
	"log/slog"
	"strings"

	allErrors "server/src/all_errors"
	"server/src/auth"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			slog.Error("Authorization header is required")
			err := allErrors.AuthenticationErrorHeaderRequired()
			_ = c.AbortWithError(err.Status(), err)
			c.Next()
			return
		}

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 || strings.ToLower(bearerToken[0]) != "bearer" {
			slog.Error("Invalid authorization header")
			err := allErrors.AuthenticationErrorInvalidHeader()
			_ = c.AbortWithError(err.Status(), err)
			c.Next()
			return
		}

		tokenString := bearerToken[1]
		claims, err := auth.ValidateToken(tokenString)
		if err != nil {
			slog.Error("Invalid token")
			err := allErrors.AuthenticationErrorInvalidToken(fmt.Sprintf("%v", err))
			_ = c.AbortWithError(err.Status(), err)
			return
		}

		c.Set("session_user_id", claims.UserId)

		c.Next()
	}
}
