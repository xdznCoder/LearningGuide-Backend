package global

import (
	"LearningGuide/file_api/config"
	FileProto "LearningGuide/file_api/proto/.FileProto"
	"github.com/go-redis/redis/v8"
)

var (
	ServerConfig  config.MainConfig
	RDB           redis.Client
	FileSrvClient FileProto.FileClient
)
