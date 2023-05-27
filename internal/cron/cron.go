package cron

import (
	"time"

	"xarr-proxy/internal/config"

	"github.com/go-co-op/gocron"
)

var s *gocron.Scheduler

func Init(cfg *config.Config) {
	s = gocron.NewScheduler(time.UTC)
	Register()
}

func Register() {
}

func StartAsync() {
	s.StartAsync()
}
