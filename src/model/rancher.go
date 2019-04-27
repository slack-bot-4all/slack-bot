package model

import "github.com/jinzhu/gorm"

// Rancher : model to w&r on db
type Rancher struct {
	gorm.Model
	Name      string `json:"name" gorm:"unique;not null;type:varchar(50)"`
	URL       string `json:"url" gorm:"not null"`
	AccessKey string `json:"accessKey" gorm:"not null"`
	SecretKey string `json:"secretKey" gorm:"not null"`
}

// TableName : setting the tablename on migrate
func (Rancher) TableName() string {
	return "rancher"
}
