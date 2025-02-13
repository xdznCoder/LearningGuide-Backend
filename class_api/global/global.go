package global

import (
	"LearningGuide/class_api/config"
	proto "LearningGuide/class_api/proto/.ClassProto"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-redis/redis/v8"
)

var (
	ServerConfig   config.MainConfig
	RDB            redis.Client
	Trans          ut.Translator
	ClassSrvClient proto.ClassClient
)
