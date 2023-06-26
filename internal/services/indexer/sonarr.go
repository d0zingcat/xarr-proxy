package indexer

import (
	"regexp"

	"xarr-proxy/internal/consts"
	"xarr-proxy/internal/services"
)

type sonarrIndexer struct {
	baseIndexer
}

type sonarrProwlarrIndexer struct {
	sonarrIndexer
}

func NewSonarrProwlarrService() *sonarrProwlarrIndexer {
	return &sonarrProwlarrIndexer{}
}

func (s *sonarrProwlarrIndexer) GetIndexerUrl(path string) string {
	url := services.SystemConfig.MustConfigQueryByKey(consts.PROWLARR_URL)
	r := regexp.MustCompile("/(sonarr|radarr)/prowlarr")
	if path != "" {
		path = r.ReplaceAllString(path, "")
	}
	return url + path
}

func (s *sonarrProwlarrIndexer) GetTitle(key string) string {
	return s.baseIndexer.GetTitle(key)
}
