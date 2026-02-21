package models

import "time"

type Sponsor struct {
	BaseModel
	Name        string    `json:"name" gorm:"not null"`
	Logo        string    `json:"logo"`
	Website     string    `json:"website"`
	Tier        string    `json:"tier"`
	ContractEnd time.Time `json:"contract_end"`
	Active      bool      `json:"active" gorm:"default:true"`
}
