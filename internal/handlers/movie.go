package handlers

import (
	"database/sql"
	"movie-api/internal/database"
	"movie-api/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetMovies(c *gin.Context) {
	// Get query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	isShowStr := c.Query("is_show")
	
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	// Build query
	query := `
		SELECT 
			 name, slug, image_url, banner_url, year, 
			description, duration_formatted, age_rating_formatted,
			release_date, is_released, is_family_friendly, is_show,
			trailer_video_id, count_watched, number_of_seasons,
			countries, languages, genres, categories, awards,
			actors, directors, created_at
		FROM movies 
		WHERE 1=1
	`

	args := []interface{}{}
	
	if isShowStr != "" {
		isShow, err := strconv.ParseBool(isShowStr)
		if err == nil {
			query += " AND is_show = ?"
			args = append(args, isShow)
		}
	}
	
	query += " ORDER BY created_at DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.MovieResponse{
			Success: false,
			Message: "Failed to fetch movies",
		})
		return
	}
	defer rows.Close()

	var movies []models.Movie

	for rows.Next() {
		var movie models.Movie
		var imageURL, bannerURL, description, durationFormatted, ageRatingFormatted, 
			releaseDate, trailerVideoID, awards, countriesStr, languagesStr, genresStr, 
			categoriesStr, actorsStr, directorsStr sql.NullString
		var year sql.NullInt64
		
		err := rows.Scan(
			&movie.Name, &movie.Slug, &imageURL, &bannerURL, 
			&year, &description, &durationFormatted, &ageRatingFormatted,
			&releaseDate, &movie.IsReleased, &movie.IsFamilyFriendly, &movie.IsShow,
			&trailerVideoID, &movie.CountWatched, &movie.NumberOfSeasons,
			&countriesStr, &languagesStr, &genresStr, &categoriesStr, &awards,
			&actorsStr, &directorsStr, &movie.CreatedAt,
		)
		
		if err != nil {
			continue
		}
		
		// Convert sql.Null types to simple types
		movie.ImageURL = models.NullStringToString(imageURL)
		movie.BannerURL = models.NullStringToString(bannerURL)
		movie.Year = models.NullInt64ToPtr(year)
		movie.Description = models.NullStringToString(description)
		movie.DurationFormatted = models.NullStringToString(durationFormatted)
		movie.AgeRatingFormatted = models.NullStringToString(ageRatingFormatted)
		movie.ReleaseDate = models.NullStringToString(releaseDate)
		movie.TrailerVideoID = models.NullStringToString(trailerVideoID)
		movie.Awards = models.NullStringToString(awards)
		
		// Parse JSON arrays of objects
		movie.Countries = models.ParseCountries(countriesStr)
		movie.Languages = models.ParseLanguages(languagesStr)
		movie.Genres = models.ParseGenres(genresStr)
		movie.Categories = models.ParseCategories(categoriesStr)
		movie.Actors = models.ParsePeople(actorsStr)
		movie.Directors = models.ParsePeople(directorsStr)
		
		movies = append(movies, movie)
	}

	// Get total count
	var total int
	countQuery := "SELECT COUNT(*) FROM movies WHERE 1=1"
	countArgs := []interface{}{}
	
	if isShowStr != "" {
		isShow, err := strconv.ParseBool(isShowStr)
		if err == nil {
			countQuery += " AND is_show = ?"
			countArgs = append(countArgs, isShow)
		}
	}
	
	database.DB.QueryRow(countQuery, countArgs...).Scan(&total)

	response := map[string]interface{}{
		"movies": movies,
		"pagination": models.Pagination{
			Page:  page,
			Limit: limit,
			Total: total,
		},
	}

	c.JSON(http.StatusOK, models.MovieResponse{
		Success: true,
		Data:    response,
	})
}

func GetMovieBySlug(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		c.JSON(http.StatusBadRequest, models.MovieResponse{
			Success: false,
			Message: "Slug parameter is required",
		})
		return
	}

	query := `
		SELECT 
			name, slug, image_url, banner_url, year, 
			description, duration_formatted, age_rating_formatted,
			release_date, is_released, is_family_friendly, is_show,
			trailer_video_id, count_watched, number_of_seasons,
			countries, languages, genres, categories, awards,
			actors, directors, created_at
		FROM movies 
		WHERE slug = ?
	`

	var movie models.Movie
	var imageURL, bannerURL, description, durationFormatted, ageRatingFormatted, 
		releaseDate, trailerVideoID, awards, countriesStr, languagesStr, genresStr, 
		categoriesStr, actorsStr, directorsStr sql.NullString
	var year sql.NullInt64
	
	err := database.DB.QueryRow(query, slug).Scan(
		&movie.Name, &movie.Slug, &imageURL, &bannerURL, 
		&year, &description, &durationFormatted, &ageRatingFormatted,
		&releaseDate, &movie.IsReleased, &movie.IsFamilyFriendly, &movie.IsShow,
		&trailerVideoID, &movie.CountWatched, &movie.NumberOfSeasons,
		&countriesStr, &languagesStr, &genresStr, &categoriesStr, &awards,
		&actorsStr, &directorsStr, &movie.CreatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, models.MovieResponse{
				Success: false,
				Message: "Movie not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.MovieResponse{
			Success: false,
			Message: "Failed to fetch movie",
		})
		return
	}
	
	// Convert sql.Null types to simple types
	movie.ImageURL = models.NullStringToString(imageURL)
	movie.BannerURL = models.NullStringToString(bannerURL)
	movie.Year = models.NullInt64ToPtr(year)
	movie.Description = models.NullStringToString(description)
	movie.DurationFormatted = models.NullStringToString(durationFormatted)
	movie.AgeRatingFormatted = models.NullStringToString(ageRatingFormatted)
	movie.ReleaseDate = models.NullStringToString(releaseDate)
	movie.TrailerVideoID = models.NullStringToString(trailerVideoID)
	movie.Awards = models.NullStringToString(awards)
	
	// Parse JSON arrays of objects
	movie.Countries = models.ParseCountries(countriesStr)
	movie.Languages = models.ParseLanguages(languagesStr)
	movie.Genres = models.ParseGenres(genresStr)
	movie.Categories = models.ParseCategories(categoriesStr)
	movie.Actors = models.ParsePeople(actorsStr)
	movie.Directors = models.ParsePeople(directorsStr)

	c.JSON(http.StatusOK, models.MovieResponse{
		Success: true,
		Data:    movie,
	})
}