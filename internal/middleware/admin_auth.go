// middleware/admin_auth.go
package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func AdminAuth() gin.HandlerFunc {
	adminKey := getEnv("ADMIN_API_KEY", "super_secret_admin_key_2024")

	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-Admin-API-Key")
		
		if apiKey == "" || apiKey != adminKey {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Invalid or missing admin API key",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// Add this helper function locally
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}