package models

import (
	"database/sql"
	"encoding/json"
	"time"
)

// Country represents the country object in JSON
type Country struct {
	Name  string `json:"name"`
	Slug  string `json:"slug"`
	Image string `json:"image"`
}

// Language represents the language object in JSON  
type Language struct {
	Name  string `json:"name"`
	Slug  string `json:"slug"`
	Image string `json:"image"`
}

// Genre represents the genre object in JSON
type Genre struct {
	Name  string `json:"name"`
	Slug  string `json:"slug"`
	Color string `json:"color"`
}

// Category represents the category object in JSON
type Category struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// Person represents actor/director object in JSON
type Person struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
	Image string `json:"image"`
}

type Movie struct {
	// ID                  int64     `json:"id"`
	Name                string    `json:"name"`
	Slug                string    `json:"slug"`
	ImageURL            string    `json:"image_url,omitempty"`
	BannerURL           string    `json:"banner_url,omitempty"`
	Year                *int64    `json:"year,omitempty"`
	Description         string    `json:"description,omitempty"`
	DurationFormatted   string    `json:"duration_formatted,omitempty"`
	AgeRatingFormatted  string    `json:"age_rating_formatted,omitempty"`
	ReleaseDate         string    `json:"release_date,omitempty"`
	IsReleased          bool      `json:"is_released"`
	IsFamilyFriendly    bool      `json:"is_family_friendly"`
	IsShow              bool      `json:"is_show"`
	TrailerVideoID      string    `json:"trailer_video_id,omitempty"`
	CountWatched        int64     `json:"count_watched"`
	NumberOfSeasons     int64     `json:"number_of_seasons"`
	Countries           []Country `json:"countries"`
	Languages           []Language `json:"languages"`
	Genres              []Genre   `json:"genres"`
	Categories          []Category `json:"categories"`
	Awards              string    `json:"awards,omitempty"`
	Actors              []Person  `json:"actors"`
	Directors           []Person  `json:"directors"`
	CreatedAt           time.Time `json:"created_at"`
}

type MovieRating struct {
	MovieSlug   string    `json:"movie_slug"`
	Option0     int64     `json:"negative_reviews"`
	Option1     int64     `json:"neutral_reviews"`
	Option2     int64     `json:"positive_reviews"`
	Option3     int64     `json:"perfect_reviews"`
	TotalVotes  int64     `json:"total_votes"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type MovieResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Message string      `json:"message,omitempty"`
}

type Pagination struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Total int `json:"total"`
}

// Helper functions to parse JSON arrays of objects
func ParseCountries(jsonStr sql.NullString) []Country {
	if !jsonStr.Valid || jsonStr.String == "" {
		return []Country{}
	}
	
	var result []Country
	if err := json.Unmarshal([]byte(jsonStr.String), &result); err != nil {
		return []Country{}
	}
	return result
}

func ParseLanguages(jsonStr sql.NullString) []Language {
	if !jsonStr.Valid || jsonStr.String == "" {
		return []Language{}
	}
	
	var result []Language
	if err := json.Unmarshal([]byte(jsonStr.String), &result); err != nil {
		return []Language{}
	}
	return result
}

func ParseGenres(jsonStr sql.NullString) []Genre {
	if !jsonStr.Valid || jsonStr.String == "" {
		return []Genre{}
	}
	
	var result []Genre
	if err := json.Unmarshal([]byte(jsonStr.String), &result); err != nil {
		return []Genre{}
	}
	return result
}

func ParseCategories(jsonStr sql.NullString) []Category {
	if !jsonStr.Valid || jsonStr.String == "" {
		return []Category{}
	}
	
	var result []Category
	if err := json.Unmarshal([]byte(jsonStr.String), &result); err != nil {
		return []Category{}
	}
	return result
}

func ParsePeople(jsonStr sql.NullString) []Person {
	if !jsonStr.Valid || jsonStr.String == "" {
		return []Person{}
	}
	
	var result []Person
	if err := json.Unmarshal([]byte(jsonStr.String), &result); err != nil {
		return []Person{}
	}
	return result
}

// Convert sql.NullString to simple string
func NullStringToString(nullStr sql.NullString) string {
	if nullStr.Valid {
		return nullStr.String
	}
	return ""
}

// Convert sql.NullInt64 to pointer (for optional fields)
func NullInt64ToPtr(nullInt sql.NullInt64) *int64 {
	if nullInt.Valid {
		return &nullInt.Int64
	}
	return nil
}