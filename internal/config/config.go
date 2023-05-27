package config

import (
	"github.com/caarlos0/env/v8"
)

type (
	Config struct {
		Debug      bool   `json:"debug" env:"DEBUG" envDefault:"false"`
		ServerPort int    `json:"server_port" env:"SERVER_PORT" envDefault:"8117"`
		ConfDir    string `json:"config_dir" env:"CONFIG_DIR" envDefault:"tmp/"`
		DbName     string `json:"db_name" env:"DB_NAME" envDefault:"xarr-proxy.db"`
		JWTSecret  string `json:"jwt_secret" env:"jwt_secret" envDefault:"secret"`
	}
)

var (
	cfg *Config

	Debug      bool
	ServerPort int
	ConfDir    string
	DbName     string
	JWTSecret  string
)

func Init() *Config {
	cfg = &Config{}
	if err := env.Parse(cfg); err != nil {
		panic(err)
	}

	// for quick visits
	ConfDir = cfg.ConfDir
	DbName = cfg.DbName
	Debug = cfg.Debug
	ServerPort = cfg.ServerPort
	JWTSecret = cfg.JWTSecret
	return cfg
}

func Get() *Config {
	return cfg
}
