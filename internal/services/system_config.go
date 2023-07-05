package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"xarr-proxy/internal/consts"
	"xarr-proxy/internal/db"
	"xarr-proxy/internal/model"
	"xarr-proxy/internal/utils"

	"github.com/pelletier/go-toml/v2"
	"github.com/samber/lo"

	"github.com/rs/zerolog/log"
)

var SystemConfig = &systemConfig{}

type systemConfig struct{}

func (*systemConfig) ApiVersion() string {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", consts.AUTHOR, consts.REPO)
	version := strings.Replace(consts.VERSION, "v", "", 1)
	req, _ := http.NewRequest("GET", url, nil)
	resp, err := utils.GetClient(nil).Do(req)
	if err != nil {
		log.Err(err).Msg("fail to get latest release info")
		return ""
	}
	defer resp.Body.Close()
	data := make(map[string]any)
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		log.Err(err).Msg("fail to decode response body")
		return ""
	}
	if data["tag_name"] != nil {
		v := data["tag_name"].(string)
		if v != version {
			return version + " 🚨"
		}
	}
	return version
}

func (*systemConfig) ApiAuthorList() []string {
	authorList := []string{"LuckyPuppy514", "d0zingcat"}
	configAuthorList := make(map[string][]string)
	content, err := utils.ReadFile(filepath.Join(consts.RULE_FILE_DIR, "author.toml"))
	if err != nil {
		log.Err(err).Msg("fail to read author.toml")
		return authorList
	}
	err = toml.Unmarshal([]byte(content), &configAuthorList)
	if err != nil {
		log.Err(err).Msg("fail to unmarshal author.toml")
		return authorList
	}
	if len(configAuthorList["author"]) > 0 {
		return configAuthorList["author"]
	}
	return authorList
}

func (*systemConfig) ApiConfigQuery() []model.SystemConfig {
	configs := make([]model.SystemConfig, 0)
	if err := db.Get().Find(&configs, "valid_status IS NOT NULL").Error; err != nil {
		log.Err(err).Msg("fail to query system config")
		return configs
	}
	return configs
}

func (*systemConfig) ApiGetSonarrRules() []model.Rule {
	configRuleList := make(map[string][]model.Rule)
	content, err := utils.ReadFile(filepath.Join(consts.RULE_FILE_DIR, "sonarr.toml"))
	if err != nil {
		log.Err(err).Msg("fail to read sonarr.toml")
		return nil
	}
	err = toml.Unmarshal([]byte(content), &configRuleList)
	if err != nil {
		log.Err(err).Msg("fail to unmarshal sonarr.toml")
	}
	return configRuleList["sonarr"]
}

func (*systemConfig) ApiGetRadarrRules() []model.Rule {
	rules := make([]model.Rule, 0)
	content, err := utils.ReadFile(filepath.Join(consts.RULE_FILE_DIR, "radarr.toml"))
	if err != nil {
		log.Err(err).Msg("fail to read radarr.toml")
		return rules
	}
	err = toml.Unmarshal([]byte(content), &rules)
	if err != nil {
		log.Err(err).Msg("fail to unmarshal radarr.toml")
	}
	return rules
}

func (*systemConfig) ApiConfigUpdate(userInfo model.SystemUser, configs []model.SystemConfig) error {
	configMap := lo.Reduce(configs, func(agg map[string]*model.SystemConfig, item model.SystemConfig, index int) map[string]*model.SystemConfig {
		agg[item.Key] = &item
		return agg
	}, map[string]*model.SystemConfig{})

	sonarrUrl := configMap[consts.SONARR_URL]
	sonarrApiKey := configMap[consts.SONARR_API_KEY]
	sonarrIndexerFormat := configMap[consts.SONARR_INDEXER_FORMAT]
	sonarrLanguage1 := configMap[consts.SONARR_LANGUAGE_1]
	sonarrLanguage2 := configMap[consts.SONARR_LANGUAGE_2]
	radarrUrl := configMap[consts.RADARR_URL]
	radarrApiKey := configMap[consts.RADARR_API_KEY]
	radarrIndexerFormat := configMap[consts.RADARR_INDEXER_FORMAT]
	jackettUrl := configMap[consts.JACKETT_URL]
	prowlarrUrl := configMap[consts.PROWLARR_URL]
	qbittorrentUrl := configMap[consts.QBITTORRENT_URL]
	qbittorrentUsername := configMap[consts.QBITTORRENT_USERNAME]
	qbittorrentPassword := configMap[consts.QBITTORRENT_PASSWORD]
	transmissionUrl := configMap[consts.TRANSMISSION_URL]
	transmissionUsername := configMap[consts.TRANSMISSION_USERNAME]
	transmissionPassword := configMap[consts.TRANSMISSION_PASSWORD]
	tmdbUrl := configMap[consts.TMDB_URL]
	tmdbApikey := configMap[consts.TMDB_API_KEY]
	cleanTitleRegex := configMap[consts.CLEAN_TITLE_REGEX]
	ruleSyncAuthors := configMap[consts.RULE_SYNC_AUTHORS]

	_ = sonarrLanguage1
	_ = sonarrLanguage2
	_ = jackettUrl
	_ = prowlarrUrl
	_ = transmissionUrl
	_ = transmissionUsername
	_ = transmissionPassword
	_ = tmdbUrl
	_ = tmdbApikey
	_ = cleanTitleRegex
	_ = ruleSyncAuthors

	if Sonarr.ExternalCheckHealth(sonarrUrl.Value, sonarrApiKey.Value) == nil {
		sonarrUrl.ValidStatus = consts.VALID_STATUS
		sonarrApiKey.ValidStatus = consts.VALID_STATUS
	} else {
		sonarrUrl.ValidStatus = consts.INVALID_STATUS
		sonarrApiKey.ValidStatus = consts.INVALID_STATUS
	}

	if Sonarr.CheckIndexerFormat(sonarrIndexerFormat.Value) {
		sonarrIndexerFormat.ValidStatus = consts.VALID_STATUS
	} else {
		sonarrIndexerFormat.ValidStatus = consts.INVALID_STATUS
	}

	if Radarr.ExternalCheckHealth(radarrUrl.Value, radarrApiKey.Value) {
		radarrUrl.ValidStatus = consts.VALID_STATUS
		radarrApiKey.ValidStatus = consts.VALID_STATUS
	} else {
		radarrUrl.ValidStatus = consts.INVALID_STATUS
		radarrApiKey.ValidStatus = consts.INVALID_STATUS
	}
	if Radarr.CheckIndexerFormat(radarrIndexerFormat.Value) {
		radarrIndexerFormat.ValidStatus = consts.VALID_STATUS
	} else {
		radarrIndexerFormat.ValidStatus = consts.INVALID_STATUS
	}

	if Qbittorrent.Login(qbittorrentUrl.Value, qbittorrentUsername.Value, qbittorrentPassword.Value) {
		qbittorrentUrl.ValidStatus = consts.VALID_STATUS
		qbittorrentUsername.ValidStatus = consts.VALID_STATUS
		qbittorrentPassword.ValidStatus = consts.VALID_STATUS
	} else {
		qbittorrentUrl.ValidStatus = consts.INVALID_STATUS
		qbittorrentUsername.ValidStatus = consts.INVALID_STATUS
		qbittorrentPassword.ValidStatus = consts.INVALID_STATUS
	}
	// TODO: transmission

	if TMDB.ExternalCheckHealth(tmdbUrl.Value, tmdbApikey.Value) {
		tmdbUrl.ValidStatus = consts.VALID_STATUS
		tmdbApikey.ValidStatus = consts.VALID_STATUS
	} else {
		tmdbUrl.ValidStatus = consts.INVALID_STATUS
		tmdbApikey.ValidStatus = consts.INVALID_STATUS
	}
	if Jackett.ExternalCheckHealth(jackettUrl.Value) {
		jackettUrl.ValidStatus = consts.VALID_STATUS
	} else {
		jackettUrl.ValidStatus = consts.INVALID_STATUS
	}

	if Prowlarr.ExternalCheckHealth(prowlarrUrl.Value) {
		prowlarrUrl.ValidStatus = consts.VALID_STATUS
	} else {
		prowlarrUrl.ValidStatus = consts.INVALID_STATUS
	}

	if utils.IsRegex(cleanTitleRegex.Value) {
		cleanTitleRegex.ValidStatus = consts.VALID_STATUS
	} else {
		cleanTitleRegex.ValidStatus = consts.INVALID_STATUS
	}

	for _, systemConfig := range configMap {
		if err := db.Get().Save(&systemConfig).Error; err != nil {
			return err
		}
	}
	return nil
}

func (*systemConfig) MustConfigQueryByKey(key string) string {
	config := model.SystemConfig{}
	if err := db.Get().First(&config, "key = ?", key).Error; err != nil {
		log.Fatal().Err(err).Msg("fail to query system config")
	}
	if config.ValidStatus == consts.INVALID_STATUS {
		log.Fatal().Msg("key invalid: " + key)
	}
	return config.Value
}
