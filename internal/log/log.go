package log

import (
	"os"

	"xarr-proxy/internal/config"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Init(cfg *config.Config) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	if cfg.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
}
