package global

import (
	"LearningGuide/post_api/config"
	PostProto "LearningGuide/post_api/proto/.PostProto"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	ServerConfig  config.MainConfig
	RDB           redis.Client
	PostSrvClient PostProto.PostClient
)

func InitConfig(path string) {
	v := viper.New()
	v.SetConfigFile(path)
	if err := v.ReadInConfig(); err != nil {
		zap.S().Panicw("Viper Read YAMLFile failed")
	}

	var sc config.MainConfig

	if err := v.Unmarshal(&sc); err != nil {
		zap.S().Panicw("Viper UnMarshal YAMLFile failed")
	}

	ServerConfig = sc
}
