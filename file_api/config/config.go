package config

import "github.com/OuterCyrex/Gorra/GorraAPI"

type MainConfig struct {
	Address   string                `json:"address"`
	Consul    GorraAPI.ConsulConfig `json:"consul"`
	JwtKey    string                `json:"jwtKey"`
	Name      string                `json:"name"`
	Port      int64                 `json:"port"`
	Redis     RedisConfig           `json:"redis"`
	SrvList   []string              `json:"srvList"`
	AliyunOss OssConfig             `json:"aliyun_oss"`
	ChatGLM   ChatGlmConfig         `json:"chat_glm"`
	Jaeger    JaegerConfig          `json:"jaeger"`
	Tags      []string              `json:"tags"`
}

type JaegerConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type RedisConfig struct {
	DB       int64  `json:"db"`
	Host     string `json:"host"`
	Password string `json:"password"`
	Port     int64  `json:"port"`
}

type OssConfig struct {
	Region     string `json:"region"`
	BucketName string `json:"bucket_name"`
	AccessKey  string `json:"access_key"`
	SecretKey  string `json:"secret_key"`
	EndPoint   string `json:"end_point"`
}

type ChatGlmConfig struct {
	AccessKey string `json:"access_key"`
}

func (c MainConfig) GetRegistryInfo() GorraAPI.RegistryInfo {
	return GorraAPI.RegistryInfo{
		Name:    c.Name,
		Address: c.Address,
		Port:    int(c.Port),
		Tags:    c.Tags,
		Consul:  c.Consul,
	}
}
