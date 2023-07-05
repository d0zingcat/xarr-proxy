package cron

import (
	"time"

	"xarr-proxy/internal/config"
	"xarr-proxy/internal/consts"
	"xarr-proxy/internal/services"

	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"
)

var s *gocron.Scheduler

func Init(cfg *config.Config) {
	s = gocron.NewScheduler(time.UTC)
	Register(s)
}

func Register(s *gocron.Scheduler) {
	// sync tmdb titles every day
	s.Cron("0 0 * * *").StartImmediately().Do(func() {
		services.Sonarr.ApiSync()
		// TODO: remove
		// services.TMDB.ApiSync()
	})
	// sync sonarr titles every 15 minutes
	s.Cron("*/15 * * * *").StartImmediately().Do(func() {
		services.Sonarr.ApiSync()
	})
	// login into qbittorrent/transmission every 30 minutes
	s.Cron("*/30 * * * *").StartImmediately().Do(func() {
		qbUrl := services.SystemConfig.MustConfigQueryByKey(consts.QBITTORRENT_URL)
		if qbUrl != "" {
			services.Qbittorrent.IsLogin = false
			qbUsername := services.SystemConfig.MustConfigQueryByKey(consts.QBITTORRENT_USERNAME)
			qbPasswd := services.SystemConfig.MustConfigQueryByKey(consts.QBITTORRENT_PASSWORD)
			if ok := services.Qbittorrent.Login(qbUrl, qbUsername, qbPasswd); ok {
				log.Info().Msg("qbittorrent login success")
			} else {
				log.Error().Msg("qbittorrent login failed")
			}
		}
		transUrl := services.SystemConfig.MustConfigQueryByKey(consts.TRANSMISSION_URL)
		if transUrl != "" {
			services.Transmission.IsLogin = false
			transUsername := services.SystemConfig.MustConfigQueryByKey(consts.TRANSMISSION_USERNAME)
			transPasswd := services.SystemConfig.MustConfigQueryByKey(consts.TRANSMISSION_PASSWORD)
			_ = transUsername
			_ = transPasswd
		}
	})
}

func StartAsync() {
	s.StartAsync()
}
