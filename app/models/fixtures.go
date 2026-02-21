package models

import "time"

type Fixture struct {
	BaseModel
	HomeTeam       string    `json:"home_team" gorm:"not null"`
	AwayTeam       string    `json:"away_team" gorm:"not null"`
	MatchDate      time.Time `json:"match_date" gorm:"not null"`
	Venue          string    `json:"venue"`
	HomeScore      int       `json:"home_score" gorm:"default:0"`
	AwayScore      int       `json:"away_score" gorm:"default:0"`
	Status         string    `json:"status" gorm:"default:'upcoming'"` // upcoming, live, completed
	PreviewImage   string    `json:"preview_image"`                    // Cloudinary URL — pre-match preview photo
	PreviewCaption string    `json:"preview_caption"`                  // e.g. "Webuye Sportif vs AFC Leopards - Friday Night"
	MatchPhotos    string    `json:"match_photos" gorm:"type:text"`    // JSON array of Cloudinary URLs — post-match action shots
}
