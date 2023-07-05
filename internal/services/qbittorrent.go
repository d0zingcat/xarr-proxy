package services

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	URL "net/url"
	"strings"

	"xarr-proxy/internal/utils"

	"github.com/rs/zerolog/log"
)

var Qbittorrent = &qbittorrent{}

type qbittorrent struct {
	IsLogin bool
	Cookie  string
}

func (q *qbittorrent) Login(url, username, password string) bool {
	if url == "" || !strings.HasPrefix(url, "http") &&
		!strings.HasPrefix(url, "https") {
		return false
	}
	form := URL.Values{}
	form.Add("username", username)
	form.Add("password", password)
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/api/v2/auth/login", url), strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := utils.GetClient(nil).Do(req)
	if err != nil {
		log.Err(err).Msg("fail to login for qbittorrent")
		return false
	}
	defer resp.Body.Close()
	body := bytes.NewBuffer(make([]byte, 0))
	io.Copy(body, resp.Body)
	if resp.StatusCode != http.StatusOK {
		log.Debug().Any("code", resp.StatusCode).Any("body", body.String()).Msg("fail")
		return false
	}
	for key, value := range resp.Header {
		log.Debug().Msg(key + value[0])
		if key == "Set-Cookie" {
			q.Cookie = value[0]
			q.IsLogin = true
			return true
		}
	}
	return false
}

func (q *qbittorrent) Rename(url string) {
	return
}
