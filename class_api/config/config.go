package config

type MainConfig struct {
	Name     string         `mapstructure:"name" json:"name"`
	Address  string         `mapstructure:"address" json:"address"`
	Port     int            `mapstructure:"port" json:"port"`
	Tags     []string       `mapstructure:"tags" json:"tags"`
	ClassSrv ClassSrvConfig `mapstructure:"classSrv" json:"classSrv"`
	JwtKey   string         `mapstructure:"jwtKey" json:"jwtKey"`
	Redis    RedisConfig    `mapstructure:"redis" json:"redis"`
	Jaeger   JaegerConfig   `mapstructure:"jaeger" json:"jaeger"`
	Consul   ConsulConfig   `mapstructure:"consul" json:"consul"`
}

type JaegerConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

type ClassSrvConfig struct {
	Name string `mapstructure:"name" json:"name"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host" json:"host"`
	Port     int    `mapstructure:"port" json:"port"`
	Password string `mapstructure:"password" json:"password"`
	DB       int    `mapstructure:"db" json:"db"`
}

type ConsulConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}
