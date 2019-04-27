package core

import (
	"github.com/jinzhu/gorm"
)

// is a single variable to connect on DB on all project
var DB *gorm.DB
