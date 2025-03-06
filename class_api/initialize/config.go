package initialize

import (
	"LearningGuide/class_api/config"
	"LearningGuide/class_api/global"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func InitConfig(path string) {
	v := viper.New()
	v.SetConfigFile(path)
	if err := v.ReadInConfig(); err != nil {
		zap.S().Panicw("Viper Read YAMLFile failed")
	}

	var mc config.MainConfig

	if err := v.Unmarshal(&mc); err != nil {
		zap.S().Panicw("Viper UnMarshal YAMLFile failed")
	}

	global.ServerConfig = mc
}
