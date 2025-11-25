// handlers/admin_handler.go
package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"movie-api/internal/database"
	"movie-api/internal/models"

	"github.com/gin-gonic/gin"
)

// AdminCreateMovie - Create new movie (FIXED STRUCT)
func AdminCreateMovie(c *gin.Context) {
	var request struct {
		Slug               string             `json:"slug" binding:"required"`
		Name               string             `json:"name" binding:"required"`
		ImageURL           string             `json:"image_url"`
		BannerURL          string             `json:"banner_url"`
		Year               int                `json:"year"`
		Description        string             `json:"description"`
		DurationFormatted  string             `json:"duration_formatted"`
		AgeRatingFormatted string             `json:"age_rating_formatted"`
		ReleaseDate        string             `json:"release_date"`
		IsReleased         bool               `json:"is_released"`
		IsFamilyFriendly   bool               `json:"is_family_friendly"`
		IsShow             bool               `json:"is_show"`
		TrailerVideoID     string             `json:"trailer_video_id"`
		Countries          []models.Country   `json:"countries"`
		Languages          []models.Language  `json:"languages"`
		Genres             []models.Genre     `json:"genres"`
		Categories         []models.Category  `json:"categories"`
		Awards             string             `json:"awards"`
		Actors             []models.Person    `json:"actors"`
		Directors          []models.Person    `json:"directors"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.MovieResponse{
			Success: false,
			Message: "Invalid request: " + err.Error(),
		})
		return
	}

	// Convert structs to JSON strings for database storage
	countriesJSON, _ := json.Marshal(request.Countries)
	languagesJSON, _ := json.Marshal(request.Languages)
	genresJSON, _ := json.Marshal(request.Genres)
	categoriesJSON, _ := json.Marshal(request.Categories)
	actorsJSON, _ := json.Marshal(request.Actors)
	directorsJSON, _ := json.Marshal(request.Directors)

	// Insert movie
	_, err := database.DB.Exec(`
		INSERT INTO movies (
			slug, name, image_url, banner_url, year, description,
			duration_formatted, age_rating_formatted, release_date, is_released,
			is_family_friendly, is_show, trailer_video_id,
			countries, languages, genres, categories, awards, actors, directors
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, 
		request.Slug, request.Name, request.ImageURL, request.BannerURL, request.Year, request.Description,
		request.DurationFormatted, request.AgeRatingFormatted, request.ReleaseDate, request.IsReleased,
		request.IsFamilyFriendly, request.IsShow, request.TrailerVideoID,
		string(countriesJSON), string(languagesJSON), string(genresJSON), 
		string(categoriesJSON), request.Awards, string(actorsJSON), string(directorsJSON),
	)

	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			c.JSON(http.StatusConflict, models.MovieResponse{
				Success: false,
				Message: "Movie with this slug already exists",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.MovieResponse{
			Success: false,
			Message: "Failed to create movie: " + err.Error(),
		})
		return
	}

	// Initialize movie responses
	_, _ = database.DB.Exec(`INSERT OR IGNORE INTO movie_responses (movie_slug) VALUES (?)`, request.Slug)

	c.JSON(http.StatusCreated, models.MovieResponse{
		Success: true,
		Message: "Movie created successfully",
		Data:    request,
	})
}

// AdminUpdateMovie - Update existing movie (FIXED STRUCT)
func AdminUpdateMovie(c *gin.Context) {
	slug := c.Param("slug")
	
	var request struct {
		Name               string             `json:"name"`
		ImageURL           string             `json:"image_url"`
		BannerURL          string             `json:"banner_url"`
		Year               int                `json:"year"`
		Description        string             `json:"description"`
		DurationFormatted  string             `json:"duration_formatted"`
		AgeRatingFormatted string             `json:"age_rating_formatted"`
		ReleaseDate        string             `json:"release_date"`
		IsReleased         bool               `json:"is_released"`
		IsFamilyFriendly   bool               `json:"is_family_friendly"`
		IsShow             bool               `json:"is_show"`
		TrailerVideoID     string             `json:"trailer_video_id"`
		Countries          []models.Country   `json:"countries"`
		Languages          []models.Language  `json:"languages"`
		Genres             []models.Genre     `json:"genres"`
		Categories         []models.Category  `json:"categories"`
		Awards             string             `json:"awards"`
		Actors             []models.Person    `json:"actors"`
		Directors          []models.Person    `json:"directors"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.MovieResponse{
			Success: false,
			Message: "Invalid request: " + err.Error(),
		})
		return
	}

	// Convert structs to JSON strings for database storage
	countriesJSON, _ := json.Marshal(request.Countries)
	languagesJSON, _ := json.Marshal(request.Languages)
	genresJSON, _ := json.Marshal(request.Genres)
	categoriesJSON, _ := json.Marshal(request.Categories)
	actorsJSON, _ := json.Marshal(request.Actors)
	directorsJSON, _ := json.Marshal(request.Directors)

	// Update movie
	result, err := database.DB.Exec(`
		UPDATE movies SET
			name = ?, image_url = ?, banner_url = ?, year = ?, description = ?,
			duration_formatted = ?, age_rating_formatted = ?, release_date = ?, is_released = ?,
			is_family_friendly = ?, is_show = ?, trailer_video_id = ?,
			countries = ?, languages = ?, genres = ?, categories = ?, awards = ?, actors = ?, directors = ?
		WHERE slug = ?
	`, 
		request.Name, request.ImageURL, request.BannerURL, request.Year, request.Description,
		request.DurationFormatted, request.AgeRatingFormatted, request.ReleaseDate, request.IsReleased,
		request.IsFamilyFriendly, request.IsShow, request.TrailerVideoID,
		string(countriesJSON), string(languagesJSON), string(genresJSON), 
		string(categoriesJSON), request.Awards, string(actorsJSON), string(directorsJSON),
		slug,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.MovieResponse{
			Success: false,
			Message: "Failed to update movie: " + err.Error(),
		})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, models.MovieResponse{
			Success: false,
			Message: "Movie not found",
		})
		return
	}

	c.JSON(http.StatusOK, models.MovieResponse{
		Success: true,
		Message: "Movie updated successfully",
		Data:    request,
	})
}


// AdminGetHomepageSections - Get all homepage sections for editing
func AdminGetHomepageSections(c *gin.Context) {
	rows, err := database.DB.Query(`
		SELECT id, section_type, title, subtitle, section_data, 
		       display_order, is_active, created_at, updated_at
		FROM homepage_sections 
		ORDER BY display_order ASC
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.MovieResponse{
			Success: false,
			Message: "Failed to fetch homepage sections",
		})
		return
	}
	defer rows.Close()

	var sections []map[string]interface{}
	for rows.Next() {
		var section struct {
			ID           int
			SectionType  string
			Title        string
			Subtitle     sql.NullString
			SectionData  string
			DisplayOrder int
			IsActive     bool
			CreatedAt    string
			UpdatedAt    string
		}
		
		err := rows.Scan(
			&section.ID, &section.SectionType, &section.Title, 
			&section.Subtitle, &section.SectionData, &section.DisplayOrder,
			&section.IsActive, &section.CreatedAt, &section.UpdatedAt,
		)
		
		if err != nil {
			continue
		}
		
		// Parse section data JSON
		var sectionData []map[string]interface{}
		json.Unmarshal([]byte(section.SectionData), &sectionData)

		sections = append(sections, map[string]interface{}{
			"id":            section.ID,
			"section_type":  section.SectionType,
			"title":         section.Title,
			"subtitle":      section.Subtitle.String,
			"section_data":  sectionData,
			"display_order": section.DisplayOrder,
			"is_active":     section.IsActive,
			"created_at":    section.CreatedAt,
			"updated_at":    section.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, models.MovieResponse{
		Success: true,
		Data:    sections,
	})
}

// AdminUpdateHomepage - Update entire homepage (all sections at once)
func AdminUpdateHomepage(c *gin.Context) {
	var request struct {
		Sections []struct {
			ID           int                      `json:"id"`
			SectionType  string                   `json:"section_type"`
			Title        string                   `json:"title"`
			Subtitle     string                   `json:"subtitle"`
			SectionData  []map[string]interface{} `json:"section_data"`
			DisplayOrder int                      `json:"display_order"`
			IsActive     bool                     `json:"is_active"`
		} `json:"sections" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.MovieResponse{
			Success: false,
			Message: "Invalid request: " + err.Error(),
		})
		return
	}

	// Start transaction
	tx, err := database.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.MovieResponse{
			Success: false,
			Message: "Failed to start transaction",
		})
		return
	}

	// Clear existing sections (optional - or update existing)
	_, err = tx.Exec("DELETE FROM homepage_sections")
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, models.MovieResponse{
			Success: false,
			Message: "Failed to clear existing sections: " + err.Error(),
		})
		return
	}

	// Insert all new sections
	for _, section := range request.Sections {
		sectionDataJSON, _ := json.Marshal(section.SectionData)
		
		_, err := tx.Exec(`
			INSERT INTO homepage_sections 
			(id, section_type, title, subtitle, section_data, display_order, is_active)
			VALUES (?, ?, ?, ?, ?, ?, ?)
		`, section.ID, section.SectionType, section.Title, section.Subtitle, 
		   string(sectionDataJSON), section.DisplayOrder, section.IsActive)
		
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, models.MovieResponse{
				Success: false,
				Message: "Failed to insert section: " + err.Error(),
			})
			return
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, models.MovieResponse{
			Success: false,
			Message: "Failed to save homepage: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.MovieResponse{
		Success: true,
		Message: "Homepage updated successfully",
		Data:    request.Sections,
	})
}

// AdminResetHomepage - Reset homepage to default sections
func AdminResetHomepage(c *gin.Context) {
	// Default sections (you can customize these)
	defaultSections := []map[string]interface{}{
		{
			"section_type": "trending_movies",
			"title":        "Trending Movies",
			"subtitle":     "Most watched movies this week",
			"section_data": []map[string]interface{}{},
			"display_order": 1,
			"is_active":    true,
		},
		{
			"section_type": "trending_updates", 
			"title":        "Latest Updates",
			"subtitle":     "New trailers and teasers",
			"section_data": []map[string]interface{}{},
			"display_order": 2,
			"is_active":    true,
		},
		{
			"section_type": "netflix_latest",
			"title":        "New on Netflix", 
			"subtitle":     "Latest additions to Netflix",
			"section_data": []map[string]interface{}{},
			"display_order": 3,
			"is_active":    true,
		},
		{
			"section_type": "prime_latest",
			"title":        "Prime Video Originals",
			"subtitle":     "Exclusive on Amazon Prime", 
			"section_data": []map[string]interface{}{},
			"display_order": 4,
			"is_active":    true,
		},
	}

	// Start transaction
	tx, err := database.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.MovieResponse{
			Success: false,
			Message: "Failed to start transaction",
		})
		return
	}

	// Clear existing sections
	_, err = tx.Exec("DELETE FROM homepage_sections")
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, models.MovieResponse{
			Success: false,
			Message: "Failed to clear sections: " + err.Error(),
		})
		return
	}

	// Insert default sections
	for _, section := range defaultSections {
		sectionDataJSON, _ := json.Marshal(section["section_data"])
		
		_, err := tx.Exec(`
			INSERT INTO homepage_sections 
			(section_type, title, subtitle, section_data, display_order, is_active)
			VALUES (?, ?, ?, ?, ?, ?)
		`, section["section_type"], section["title"], section["subtitle"], 
		   string(sectionDataJSON), section["display_order"], section["is_active"])
		
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, models.MovieResponse{
				Success: false,
				Message: "Failed to insert default section: " + err.Error(),
			})
			return
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, models.MovieResponse{
			Success: false,
			Message: "Failed to reset homepage: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.MovieResponse{
		Success: true,
		Message: "Homepage reset to default successfully",
		Data:    defaultSections,
	})
}