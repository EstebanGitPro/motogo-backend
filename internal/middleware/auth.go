package middleware

import (
	"log/slog"
	"net/http"

	"github.com/EstebanGitPro/motogo-backend/config"
	"github.com/EstebanGitPro/motogo-backend/internal/platform/jwt"
	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware(cfg config.JWTConfig) gin.HandlerFunc {
	tokenGenerator := jwt.New(cfg.SecretKey)
	
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		
		var tokenString string
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			tokenString = authHeader[7:]
		} else {
			tokenString = authHeader
		}

		if tokenString == "" {
			slog.Warn("Authentication attempt without token",
				slog.String("client_ip", c.ClientIP()),
				slog.String("user_agent", c.GetHeader("User-Agent")))
			
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Token not provided",
			})
			return
		}

		claims, err := tokenGenerator.Validate(tokenString)
		if err != nil {
			slog.Warn("Invalid token validation attempt",
				slog.String("error", err.Error()),
				slog.String("client_ip", c.ClientIP()),
				slog.String("user_agent", c.GetHeader("User-Agent")))
			
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token",
			})
			return
		}


		c.Set("user_id", claims.ID)
		c.Set("user_email", claims.Email)
		
		slog.Debug("User authenticated successfully",
			slog.String("user_id", claims.ID),
			slog.String("user_email", claims.Email))
		
		c.Next()
	}
}