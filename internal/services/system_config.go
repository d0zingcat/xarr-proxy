package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"

	"xarr-proxy/internal/consts"
	"xarr-proxy/internal/utils"

	"github.com/pelletier/go-toml/v2"

	"github.com/rs/zerolog/log"
)

var SystemConfig = &systemConfig{}

type systemConfig struct{}

func (*systemConfig) Version() string {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", consts.AUTHOR, consts.REPO)
	resp, err := http.Get(url)
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
		if v != consts.VERSION {
			return consts.VERSION + " 🚨"
		}
	}
	return consts.VERSION
}

func (*systemConfig) AuthorList() []string {
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
