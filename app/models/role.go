package models

type Role struct {
	BaseModel
	Name        string       `json:"name" gorm:"uniqueIndex;not null"`
	Description string       `json:"description"`
	Permissions []Permission `json:"permissions" gorm:"many2many:role_permissions;"`
}

type Permission struct {
	BaseModel
	Name        string `json:"name" gorm:"uniqueIndex;not null"`
	Description string `json:"description"`
}
