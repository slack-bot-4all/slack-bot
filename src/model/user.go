package model

import "github.com/jinzhu/gorm"

// User : system user
type User struct {
	gorm.Model
	Username string `json:"username" gorm:"not null"`
	Password string `json:"password" gorm:"not null"`
}

// TableName : setting the tablename on migrate
func (User) TableName() string {
	return "user"
}
