package main

import (
	"fmt"

	"xarr-proxy/internal/api"
	"xarr-proxy/internal/config"
	"xarr-proxy/internal/consts"
	"xarr-proxy/internal/cron"
	"xarr-proxy/internal/log"
)

func Init() {
	cfg := config.Init()
	log.Init(cfg)

	cron.Init(cfg)
	cron.StartAsync()

	api.Init(cfg)
	api.Start(cfg)
}

func main() {
	fmt.Println(consts.LOGO)

	Init()
}
