package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"xarr-proxy/internal/api/req"
	"xarr-proxy/internal/api/resp"
	"xarr-proxy/internal/consts"
	"xarr-proxy/internal/db"
	"xarr-proxy/internal/utils"

	"github.com/jinzhu/copier"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	"gorm.io/gorm/clause"
)

var Sonarr = &sonarr{}

type sonarr struct{}

func (*sonarr) QueryByToken(token string) []db.SonarrRule {
	rules := make([]db.SonarrRule, 0)
	if err := db.Get().Where("token = ? AND valid_status = ?", token, consts.VALID_STATUS).Order("priority").Find(&rules).Error; err != nil {
		log.Err(err).Msg("fail to query sonarr rules")
	}
	return rules
}

func (*sonarr) ExternalQueryHistory(fallbackTimeSeconds int) (*SonarrHistory, error) {
	sonarrUrl := SystemConfig.MustConfigQueryByKey(consts.SONARR_URL)
	sonarrApiKey := SystemConfig.MustConfigQueryByKey(consts.SONARR_API_KEY)
	queryTime := time.Now().UTC().Add(time.Second * time.Duration(fallbackTimeSeconds)).Format("2006-01-02T15:04:05Z")

	values := url.Values{}
	values.Add("apikey", sonarrApiKey)
	values.Add("eventType", "1")
	values.Add("date", queryTime)

	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/api/v3/history/since?%s", sonarrUrl, values.Encode()), nil)
	resp, err := utils.GetClient(nil).Do(req)
	if err != nil {
		log.Err(err).Msg("fail to get series for sonarr")
		return nil, err
	}
	defer resp.Body.Close()
	log.Debug().Msg(utils.DumpRequest(*req))
	log.Debug().Msg(utils.DumpResponse(resp))
	history := make(SonarrHistory, 0)
	if err := json.NewDecoder(resp.Body).Decode(&history); err != nil {
		return nil, err
	}
	return &history, nil
}

func (*sonarr) ExternalCheckHealth(url, apiKey string) error {
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/api/v3/health?apikey=%s", url, apiKey), nil)
	resp, err := utils.GetClient(nil).Do(req)
	if err != nil {
		log.Err(err).Msg("fail to get health for sonarr")
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusUnauthorized {
		return err
	}
	body := bytes.NewBuffer(make([]byte, 0))
	io.Copy(body, resp.Body)
	if !strings.Contains(body.String(), "Unauthorized") {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New("not ok")
	}
	return nil
}

func (*sonarr) ExternalGetSeries() (SonarrSeries, error) {
	sonarrUrl := SystemConfig.MustConfigQueryByKey(consts.SONARR_URL)
	sonarrApiKey := SystemConfig.MustConfigQueryByKey(consts.SONARR_API_KEY)
	values := url.Values{}
	values.Add("apikey", sonarrApiKey)
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/api/v3/series?%s", sonarrUrl, values.Encode()), nil)
	resp, err := utils.GetClient(nil).Do(req)
	if err != nil {
		log.Err(err).Msg("fail to get series for sonarr")
		return nil, err
	}
	defer resp.Body.Close()
	log.Debug().Msg(utils.DumpResponse(resp))
	series := make(SonarrSeries, 0)
	if err := json.NewDecoder(resp.Body).Decode(&series); err != nil {
		return nil, err
	}
	return series, err
}

func (*sonarr) ValidRuleQueryByPriority(token string) ([]db.SonarrRule, error) {
	rules := make([]db.SonarrRule, 0)
	if err := db.Get().Where("token = ? AND valid_status = ?", token, consts.VALID_STATUS).Order("priority").Find(&rules).Error; err != nil {
		return nil, err
	}
	return rules, nil
}

func (s *sonarr) GenerateSonarrTitleID(id, sno int) int {
	idS := fmt.Sprintf("%d%d", id, sno)
	id, err := strconv.Atoi(idS)
	if err != nil {
		log.Err(err).Msg("fail to generate sonarr title id")
	}
	return id
}

func (s *sonarr) CheckHealth(url, apiKey string) error {
	err := s.ExternalCheckHealth(url, apiKey)
	return err
}

func (*sonarr) CheckIndexerFormat(format string) bool {
	if format == "" {
		return false
	}
	return strings.Contains(format, "{title}") &&
		strings.Contains(format, "{season}") &&
		strings.Contains(format, "{episode}")
}

func (*sonarr) GetNeedSyncTMDBTitleTvdbIDs() []int {
	tvdbIDs := make([]int, 0)
	if err := db.Get().
		Select("st.tvdb_id").
		Table("sonarr_title AS st").
		Joins("LEFT JOIN tmdb_title AS tt ON st.tvdb_id = tt.tvdb_id").
		Where("st.sno = 0 AND tt.tvdb_id IS NULL").
		Scan(&tvdbIDs).Error; err != nil {
		log.Err(err).Msg("fail to get need sync tvdb ids")
		return nil
	}
	return tvdbIDs
}

func (*sonarr) QueryByTitle(title string) (*db.SonarrTitle, error) {
	sonarrTitle := &db.SonarrTitle{Title: title}
	cleanTitleRegex := SystemConfig.MustConfigQueryByKey(consts.CLEAN_TITLE_REGEX)
	cleanTitle := utils.CleanTitle(title, cleanTitleRegex)
	sonarrTitleList := make([]*db.SonarrTitle, 0)
	fn := func(title string, sonarrTitleList *[]*db.SonarrTitle) error {
		if err := db.Get().Where("clean_title = ?", cleanTitle).Find(sonarrTitleList).Error; err != nil {
			return err
		}
		return nil
	}
	if err := fn(title, &sonarrTitleList); err != nil {
		return sonarrTitle, err
	}
	if len(sonarrTitleList) > 0 {
		sonarrTitle = sonarrTitleList[0]
	} else {
		title = utils.RemoveSeason(title)
		cleanTitle = utils.CleanTitle(title, cleanTitleRegex)
		if err := fn(title, &sonarrTitleList); err != nil {
			return sonarrTitle, err
		}
		if len(sonarrTitleList) > 0 {
			sonarrTitle = sonarrTitleList[0]
		}
	}
	if sonarrTitle != nil {
		sonarrTitle.Title = title
	}
	return sonarrTitle, nil
}

func (*sonarr) QueryAll() ([]db.SonarrTitle, error) {
	subquery1 := db.Get().
		Model(&db.SonarrTitle{}).
		Select("main_title, title, clean_title, season_number, monitored").Group("clean_title")
	subquery2 := db.Get().
		Table("sonarr_title AS st").
		Select("st.main_title, tt.title, NULL AS clean_title, -1 AS season_number, st.monitored").
		Joins("LEFT JOIN tmdb_title AS tt ON st.tvdb_id = tt.tvdb_id").
		Where("st.sno = ?", 0).
		Group("tt.title")
	sonarrTitleList := make([]db.SonarrTitle, 0)
	if err := db.Get().
		Raw("? UNION ?", subquery1, subquery2).
		Select("main_title, title, clean_title, season_number").
		Where("tital IS NOT NULL").
		Order("monitored DESC, LENGTH(title) DESC").
		Scan(&sonarrTitleList).Error; err != nil {
		return nil, err
	}
	return sonarrTitleList, nil
}

func (*sonarr) FormatTitle(text string, format string, cleanTitleRegex string, sonarrRuleList []db.SonarrRule, sonarrTitleList []db.SonarrTitle) string {
	cleanTitleKey := "{cleanTitle}"
	placeholder := `\s`
	for _, sonarrRule := range sonarrRuleList {
		if strings.Contains(sonarrRule.Regex, cleanTitleKey) {
			for _, sonarrTitle := range sonarrTitleList {
				cleanTitle := sonarrTitle.CleanTitle
				if cleanTitle == "" {
					cleanTitle = utils.CleanTitle(sonarrTitle.Title, cleanTitleRegex)
				}
				cleanTitle = regexp.MustCompile(placeholder).ReplaceAllString(cleanTitle, ".?")
				cleanText := utils.CleanTitle(text, cleanTitleRegex)
				regex := strings.ReplaceAll(sonarrRule.Regex, cleanTitleKey, cleanTitle)
				if matched, _ := regexp.MatchString(regex, cleanText); matched {
					if matched, _ := regexp.MatchString("["+placeholder+".?a-zA-Z]+", cleanTitle); matched {
						prefixRegex := strings.ReplaceAll(sonarrRule.Regex, cleanTitleKey, "[a-zA-Z]+"+placeholder+cleanTitle)
						suffixRegex := strings.ReplaceAll(sonarrRule.Regex, cleanTitleKey, cleanTitle+placeholder+"[a-zA-Z]+")
						cleanText = strings.NewReplacer("( season \\d+| episode \\d+| ep \\d+| aka )", "/").Replace(cleanText)
						matched1, _ := regexp.MatchString(prefixRegex, cleanText)
						matched2, _ := regexp.MatchString(suffixRegex, cleanText)
						if matched1 || matched2 {
							log.Printf("英文标题前或后有英文单词：%s，不匹配：%s\n", cleanText, cleanTitle)
							continue
						}
					}
					format = utils.ReplaceToken("title", sonarrTitle.MainTitle, format)
					seasonNumber := sonarrTitle.SeasonNumber
					if seasonNumber != -1 && seasonNumber != 1 {
						format = utils.ReplaceToken("season", "S"+strconv.Itoa(seasonNumber), format)
					}
					break
				}
			}
		} else {
			matcher := regexp.MustCompile(sonarrRule.Regex)
			if matcher.MatchString(text) {
				value := matcher.ReplaceAllString(text, sonarrRule.Replacement)
				format = utils.ReplaceToken("title", value, format)
				break
			}
		}
	}
	return strings.TrimSpace(format)
}

func (*sonarr) Format(text string, format string, tokenRuleMap map[string][]db.SonarrRule) string {
	tokenRegex := regexp.MustCompile(`\{([^}]+)\}`)
	matches := tokenRegex.FindAllStringSubmatch(format, -1)
	for _, match := range matches {
		token := match[1]
		sonarrRuleList := tokenRuleMap[token]
		for _, sonarrRule := range sonarrRuleList {
			tokenMatcher := regexp.MustCompile(sonarrRule.Regex)
			if tokenMatcher.MatchString(text) {
				value := tokenMatcher.ReplaceAllString(text, sonarrRule.Replacement)
				format = utils.ReplaceTokenOffset(token, value, format, &sonarrRule.Offset)
				break
			}
		}
	}

	if strings.Contains(format, "{episode}") {
		return text
	}

	format = utils.RemoveAllToken(format)
	return strings.TrimSpace(format)
}

// for controller
func (s *sonarr) ApiSync() (bool, error) {
	series, err := s.ExternalGetSeries()
	if err != nil {
		return false, err
	}

	cleanTitleRegex := SystemConfig.MustConfigQueryByKey(consts.CLEAN_TITLE_REGEX)
	if err != nil {
		return false, err
	}
	titles := make([]db.SonarrTitle, 0)
	for _, sv := range series {
		sno := 0
		monitored := 0
		if sv.Monitored {
			monitored = 1
		}
		mainT := db.SonarrTitle{
			ID:           s.GenerateSonarrTitleID(sv.ID, sno),
			SeriesID:     sv.ID,
			TvdbID:       sv.TvdbID,
			Sno:          sno,
			MainTitle:    sv.Title,
			Title:        sv.Title,
			CleanTitle:   utils.CleanTitle(sv.Title, cleanTitleRegex),
			SeasonNumber: -1,
			Monitored:    monitored,
		}
		sno++
		titles = append(titles, mainT)

		slugT := db.SonarrTitle{}
		err := copier.Copy(&slugT, &mainT)
		if err != nil {
			log.Err(err).Msg("fail to copy sonarr title")
		}
		slugT.ID = s.GenerateSonarrTitleID(sv.ID, sno)
		slugT.CleanTitle = utils.CleanTitle(sv.TitleSlug, cleanTitleRegex)
		slugT.Sno = sno
		sno++
		titles = append(titles, slugT)

		for _, alias := range sv.AlternateTitles {
			aliasT := db.SonarrTitle{
				ID:           s.GenerateSonarrTitleID(sv.ID, sno),
				SeriesID:     sv.ID,
				TvdbID:       sv.TvdbID,
				Sno:          sno,
				MainTitle:    sv.Title,
				Title:        alias.Title,
				CleanTitle:   utils.CleanTitle(alias.Title, cleanTitleRegex),
				SeasonNumber: alias.SceneSeasonNumber,
				Monitored:    monitored,
			}
			sno++
			titles = append(titles, aliasT)
		}
	}
	if err := db.Get().Table("sonarr_title").Clauses(clause.OnConflict{
		UpdateAll: true,
	}).CreateInBatches(titles, 100).Error; err != nil {
		log.Err(err).Msg("fail to batch create sonarr titles")
	}
	return true, nil
}

func (*sonarr) ApiQuery(req *req.SonarrTitleQueryReq) (any, error) {
	query := db.Get().Model(&db.SonarrTitle{})
	if req.Title != "" {
		query = query.Where("title LIKE ?", "%"+req.Title+"%")
	}
	if req.TvdbID != "" {
		query = query.Where("tvdb_id = ?", req.TvdbID)
	}
	var cnt int64
	sonarrTitles := make([]db.SonarrTitle, 0)
	if err := query.Count(&cnt).Error; err != nil {
		return nil, err
	}
	query = query.Order("update_time DESC").Offset(req.PageSize * (req.Current - 1)).Limit(req.PageSize)
	if err := query.Find(&sonarrTitles).Error; err != nil {
		return nil, err
	}
	return resp.RowsResponse[db.SonarrTitle]{
		Pagination: resp.Pagination{
			Current:  req.Current,
			PageSize: req.PageSize,
			Total:    int(cnt),
		},
		List: sonarrTitles,
	}, nil
}

func (*sonarr) ApiBulkDelete(req req.SonarrTitleDeleteReq) (any, error) {
	if err := db.Get().Where("id IN ?", req).Delete(&db.SonarrTitle{}).Error; err != nil {
		return nil, err
	}
	return true, nil
}

func (*sonarr) ApiRuleSync() (any, error) {
	ruleSyncAuthors := SystemConfig.MustConfigQueryByKey(consts.RULE_SYNC_AUTHORS)
	var authorList []string
	if ruleSyncAuthors == "" {
		return false, errors.New("rule sync author list is empty")
	}
	if ruleSyncAuthors == "ALL" {
		authorList = SystemConfig.ApiAuthorList()
	} else {
		authorList = strings.Split(ruleSyncAuthors, ",")
	}
	authorMap := lo.Reduce(authorList, func(acc map[string]bool, item string, index int) map[string]bool {
		acc[item] = true
		return acc
	}, map[string]bool{})
	log.Info().Msgf("rule sync author list: %v", authorList)
	rules := make([]db.SonarrRule, 0)
	sonarrRules := SystemConfig.ApiGetSonarrRules()
	for _, rule := range sonarrRules {
		// author is disabled
		if _, ok := authorMap[rule.Author]; !ok {
			continue
		}
		rules = append(rules, db.SonarrRule{
			ID:          rule.ID,
			Token:       rule.Token,
			Priority:    rule.Priority,
			Regex:       rule.Regex,
			Replacement: rule.Replacement,
			Offset:      rule.Offset,
			Example:     rule.Example,
			Remark:      rule.Remark,
			Author:      rule.Author,
			ValidStatus: consts.VALID_STATUS,
		})
	}
	if err := db.Get().Clauses(clause.OnConflict{
		UpdateAll: true,
	}).CreateInBatches(rules, 50).Error; err != nil {
		return false, err
	}
	return true, nil
}

func (*sonarr) ApiRuleQuery(req *req.SonarrRuleQueryReq) (any, error) {
	query := db.Get().Model(&db.SonarrRule{})
	if req.Token != "" {
		query = query.Where("token LIKE ?", "%"+req.Token+"%")
	}
	if req.Remark != "" {
		query = query.Where("remark LIKE ?", "%"+req.Remark+"%")
	}
	var cnt int64
	if err := query.Count(&cnt).Error; err != nil {
		return nil, err
	}
	rules := make([]db.SonarrRule, 0)
	query = query.Order("update_time DESC").Offset(req.PageSize * (req.Current - 1)).Limit(req.PageSize)
	if err := query.Find(&rules).Error; err != nil {
		return nil, err
	}
	return resp.RowsResponse[db.SonarrRule]{
		Pagination: resp.Pagination{
			Total:    int(cnt),
			PageSize: req.PageSize,
			Current:  req.Current,
		},
		List: rules,
	}, nil
}

func (*sonarr) ApiSwitchValidStatus(req req.SonarrRuleSwitchValidStatusReq, status int) (any, error) {
	if err := db.Get().Model(&db.SonarrRule{}).Where("id IN ?", req).Update("valid_status", status).Error; err != nil {
		return false, err
	}
	return true, nil
}
