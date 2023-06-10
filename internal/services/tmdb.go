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

var TMDb = &tmdb{}

type tmdb struct{}

func (t *tmdb) CheckHealth(url, apiKey string) bool {
	if url == "" || apiKey == "" || !strings.HasPrefix(url, "http") &&
		!strings.HasPrefix(url, "https") {
		return false
	}
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/3/movie/550?api_key=%s", url, apiKey), nil)
	resp, err := utils.GetClient(nil).Do(req)
	if err != nil {
		log.Err(err).Msg("fail to get health for tmdb")
		return false
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return false
	}
	body := bytes.NewBuffer(make([]byte, 0))
	io.Copy(body, resp.Body)
	return !strings.Contains(body.String(), "Invalid API key")
}
