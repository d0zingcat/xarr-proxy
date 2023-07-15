package cron

import (
	"time"

	"xarr-proxy/internal/config"
	"xarr-proxy/internal/consts"
	"xarr-proxy/internal/services"

	"github.com/go-co-op/gocron"
)

var s *gocron.Scheduler

func Init(cfg *config.Config) {
	s = gocron.NewScheduler(time.UTC)
	Register(s)
}

func Register(s *gocron.Scheduler) {
	// sync tmdb titles every hour
	s.Cron("0 * * * *").StartImmediately().Do(func() {
		// TODO: remove
		// services.Sonarr.ApiSync()
		// services.TMDB.ApiSync()
	})
	// sync sonarr titles every 15 minutes
	s.Cron("*/15 * * * *").StartImmediately().Do(func() {
		// services.Sonarr.ApiSync()
	})
	// login into qbittorrent/transmission every 30 minutes
	s.Cron("*/30 * * * *").StartImmediately().Do(func() {
		services.Qbittorrent.Login()
		transUrl := services.SystemConfig.MustConfigQueryByKey(consts.TRANSMISSION_URL)
		if transUrl != "" {
			services.Transmission.IsLogin = false
			transUsername := services.SystemConfig.MustConfigQueryByKey(consts.TRANSMISSION_USERNAME)
			transPasswd := services.SystemConfig.MustConfigQueryByKey(consts.TRANSMISSION_PASSWORD)
			_ = transUsername
			_ = transPasswd
		}
	})
	// rename qbittorrent for sonarr
	s.Cron("*/30 * * * *").StartImmediately().Do(func() {
		services.Qbittorrent.Login()
		services.Qbittorrent.Rename()
	})
}

func StartAsync() {
	s.StartAsync()
}
