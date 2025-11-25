package handlers

import (
	"movie-api/internal/auth"
	"movie-api/internal/database"
	"movie-api/internal/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	// "github.com/google/uuid"
)

// GetMovieRating - Public endpoint for aggregated ratings
func GetMovieRating(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		c.JSON(http.StatusBadRequest, models.MovieResponse{
			Success: false,
			Message: "Slug parameter is required",
		})
		return
	}

	// USE RESPONSE MANAGER INSTEAD OF DB QUERY
	counts := database.ResponseManagerInstance.GetMovieCounts(slug)
	
	rating := models.MovieRating{
		MovieSlug:   slug,
		Option0:     int64(counts.Option0),
		Option1:     int64(counts.Option1), 
		Option2:     int64(counts.Option2),
		Option3:     int64(counts.Option3),
		TotalVotes:  int64(counts.Total),
		UpdatedAt:   time.Now(),
	}

	c.JSON(http.StatusOK, models.MovieResponse{
		Success: true,
		Data:    rating,
	})
}

// GetUserVoteStatus - Check if user has voted for a movie
func GetUserVoteStatus(c *gin.Context) {
	slug := c.Param("slug")
	
	// Get user payload from context (set by middleware)
	userPayload, exists := c.Get("user_payload")
	if !exists {
		c.JSON(http.StatusOK, models.MovieResponse{
			Success: true,
			Data: map[string]interface{}{
				"has_voted":   false,
				"user_choice": nil,
			},
		})
		return
	}

	payload := userPayload.(*auth.TokenPayload)

	// USE RESPONSE MANAGER INSTEAD OF DB QUERY
	hasVoted, userChoice := database.ResponseManagerInstance.HasUserVoted(payload.UserID, slug)

	if hasVoted {
		c.JSON(http.StatusOK, models.MovieResponse{
			Success: true,
			Data: map[string]interface{}{
				"has_voted":   true,
				"user_choice": userChoice,
			},
		})
	} else {
		c.JSON(http.StatusOK, models.MovieResponse{
			Success: true,
			Data: map[string]interface{}{
				"has_voted":   false,
				"user_choice": nil,
			},
		})
	}
}

// SubmitVote - Submit user's vote
func SubmitVote(c *gin.Context) {
	var request struct {
		MovieSlug    string `json:"movie_slug" binding:"required"`
		OptionChosen *int   `json:"option_chosen" binding:"required,min=0,max=3"` // Use pointer to allow 0
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.MovieResponse{
			Success: false,
			Message: "Invalid request data: " + err.Error(),
		})
		return
	}

	// Get user payload from context
	userPayload, exists := c.Get("user_payload")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.MovieResponse{
			Success: false,
			Message: "User token required",
		})
		return
	}
	//for testing without auth
	payload := userPayload.(*auth.TokenPayload)

	// userId := uuid.New().String()

	// Add to in-memory buffer - dereference the pointer
	//for testing without auth
	database.ResponseManagerInstance.AddResponse(payload.UserID, request.MovieSlug, *request.OptionChosen)
	// database.ResponseManagerInstance.AddResponse(userId, request.MovieSlug, *request.OptionChosen)

	c.JSON(http.StatusOK, models.MovieResponse{
		Success: true,
		Message: "Vote submitted successfully",
	})
}