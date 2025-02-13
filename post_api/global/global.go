package global

import (
	"LearningGuide/post_api/config"
	PostProto "LearningGuide/post_api/proto/.PostProto"
	UserProto "LearningGuide/post_api/proto/.UserProto"
	"github.com/go-redis/redis/v8"
)

var (
	ServerConfig  config.MainConfig
	RDB           redis.Client
	PostSrvClient PostProto.PostClient
	UserSrvClient UserProto.UserClient
)
