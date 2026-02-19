package models

type LeagueTable struct {
	BaseModel
	TeamName       string `json:"team_name" gorm:"not null"`
	Played         int    `json:"played" gorm:"default:0"`
	Wins           int    `json:"wins" gorm:"default:0"`
	Draws          int    `json:"draws" gorm:"default:0"`
	Losses         int    `json:"losses" gorm:"default:0"`
	GoalsFor       int    `json:"goals_for" gorm:"default:0"`
	GoalsAgainst   int    `json:"goals_against" gorm:"default:0"`
	GoalDifference int    `json:"goal_difference" gorm:"default:0"`
	Points         int    `json:"points" gorm:"default:0"`
	Position       int    `json:"position"`
}

func (l *LeagueTable) Calculate() {
	l.Points = (l.Wins * 3) + l.Draws
	l.GoalDifference = l.GoalsFor - l.GoalsAgainst
}
