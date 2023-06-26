package req

import (
	"net/http"

	"github.com/mcuadros/go-defaults"
)

type (
	ReqPage struct {
		Current  int `json:"current" default:"1"`
		PageSize int `json:"pageSize" default:"10"`
	}
	SonarrTitleQueryReq struct {
		ReqPage
		TvdbID string `json:"tvdbId"`
		Title  string `json:"title"`
	}

	SonarrTitleDeleteReq []int

	SonarrRuleQueryReq struct {
		ReqPage
		Token  string `json:"token"`
		Remark string `json:"remark"`
	}
	SonarrRuleSwitchValidStatusReq []string
)

func (s *SonarrTitleQueryReq) Bind(r *http.Request) error {
	defaults.SetDefaults(s)
	return nil
}

func (s *SonarrTitleDeleteReq) Bind(r *http.Request) error {
	return nil
}

func (s *SonarrRuleQueryReq) Bind(r *http.Request) error {
	defaults.SetDefaults(s)
	return nil
}

func (s *SonarrRuleSwitchValidStatusReq) Bind(r *http.Request) error {
	return nil
}
