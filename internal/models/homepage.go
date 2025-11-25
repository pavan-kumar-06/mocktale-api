package models

type HomepageSection struct {
	ID          int64         `json:"id"`
	SectionType string        `json:"section_type"`
	Title       string        `json:"title"`
	Subtitle    string        `json:"subtitle,omitempty"`
	SectionData []SectionItem `json:"section_data"`
	DisplayOrder int          `json:"display_order"`
}

type SectionItem struct {
	Slug       string `json:"slug"`
	Name       string `json:"name"`
	PosterURL  string `json:"poster_url"`
	UpdateType string `json:"update_type,omitempty"`  // Only for updates section
	UpdateTitle string `json:"update_title,omitempty"` // Only for updates section
}