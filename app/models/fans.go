package models

import "time"

type Fan struct {
	BaseModel
	Name         string    `json:"name" gorm:"not null"`
	Email        string    `json:"email" gorm:"type:varchar(191);uniqueIndex"`
	Tier         string    `json:"tier"`
	JoinDate     time.Time `json:"join_date"`
	Location     string    `json:"location"`
	MembershipID string    `json:"membership_id" gorm:"type:varchar(191);uniqueIndex"`
}
