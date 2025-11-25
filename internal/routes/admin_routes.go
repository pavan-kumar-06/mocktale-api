// routes/admin_routes.go
package routes

import (
	"movie-api/internal/handlers"
	"movie-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupAdminRoutes(router *gin.Engine) {
	admin := router.Group("/admin")
	admin.Use(middleware.AdminAuth())
	{
		// Movies management
		admin.POST("/movies", handlers.AdminCreateMovie)
		admin.PUT("/movies/:slug", handlers.AdminUpdateMovie)
		
		// Homepage management
		admin.GET("/homepage", handlers.AdminGetHomepageSections)        // Get all sections
		admin.PUT("/homepage", handlers.AdminUpdateHomepage)             // Update entire homepage
		admin.POST("/homepage/reset", handlers.AdminResetHomepage)       // Reset to default
	}
}