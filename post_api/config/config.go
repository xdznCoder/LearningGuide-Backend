package config

import "github.com/OuterCyrex/Gorra/GorraAPI"

type MainConfig struct {
	Name    string                `json:"name"`
	Address string                `json:"address"`
	Port    int                   `json:"port"`
	Tags    []string              `json:"tags"`
	SrvList []string              `json:"srvList"`
	JwtKey  string                `json:"jwtKey"`
	Jaeger  JaegerConfig          `json:"jaeger"`
	Redis   RedisConfig           `json:"redis"`
	Consul  GorraAPI.ConsulConfig `json:"consul"`
}

type JaegerConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type RedisConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

func (c MainConfig) GetRegistryInfo() GorraAPI.RegistryInfo {
	return GorraAPI.RegistryInfo{
		Name:    c.Name,
		Address: c.Address,
		Port:    c.Port,
		Tags:    c.Tags,
		Consul:  c.Consul,
	}
}
