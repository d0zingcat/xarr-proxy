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

var Radarr = &radarr{}

type radarr struct{}

func (*radarr) CheckHealth(url, apiKey string) bool {
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
	if strings.Contains(body.String(), "Unauthorized") {
		return true
	}
	return resp.StatusCode == http.StatusOK
}

func (*radarr) CheckIndexerFormat(format string) bool {
	if format == "" {
		return false
	}
	return strings.Contains(format, "{title}") &&
		strings.Contains(format, "{year}")
}
