package config

type MainConfig struct {
	Name    string        `mapstructure:"name" json:"name"`
	Address string        `mapstructure:"address" json:"address"`
	Port    int           `mapstructure:"port" json:"port"`
	Tags    []string      `mapstructure:"tags" json:"tags"`
	UserSrv UserSrvConfig `mapstructure:"userSrv" json:"userSrv"`
	JwtKey  string        `mapstructure:"jwtKey" json:"jwtKey"`
	Redis   RedisConfig   `mapstructure:"redis" json:"redis"`
	Consul  ConsulConfig  `mapstructure:"consul" json:"consul"`
	Jaeger  JaegerConfig  `mapstructure:"jaeger" json:"jaeger"`
	Email   EmailConfig   `mapstructure:"email" json:"email"`
}

type UserSrvConfig struct {
	Name string `mapstructure:"name" json:"name"`
}

type JaegerConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
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

type EmailConfig struct {
	Host     string `mapstructure:"host" json:"host"`
	Password string `mapstructure:"password" json:"password"`
}
