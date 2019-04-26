package model

import "github.com/jinzhu/gorm"

// Task : system user
type Task struct {
	gorm.Model
	Service            string `json:"service" gorm:"not null"`
	ChannelToSendAlert string `json:"channelToSendAlert" gorm:"not null"`
	RancherURL         string `json:"rancherUrl" gorm:"not null"`
	RancherAccessKey   string `json:"rancherAccessKey" gorm:"not null"`
	RancherSecretKey   string `json:"rancherSecretKey" gorm:"not null"`
	RancherProjectID   string `json:"rancherProjectId" gorm:"not null"`
}

// TableNane : setting the tablename on migrate
func (Task) TableName() string {
	return "task"
}
