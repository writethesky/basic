package internal

import (
	"github.com/writethesky/utility/config"
)

type configEntity struct {
	MySQL       mysqlEntity  `yaml:"mysql"`
	TokenServer serverEntity `yaml:"token_server"`
	UserServer  serverEntity `yaml:"user_server"`
	Redis       redisEntity  `yaml:"redis"`
}

type mysqlEntity struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Database string `yaml:"database"`
}

type redisEntity struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Database int    `yaml:"database"`
}

type serverEntity struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

var Config *configEntity

func init() {
	Config = new(configEntity)
	config.Parse(Config, "../")

}
