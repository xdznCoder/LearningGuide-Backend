package global

import (
	"LearningGuide/class_srv/config"
	"gorm.io/gorm"
)

var (
	ServerConfig config.ServerConfig
	DB           *gorm.DB
)
