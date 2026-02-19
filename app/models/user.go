package models

type User struct {
	BaseModel
	FullName string `json:"full_name" gorm:"not null"`
	Username string `json:"username" gorm:"uniqueIndex;not null"`
	Email    string `json:"email" gorm:"uniqueIndex;not null"`
	Phone    string `json:"phone" gorm:"uniqueIndex;not null"`
	Password string `json:"-" gorm:"not null"`
	Role     string `json:"role" gorm:"default:'user'"` // admin / user
}
