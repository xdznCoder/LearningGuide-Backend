package config

import "github.com/OuterCyrex/Gorra/GorraAPI"

type MainConfig struct {
	Name    string                `mapstructure:"name"`
	Address string                `mapstructure:"address"`
	Port    int                   `mapstructure:"port"`
	Tags    []string              `mapstructure:"tags"`
	SrvList []string              `mapstructure:"srvList"`
	JwtKey  string                `mapstructure:"jwtKey"`
	Jaeger  JaegerConfig          `mapstructure:"jaeger"`
	Redis   RedisConfig           `mapstructure:"redis"`
	Consul  GorraAPI.ConsulConfig `mapstructure:"consul"`
}

type JaegerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
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
