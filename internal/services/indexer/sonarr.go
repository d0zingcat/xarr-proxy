package indexer

import (
	xmlUtil "encoding/xml"
	"regexp"
	"strings"

	"xarr-proxy/internal/consts"
	"xarr-proxy/internal/db"
	"xarr-proxy/internal/services"

	"github.com/rs/zerolog/log"
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

func (s *sonarrProwlarrIndexer) GetSearchTitle(title string) []string {
	searchTitleList := []string{}
	sonarrTitle, _ := services.Sonarr.QueryByTitle(title)
	if sonarrTitle == nil {
		sonarrTitle = syncAndGetSonarrTitle(title)
	}
	searchTitleList = append(searchTitleList, sonarrTitle.Title)

	if sonarrTitle.Sno == 0 || sonarrTitle.Sno == 1 {
		tmdbTitleList := make([]db.TMDBTitle, 0)
		if err := db.Get().Where("tvdb_id = ?", sonarrTitle.TvdbID).Find(&tmdbTitleList).Error; err != nil {
			log.Err(err).Msg("fail to query tmdb titles")
		}
		if len(tmdbTitleList) > 0 {
			for _, tmdbTitle := range tmdbTitleList {
				searchTitleList = append(searchTitleList, tmdbTitle.Title)
			}
			searchTitleList = append(searchTitleList, tmdbTitleList[0].Title)
		}
	}
	return searchTitleList
}

func (s *sonarrProwlarrIndexer) ExecuteFormatRule(xml string) string {
	if xml == "" || !strings.Contains(xml, "<item>") {
		return xml
	}

	format := services.SystemConfig.MustConfigQueryByKey(consts.SONARR_INDEXER_FORMAT)
	tokenRuleMap := make(map[string][]db.SonarrRule)
	matcher := regexp.MustCompile(`\{([^}]+)\}`).FindAllStringSubmatch(format, -1)
	for _, match := range matcher {
		log.Debug().Msgf("%v:%v", match[0], match[1])
		token := match[1]
		tokenRuleMap[token] = services.Sonarr.QueryByToken(token)
	}

	if _, ok := tokenRuleMap["title"]; !ok {
		return xml
	}
	x := &TorznabRss{}
	err := xmlUtil.Unmarshal([]byte(xml), x)
	if err != nil {
		log.Err(err).Msg("fail to parse xml")
		return xml
	}
	cleanTitleRegex := services.SystemConfig.MustConfigQueryByKey(consts.CLEAN_TITLE_REGEX)
	sonarrTitleList, _ := services.Sonarr.QueryAll()
	channel := &x.Channel
	for i, item := range channel.Items {
		formatText := services.Sonarr.FormatTitle(item.Title, format, cleanTitleRegex, tokenRuleMap["title"], sonarrTitleList)
		if strings.Contains(formatText, "{title}") {
			log.Error().Msgf("索引器格式化失败：%s ==> 未匹配到标题", item.Title)
			continue
		}
		formatText = services.Sonarr.Format(item.Title, formatText, tokenRuleMap)
		log.Debug().Msgf("索引器格式化：%s ==> %s", item.Title, formatText)
		channel.Items[i].Title = formatText
	}
	d, err := xmlUtil.MarshalIndent(x, "", "\t")
	if err != nil {
		log.Err(err).Msg("fail to marshal to xml string")
		return xml
	}
	return "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n" + string(d)
}

func (s *sonarrProwlarrIndexer) GetIndexerUrl(path string) string {
	url := services.SystemConfig.MustConfigQueryByKey(consts.PROWLARR_URL)
	r := regexp.MustCompile("/(sonarr|radarr)/prowlarr")
	if path != "" {
		path = r.ReplaceAllString(path, "")
	}
	return url + path
}

func syncAndGetSonarrTitle(title string) *db.SonarrTitle {
	sonarrTitle, _ := services.Sonarr.QueryByTitle(title)
	if sonarrTitle == nil {
		services.Sonarr.ApiSync()
		services.TMDB.ApiSync()
		sonarrTitle, _ := services.Sonarr.QueryByTitle(title)
		if sonarrTitle == nil {
			log.Error().Msgf("找不到匹配的标题：%s\n", title)
			sonarrTitle = &db.SonarrTitle{Title: title}
		}
	}
	return sonarrTitle
}
