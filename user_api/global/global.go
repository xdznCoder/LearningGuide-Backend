package global

import (
	"LearningGuide/user_api/config"
	proto "LearningGuide/user_api/proto/userProto"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-redis/redis/v8"
)

var (
	ServerConfig  config.MainConfig
	RDB           redis.Client
	Trans         ut.Translator
	UserSrvClient proto.UserClient
)
