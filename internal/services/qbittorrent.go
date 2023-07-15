package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	URL "net/url"
	"regexp"
	"strconv"
	"strings"

	"xarr-proxy/internal/config"
	"xarr-proxy/internal/consts"
	"xarr-proxy/internal/db"
	"xarr-proxy/internal/utils"

	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
)

const (
	FALLBACK_TIME_SECONDS = 3600
)

var Qbittorrent = &qbittorrent{}

type qbittorrent struct {
	Url     string
	IsLogin bool
	Cookie  string
}

func (q *qbittorrent) ExternalLogin(url, username, password string) bool {
	if url == "" || !strings.HasPrefix(url, "http") &&
		!strings.HasPrefix(url, "https") {
		return false
	}
	q.Url = url
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

func (q *qbittorrent) ExternalFiles(hash string) []string {
	form := URL.Values{}
	form.Add("hash", hash)
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/v2/torrents/files", q.Url), nil)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Cookie", q.Cookie)
	resp, err := utils.GetClient(nil).Do(req)
	if err != nil {
		log.Err(err).Msg("fail to rename for qbittorrent")
		return nil
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Error().Msg("status not ok")
		return nil
	}
	fileArray := make([]map[string]any, 0)
	if err := json.NewDecoder(resp.Body).Decode(&fileArray); err != nil {
		log.Err(err).Msg("fail to decode json")
		return nil
	}
	filenames := lo.Map(fileArray, func(item map[string]any, index int) string {
		return item["name"].(string)
	})
	return filenames
}

func (q *qbittorrent) ExternalRename(hash, name string) bool {
	form := URL.Values{}
	form.Add("hash", hash)
	form.Add("name", name)
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/v2/torrents/rename", q.Url), nil)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Cookie", q.Cookie)
	resp, err := utils.GetClient(nil).Do(req)
	if err != nil {
		log.Err(err).Msg("fail to rename for qbittorrent")
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

func (q *qbittorrent) ExternalRenameFile(hash, oldPath, newPath string) bool {
	form := URL.Values{}
	form.Add("hash", hash)
	form.Add("oldPath", oldPath)
	form.Add("newPath", newPath)
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/v2/torrents/renameFile", q.Url), nil)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Cookie", q.Cookie)
	resp, err := utils.GetClient(nil).Do(req)
	if err != nil {
		log.Err(err).Msg("fail to rename file for qbittorrent")
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

func (q *qbittorrent) Login() bool {
	qbUrl := SystemConfig.MustConfigQueryByKey(consts.QBITTORRENT_URL)
	if qbUrl == "" {
		log.Error().Msg("qbittorrent url is empty")
		return false
	}
	Qbittorrent.IsLogin = false
	qbUsername := SystemConfig.MustConfigQueryByKey(consts.QBITTORRENT_USERNAME)
	qbPasswd := SystemConfig.MustConfigQueryByKey(consts.QBITTORRENT_PASSWORD)
	if ok := Qbittorrent.ExternalLogin(qbUrl, qbUsername, qbPasswd); ok {
		log.Info().Msg("qbittorrent login success")
		return true
	} else {
		log.Error().Msg("qbittorrent login failed")
		return false
	}
}

func (q *qbittorrent) Rename() {
	log.Info().Msg("开始执行 Sonarr 重命名任务")
	if !q.IsLogin {
		log.Error().Msg("Qbittorrent 未登录")
		return
	}

	history, err := Sonarr.ExternalQueryHistory(-FALLBACK_TIME_SECONDS)
	if err != nil {
		log.Info().Msgf("Error fetching data from Sonarr API: %v\n", err)
		return
	}
	if history == nil {
		log.Info().Msg("Sonarr API returned nil")
		return
	}
	tokenRuleMap := make(map[string][]db.SonarrRule)
	tokenRuleMap["season"] = append(tokenRuleMap["season"], Sonarr.QueryByToken("season")...)
	tokenRuleMap["episode"] = append(tokenRuleMap["episode"], Sonarr.QueryByToken("episode")...)
	torrentInfoHashMap := make(map[string]int)
	for _, h := range *history {
		sourceTitle := h.SourceTitle
		torrentInfoHash := h.DownloadID
		if torrentInfoHash == "" {
			torrentInfoHash = h.Data.TorrentInfoHash
		} else {
			torrentInfoHash = strings.ToLower(torrentInfoHash)
		}
		if torrentInfoHash != "" && torrentInfoHashMap[torrentInfoHash] == 0 {
			torrentInfoHashMap[torrentInfoHash] = 1
		}
		if strings.EqualFold(h.Data.DownloadClient, consts.TRANSMISSION) {
			// TODO: transmission rename
		} else if strings.EqualFold(h.Data.DownloadClient, consts.QBITTORRENT) {
			// rename
			q.ExternalRename(torrentInfoHash, sourceTitle)
			if !config.Get().RenameFile {
				return
			}
			log.Info().Msg("renaming files...")
			subtitleNo := 1
			renamed := false
			newFilenameFormat := Sonarr.Format(sourceTitle, "{season}", tokenRuleMap)
			newFilenameFormat += "{episode}"
			files := q.ExternalFiles(torrentInfoHash)
			if files == nil {
				log.Error().Msg("fail to get files")
				return
			}
			for _, oldFilepath := range files {
				lastIndex := strings.LastIndex(oldFilepath, "/") + 1
				if lastIndex > 0 && strings.EqualFold(oldFilepath[:lastIndex-1], sourceTitle) {
					log.Debug().Msgf("文件已经重命名: %s", oldFilepath)
					renamed = true
					break
				}
				oldFilename := oldFilepath[lastIndex:]
				newFilename := oldFilename
				r := regexp.MustCompile("(?i)" + consts.VIDEO_EXT_REGEX)
				matches := r.FindAllStringSubmatch(oldFilename, -1)
				for _, match := range matches {
					extension := match[1]
					newFilename = r.ReplaceAllString(newFilename, "")
					newFilename = strings.Trim(Sonarr.Format(newFilename, newFilenameFormat, tokenRuleMap), " ")
					if !regexp.
						MustCompile(`S\d+E\d+`).
						MatchString(newFilename) || newFilename == "" {
						newFilename = oldFilename
					} else {
						if regexp.
							MustCompile(`(?i)` + consts.SUBTITLE_EXT_REGEX).
							MatchString(extension) {
							newFilename += "." + strconv.Itoa(subtitleNo)
							subtitleNo++
						}
						newFilename += extension
					}
				}
				newFilepath := ""
				i := strings.Index(oldFilepath, "/")
				if i+1 == lastIndex {
					newFilepath = sourceTitle + "/" + newFilename
				} else {
					newFilepath = sourceTitle + oldFilepath[i:lastIndex] + newFilename
				}
				q.ExternalRenameFile(torrentInfoHash, oldFilepath, newFilepath)
				log.Debug().Msgf("qBittorrent rename success: %v => %v", oldFilename, newFilename)
			}
			if !renamed {
				log.Debug().Msgf("qBittorrent torrent rename success: %v => %v", torrentInfoHash, sourceTitle)
			}

		} else {
			log.Error().Msgf("download client %s not supported", h.Data.DownloadClient)
		}
	}
}
