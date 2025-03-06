package initialize

import (
	"LearningGuide/user_api/config"
	"LearningGuide/user_api/global"
	"github.com/spf13/viper"
	"go.uber.org/zap"
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

	global.ServerConfig = sc
}
