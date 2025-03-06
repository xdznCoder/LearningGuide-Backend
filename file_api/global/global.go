package global

import (
	"LearningGuide/file_api/config"
	FileProto "LearningGuide/file_api/proto/.FileProto"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	ServerConfig  config.MainConfig
	RDB           redis.Client
	FileSrvClient FileProto.FileClient
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
