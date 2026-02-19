package models

import "time"

type Fixture struct {
	BaseModel
	HomeTeam  string    `json:"home_team" gorm:"not null"`
	AwayTeam  string    `json:"away_team" gorm:"not null"`
	MatchDate time.Time `json:"match_date" gorm:"not null"`
	Venue     string    `json:"venue"`
	HomeScore int       `json:"home_score" gorm:"default:0"`
	AwayScore int       `json:"away_score" gorm:"default:0"`
	Status    string    `json:"status" gorm:"default:'upcoming'"` // upcoming, live, completed
}
