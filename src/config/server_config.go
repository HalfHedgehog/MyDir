package config

type Cfg struct {
	Server  Server  `yaml:"Server"`
	Redis   Redis   `yaml:"Redis"`
	UserRpc UserRpc `yaml:"UserRpc"`
	JWT     JWT     `yaml:"JWT"`
}
type Server struct {
	Port string `yaml:"Port"`
}

type Redis struct {
	Address  string `yaml:"Address"`
	Password string `yaml:"Password"`
	DB       int    `yaml:"Port"`
}

type UserRpc struct {
	Address string `yaml:"Address"`
}

type JWT struct {
	Key         string `yaml:"Key"`
	ExpiresTime int64  `yaml:"ExpiresTime"`
	Issuer      string `yaml:"Issuer"`
}
