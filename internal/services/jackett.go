package services

import (
	"net/http"

	"xarr-proxy/internal/utils"

	"github.com/rs/zerolog/log"
)

var Jackett = &jackett{}

type jackett struct{}

func (*jackett) CheckHealth(url string) bool {
	if url == "" {
		return false
	}
	req, _ := http.NewRequest("GET", url, nil)
	resp, err := utils.GetClient(nil).Do(req)
	if err != nil {
		log.Err(err).Msg("fail to get health for sonarr")
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusUnauthorized
}
