package global

import (
	"github.com/OuterCyrex/Gorra/GorraSrv"
	"gorm.io/gorm"
)

var (
	ServerConfig GorraSrv.ServerConfig
	DB           *gorm.DB
)
