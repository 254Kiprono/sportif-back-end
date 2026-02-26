package models

import (
	"github.com/google/uuid"
)

type Lineup struct {
	BaseModel
	FixtureID uuid.UUID      `json:"fixture_id" gorm:"type:char(36);uniqueIndex"`
	Formation string         `json:"formation"` // e.g., "4-3-3", "4-4-2"
	Players   []LineupPlayer `json:"players" gorm:"foreignKey:LineupID"`
}

type LineupPlayer struct {
	BaseModel
	LineupID  uuid.UUID `json:"lineup_id" gorm:"type:char(36);index"`
	PlayerID  uint      `json:"player_id" gorm:"index"`
	Player    Player    `json:"player" gorm:"foreignKey:PlayerID"`
	Position  string    `json:"position"` // Specific role on the pitch
	IsStarter bool      `json:"is_starter" gorm:"default:true"`
	IsCaptain bool      `json:"is_captain" gorm:"default:false"`
}

type MatchEventType string

const (
	EventGoal         MatchEventType = "goal"
	EventAssist       MatchEventType = "assist"
	EventYellowCard   MatchEventType = "yellow_card"
	EventRedCard      MatchEventType = "red_card"
	EventSubstitution MatchEventType = "substitution"
)

type MatchEvent struct {
	BaseModel
	FixtureID      uuid.UUID      `json:"fixture_id" gorm:"type:char(36);index"`
	Type           MatchEventType `json:"type"`
	Minute         int            `json:"minute"`
	PlayerID       *uint          `json:"player_id,omitempty" gorm:"index"`
	Player         *Player        `json:"player,omitempty" gorm:"foreignKey:PlayerID"`
	AssistPlayerID *uint          `json:"assist_player_id,omitempty" gorm:"index"`
	AssistPlayer   *Player        `json:"assist_player,omitempty" gorm:"foreignKey:AssistPlayerID"`
	PlayerOutID    *uint          `json:"player_out_id,omitempty" gorm:"index"` // For substitutions
	PlayerOut      *Player        `json:"player_out,omitempty" gorm:"foreignKey:PlayerOutID"`
	IsOpponent     bool           `json:"is_opponent" gorm:"default:false"` // To track opponent goals for live score
	Commentary     string         `json:"commentary"`
}
