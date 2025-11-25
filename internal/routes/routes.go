package routes

import (
	"movie-api/internal/auth"
	"movie-api/internal/handlers"
	"strings"
	"time"

	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
	router := gin.Default()

	allowedOrigins := getEnv("ALLOWED_ORIGINS", "http://localhost:5173,http://localhost:5174,http://127.0.0.1:5173")
	originsList := strings.Split(allowedOrigins, ",")

	// CORS configuration from environment
	router.Use(cors.New(cors.Config{
		AllowOrigins:     originsList,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "X-Admin-API-Key"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Token middleware for user routes
	router.Use(auth.GetOrCreateToken())

	// API routes
	api := router.Group("/api")
	{
		homepage := api.Group("/homepage")
		{
			homepage.GET("/sections", handlers.GetHomepageSections)
		}
		
		// Movie routes
		movies := api.Group("/movies")
		{
			movies.GET("", handlers.GetMovies)
			movies.GET("/:slug", handlers.GetMovieBySlug)
		}

		// Public rating routes
		ratings := api.Group("/ratings")
		{
			ratings.GET("/:slug", handlers.GetMovieRating)
		}

		// User routes (require token)
		user := api.Group("/user")
		{
			user.GET("/vote-status/:slug", handlers.GetUserVoteStatus)
			user.POST("/vote", handlers.SubmitVote)
			user.GET("/token", handlers.GetOrCreateToken)
		}

		api.GET("/search", handlers.SearchMovies)
		api.GET("/genres", handlers.GetAllGenres)
		api.GET("/categories", handlers.GetAllCategories) 
		api.GET("/languages", handlers.GetAllLanguages)
		api.GET("/countries", handlers.GetAllCountries)
		api.GET("/people/search", handlers.SearchPeople) 

		// Health check
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status":  "ok",
				"message": "Server is running",
			})
		})
	}

	return router
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
