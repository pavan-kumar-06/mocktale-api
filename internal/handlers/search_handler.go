// handlers/search_handler.go
package handlers

import (
	"database/sql"
	"movie-api/internal/database"
	"movie-api/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// SearchMovies - Simple search by name (returns only name, year, poster)
func SearchMovies(c *gin.Context) {
	query := c.Query("q")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	if query == "" {
		c.JSON(http.StatusBadRequest, models.MovieResponse{
			Success: false,
			Message: "Search query 'q' is required",
		})
		return
	}

	searchQuery := `
		SELECT name, slug, image_url, year
		FROM movies 
		WHERE name LIKE ?
		ORDER BY 
			CASE 
				WHEN name LIKE ? THEN 1  -- Exact start matches first
				ELSE 2
			END,
			name ASC
		LIMIT ? OFFSET ?
	`

	searchPattern := "%" + query + "%"
	exactStart := query + "%"
	
	rows, err := database.DB.Query(searchQuery, searchPattern, exactStart, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.MovieResponse{
			Success: false,
			Message: "Failed to search movies",
		})
		return
	}
	defer rows.Close()

	type SimpleMovie struct {
		Name     string `json:"name"`
		Slug     string `json:"slug"`
		ImageURL string `json:"image_url"`
		Year     *int64 `json:"year,omitempty"`
	}

	var movies []SimpleMovie
	for rows.Next() {
		var movie SimpleMovie
		var year sql.NullInt64
		var imageURL sql.NullString
		
		err := rows.Scan(&movie.Name, &movie.Slug, &imageURL, &year)
		if err != nil {
			continue
		}
		
		movie.ImageURL = models.NullStringToString(imageURL)
		if year.Valid {
			movie.Year = &year.Int64
		}
		
		movies = append(movies, movie)
	}

	// Get total count
	var total int
	countQuery := "SELECT COUNT(*) FROM movies WHERE name LIKE ?"
	database.DB.QueryRow(countQuery, searchPattern).Scan(&total)

	response := map[string]interface{}{
		"movies": movies,
		"pagination": models.Pagination{
			Page:  page,
			Limit: limit,
			Total: total,
		},
		"query": query,
	}

	c.JSON(http.StatusOK, models.MovieResponse{
		Success: true,
		Data:    response,
	})
}