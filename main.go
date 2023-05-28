package main

import (
	"fmt"

	"xarr-proxy/internal/api"
	"xarr-proxy/internal/config"
	"xarr-proxy/internal/consts"
	"xarr-proxy/internal/cron"
	"xarr-proxy/internal/db"
	"xarr-proxy/internal/log"
)

func Init() {
	cfg := config.Init()
	log.Init(cfg)

	db.Init(cfg)
	db.Migrate(cfg)

	cron.Init(cfg)
	cron.StartAsync()

	api.Init(cfg)
	api.Start(cfg)
}

func main() {
	fmt.Printf(consts.LOGO, consts.VERSION)

	Init()
}
