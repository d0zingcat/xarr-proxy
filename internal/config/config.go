package config

import (
	"github.com/caarlos0/env/v8"
)

type (
	Config struct {
		Debug              bool   `json:"debug" env:"DEBUG" envDefault:"false"`
		ServerPort         int    `json:"server_port" env:"SERVER_PORT" envDefault:"8117"`
		ConfDir            string `json:"config_dir" env:"CONFIG_DIR" envDefault:"/config"`
		DbName             string `json:"db_name" env:"DB_NAME" envDefault:"xarr-proxy.db"`
		JWTSecret          string `json:"jwt_secret" env:"JWT_SECRET" envDefault:"secret"`
		TokenTTL           int    `json:"token_ttl" env:"TOKEN_TTL" envDefault:"3600"`
		TokenBlockTTL      int    `json:"token_block_ttl" env:"TOKEN_BLOCK_TTL" envDefault:"604800"`
		CacheTTL           int    `json:"cache_expire" env:"CACHE_EXPIRE" envDefault:"300"`
		CachePurgeInterval int    `json:"cache_purge_interval" env:"CACHE_PURGE_INTERVAL" envDefault:"600"`
	}
)

var (
	cfg *Config

	Debug              bool
	ServerPort         int
	ConfDir            string
	DbName             string
	JWTSecret          string
	TokenTTL           int
	TokenBlockTTL      int
	CacheTTL           int
	CachePurgeInterval int
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
	TokenTTL = cfg.TokenTTL
	CacheTTL = cfg.CacheTTL
	CachePurgeInterval = cfg.CachePurgeInterval
	TokenBlockTTL = cfg.TokenBlockTTL

	return cfg
}

func Get() *Config {
	return cfg
}
