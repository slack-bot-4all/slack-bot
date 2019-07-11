package model

import "github.com/jinzhu/gorm"

// ContainerCount : model to w&r on db
type ContainerCount struct {
	gorm.Model
	ContainerID string `json:"jsonId" gorm:"unique;not null;type:varchar(50)"`
	Count       uint   `json:"count" gorm:"not null"`
	IsService   bool   `json:"isService" gorm:"not null"`
	ServiceName string `json:"serviceName" gorm:"not null"`
	StackName   string `json:"stackName" gorm:"not null"`
}

// TableName : setting the tablename on migrate
func (ContainerCount) TableName() string {
	return "containerCount"
}
