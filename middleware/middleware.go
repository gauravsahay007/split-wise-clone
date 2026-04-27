package middleware

import (
	"net/http"
	"strings"

	"github.com/gauravsahay007/split-wise-clone/utils"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized: Missing or invalid token format",
			})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		userID, err := utils.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized: Invalid or expired token",
			})
			c.Abort()
			return
		}

		// 4. Set the ID in context for other handlers to use
		c.Set("current_user_id", userID)

		c.Next()
	}
}
