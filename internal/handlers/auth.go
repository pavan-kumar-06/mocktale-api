package handlers

import (
	"movie-api/internal/auth"
	"movie-api/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetOrCreateToken generates a new token if user doesn't have one
func GetOrCreateToken(c *gin.Context) {
	// Check if already has valid token
	if payload, valid := auth.ValidateToken(c.Request); valid {
		c.JSON(http.StatusOK, models.MovieResponse{
			Success: true,
			Message: "Token already exists",
			Data: map[string]interface{}{
				"user_id": payload.UserID,
			},
		})
		return
	}

	// Generate new token
	_, err := auth.GenerateAndSetToken(c.Writer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.MovieResponse{
			Success: false,
			Message: "Failed to generate token",
		})
		return
	}

	c.JSON(http.StatusOK, models.MovieResponse{
		Success: true,
		Message: "Token generated successfully",
	})
}