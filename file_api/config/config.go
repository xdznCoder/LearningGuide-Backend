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
	AliyunOss OssConfig             `json:"aliyunOss"`
	ChatGLM   ChatGlmConfig         `json:"chatGlm"`
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
	Region            string `json:"region"`
	FileBucketName    string `json:"fileBucketName"`
	PictureBucketName string `json:"pictureBucketName"`
	AccessKey         string `json:"accessKey"`
	SecretKey         string `json:"secretKey"`
	EndPoint          string `json:"endPoint"`
}

type ChatGlmConfig struct {
	AccessKey string `json:"accessKey"`
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
