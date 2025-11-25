package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"movie-api/internal/database"
	"movie-api/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetHomepageSections(c *gin.Context) {

	// Get all active sections ordered by display order
	query := `
		SELECT id, section_type, title, subtitle, section_data, display_order
		FROM homepage_sections 
		WHERE is_active = 1
		ORDER BY display_order ASC
	`

	rows, err := database.DB.Query(query)
	if err != nil {
		log.Printf("Database query error: %v", err)
		c.JSON(http.StatusInternalServerError, models.MovieResponse{
			Success: false,
			Message: "Failed to fetch homepage sections",
		})
		return
	}
	defer rows.Close()
	var sections []models.HomepageSection

	for rows.Next() {
		var section models.HomepageSection
		var sectionDataStr sql.NullString
		
		err := rows.Scan(
			&section.ID, &section.SectionType, &section.Title, 
			&section.Subtitle, &sectionDataStr, &section.DisplayOrder,
		)
		
		if err != nil {
			log.Printf("Row scan error: %v", err) // Add this line
			continue
		}
		
		// Parse section data JSON
		if sectionDataStr.Valid {
			if err := json.Unmarshal([]byte(sectionDataStr.String), &section.SectionData); err != nil {
				log.Printf("JSON unmarshal error: %v", err) // Add this line
				continue
			}
		}
		
		sections = append(sections, section)
	}

	c.JSON(http.StatusOK, models.MovieResponse{
		Success: true,
		Data:    sections,
	})
}