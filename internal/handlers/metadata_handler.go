// handlers/metadata_handler.go
package handlers

import (
	"database/sql"
	"movie-api/internal/database"
	"movie-api/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetAllGenres - Get all genres
func GetAllGenres(c *gin.Context) {
	rows, err := database.DB.Query("SELECT name, slug, color FROM genres ORDER BY name")
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.MovieResponse{
			Success: false,
			Message: "Failed to fetch genres",
		})
		return
	}
	defer rows.Close()

	var genres []models.Genre
	for rows.Next() {
		var genre models.Genre
		err := rows.Scan(&genre.Name, &genre.Slug, &genre.Color)
		if err != nil {
			continue
		}
		genres = append(genres, genre)
	}

	c.JSON(http.StatusOK, models.MovieResponse{
		Success: true,
		Data:    genres,
	})
}

// GetAllCategories - Get all categories
func GetAllCategories(c *gin.Context) {
	rows, err := database.DB.Query("SELECT name, slug FROM categories ORDER BY name")
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.MovieResponse{
			Success: false,
			Message: "Failed to fetch categories",
		})
		return
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var category models.Category
		err := rows.Scan(&category.Name, &category.Slug)
		if err != nil {
			continue
		}
		categories = append(categories, category)
	}

	c.JSON(http.StatusOK, models.MovieResponse{
		Success: true,
		Data:    categories,
	})
}

// GetAllLanguages - Get all languages
func GetAllLanguages(c *gin.Context) {
	rows, err := database.DB.Query("SELECT name, slug, image_url FROM languages ORDER BY name")
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.MovieResponse{
			Success: false,
			Message: "Failed to fetch languages",
		})
		return
	}
	defer rows.Close()

	var languages []models.Language
	for rows.Next() {
		var language models.Language
		err := rows.Scan(&language.Name, &language.Slug, &language.Image)
		if err != nil {
			continue
		}
		languages = append(languages, language)
	}

	c.JSON(http.StatusOK, models.MovieResponse{
		Success: true,
		Data:    languages,
	})
}

// GetAllCountries - Get all countries
func GetAllCountries(c *gin.Context) {
	rows, err := database.DB.Query("SELECT name, slug, image_url FROM countries ORDER BY name")
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.MovieResponse{
			Success: false,
			Message: "Failed to fetch countries",
		})
		return
	}
	defer rows.Close()

	var countries []models.Country
	for rows.Next() {
		var country models.Country
		err := rows.Scan(&country.Name, &country.Slug, &country.Image)
		if err != nil {
			continue
		}
		countries = append(countries, country)
	}

	c.JSON(http.StatusOK, models.MovieResponse{
		Success: true,
		Data:    countries,
	})
}


// SearchPeople - Search people by name with pagination
func SearchPeople(c *gin.Context) {
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
		SELECT name, slug, image_url 
		FROM people 
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
			Message: "Failed to search people",
		})
		return
	}
	defer rows.Close()

	var people []models.Person
	for rows.Next() {
		var person models.Person
		var imageURL sql.NullString
		
		err := rows.Scan(&person.Name, &person.Slug, &imageURL)
		if err != nil {
			continue
		}
		
		person.Image = models.NullStringToString(imageURL)
		people = append(people, person)
	}

	// Get total count for pagination
	var total int
	countQuery := "SELECT COUNT(*) FROM people WHERE name LIKE ?"
	database.DB.QueryRow(countQuery, searchPattern).Scan(&total)

	response := map[string]interface{}{
		"people": people,
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