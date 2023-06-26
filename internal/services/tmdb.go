package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"xarr-proxy/internal/api/req"
	"xarr-proxy/internal/api/resp"
	"xarr-proxy/internal/consts"
	"xarr-proxy/internal/db"
	"xarr-proxy/internal/utils"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm/clause"
)

var TMDB = &tmdb{}

type tmdb struct{}

func (t *tmdb) ExternalCheckHealth(url, apiKey string) bool {
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

func (*tmdb) ApiFindByTvdbIDLanguage(id int, language string) (tmdbID int, tmdbTitle string, err error) {
	tmdbUrl := SystemConfig.MustConfigQueryByKey(consts.TMDB_URL)
	tmdbApiKey := SystemConfig.MustConfigQueryByKey(consts.TMDB_API_KEY)
	values := url.Values{}
	values.Add("api_key", tmdbApiKey)
	values.Add("language", language)
	values.Add("external_source", "tvdb_id")
	url := fmt.Sprintf("%s/3/find/%d?%s", tmdbUrl, id, values.Encode())
	log.Debug().Msg(url)
	req, _ := http.NewRequest("GET", url, nil)
	resp, err := utils.GetClient(nil).Do(req)
	if err != nil {
		log.Err(err).Msg("fail to get tmdb by tvdb id")
		return
	}
	defer resp.Body.Close()
	var result TMDBSearchResult
	log.Debug().Msg(utils.DumpResponse(resp))
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Err(err).Msg("fail to decode tmdb search result")
		return
	}
	if len(result.TvResults) > 0 {
		tmdbID = result.TvResults[0].ID
		tmdbTitle = result.TvResults[0].Name
	} else {
		err = errors.New("no result found")
	}
	return
}

// for controller
func (t *tmdb) ApiSync() (bool, error) {
	tvdbIDs := Sonarr.GetNeedSyncTMDBTitleTvdbIDs()
	if len(tvdbIDs) == 0 {
		return true, nil
	}
	language1 := SystemConfig.MustConfigQueryByKey(consts.SONARR_LANGUAGE_1)
	language2 := SystemConfig.MustConfigQueryByKey(consts.SONARR_LANGUAGE_2)
	tmdbTitles := make([]db.TMDBTitle, 0)
	for _, tvdbID := range tvdbIDs {
		for _, language := range []string{language1, language2} {
			tmdbID, tmdbTitlte, err := t.ApiFindByTvdbIDLanguage(tvdbID, language)
			if err != nil {
				log.Err(err).Msgf("fail to get tmdb by tvdb id %d", tvdbID)
				continue
			}
			tmdbTitles = append(tmdbTitles, db.TMDBTitle{
				TvdbID:   tvdbID,
				TmdbID:   tmdbID,
				Title:    tmdbTitlte,
				Language: language,
			})
		}
	}
	if err := db.Get().Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&tmdbTitles).Error; err != nil {
		log.Err(err).Msg("fail to create tmdb titles")
		return false, err
	}
	return true, nil
}

func (*tmdb) ApiBulkDelete(req req.TMDBTitleDeleteReq) (any, error) {
	if err := db.Get().Where("id IN ?", req).Delete(&db.TMDBTitle{}).Error; err != nil {
		return nil, err
	}
	return true, nil
}

func (*tmdb) ApiQuery(req *req.TMDBTitleQueryReq) (any, error) {
	query := db.Get().Model(&db.TMDBTitle{})
	if req.Title != "" {
		query = query.Where("title LIKE ?", "%"+req.Title+"%")
	}
	if req.TvdbID != "" {
		query = query.Where("tvdb_id = ?", req.TvdbID)
	}
	var cnt int64
	tmdbTitles := make([]db.TMDBTitle, 0)
	if err := query.Count(&cnt).Error; err != nil {
		return nil, err
	}
	cleanTitleRegex := SystemConfig.MustConfigQueryByKey(consts.CLEAN_TITLE_REGEX)
	type Row struct {
		db.TMDBTitle
		CleanTitle string `json:"cleanTitle"`
	}
	query = query.Order("update_time DESC").Offset(req.PageSize * (req.Current - 1)).Limit(req.PageSize)
	if err := query.Find(&tmdbTitles).Error; err != nil {
		return nil, err
	}
	rows := make([]Row, 0)
	for i := range tmdbTitles {
		rows = append(rows, Row{
			TMDBTitle:  tmdbTitles[i],
			CleanTitle: utils.CleanTitle(tmdbTitles[i].Title, cleanTitleRegex),
		})
	}
	return resp.RowsResponse[Row]{
		Pagination: resp.Pagination{
			Current:  req.Current,
			PageSize: req.PageSize,
			Total:    int(cnt),
		},
		List: rows,
	}, nil
}
