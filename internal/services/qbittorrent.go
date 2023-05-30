package services

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"

	"xarr-proxy/internal/utils"

	"github.com/rs/zerolog/log"
)

var Qbittorrent = &qbittorrent{}

type qbittorrent struct {
	IsLogin bool
	Cookie  string
}

func (*qbittorrent) CheckHealth(url, apiKey string) bool {
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/api/v3/health?apikey=%s", url, apiKey), nil)
	resp, err := utils.GetClient(nil).Do(req)
	if err != nil {
		log.Err(err).Msg("fail to get health for sonarr")
		return false
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusUnauthorized {
		return false
	}
	body := bytes.NewBuffer(make([]byte, 0))
	io.Copy(body, resp.Body)
	return !strings.Contains(body.String(), "Unauthorized")
}

func (*qbittorrent) Login(url, username, password string) {
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/api/v2/auth/login?username=%s&password=%s", url, username, password), nil)
	resp, err := utils.GetClient(nil).Do(req)
	if err != nil {
		log.Err(err).Msg("fail to login for qbittorrent")
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusUnauthorized {
		return
	}
	body := bytes.NewBuffer(make([]byte, 0))
	io.Copy(body, resp.Body)
	log.Info().Msg(body.String())
}
