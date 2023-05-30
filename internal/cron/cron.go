package cron

import (
	"time"

	"xarr-proxy/internal/config"

	"github.com/go-co-op/gocron"
)

var s *gocron.Scheduler

func Init(cfg *config.Config) {
	s = gocron.NewScheduler(time.UTC)
	Register(s)
}

func Register(s *gocron.Scheduler) {
}

func StartAsync() {
	s.StartAsync()
}
