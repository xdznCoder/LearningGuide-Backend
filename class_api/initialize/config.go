package initialize

import (
	"LearningGuide/class_api/config"
	"LearningGuide/class_api/global"
	"encoding/json"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func InitConfig() {
	nacosConfig := getNacosConfig()
	sc := []constant.ServerConfig{
		{
			IpAddr: nacosConfig.Host,
			Port:   nacosConfig.Port,
		},
	}

	cc := constant.ClientConfig{
		NamespaceId:         nacosConfig.Namespace,
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "tmp/nacos/log",
		CacheDir:            "tmp/nacos/cache",
		LogLevel:            "debug",
	}

	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": sc,
		"clientConfig":  cc,
	})

	if err != nil {
		zap.S().Panicf("Create Nacos Client Failed: %v", err)
		return
	}

	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: nacosConfig.DataId,
		Group:  nacosConfig.Group,
	})
	if err != nil {
		zap.S().Panicf("Get Nacos JSON Failed: %v", err)
		return
	}

	mainConfig := config.MainConfig{}
	err = json.Unmarshal([]byte(content), &mainConfig)
	if err != nil {
		zap.S().Panicf("Unmarshal Nacos JSON Failed: %v", err)
		return
	}

	global.ServerConfig = mainConfig
}

func getNacosConfig() config.NacosConfig {
	YAMLFile := "class_api/config/config.yaml"
	v := viper.New()
	v.SetConfigFile(YAMLFile)
	if err := v.ReadInConfig(); err != nil {
		zap.S().Panicw("Viper Read YAMLFile failed")
	}

	var nacos config.NacosConfig

	if err := v.Unmarshal(&nacos); err != nil {
		zap.S().Panicw("Viper UnMarshal YAMLFile failed")
	}

	return nacos
}
