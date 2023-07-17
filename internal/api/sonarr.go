package api

import (
	"errors"
	"net/http"
	"strings"

	"xarr-proxy/internal/api/req"
	"xarr-proxy/internal/consts"
	"xarr-proxy/internal/services"
	"xarr-proxy/internal/services/filter"
	"xarr-proxy/internal/services/indexer"
	"xarr-proxy/internal/utils"

	"github.com/go-chi/render"
)

func sonarrProwlarrProxy(w http.ResponseWriter, r *http.Request) {
	f := filter.NewIndexerFilter(indexer.NewIndexerService(indexer.NewSonarrProwlarrService()))
	f.DoFilter(w, r)
}

func sonarrTitleQuery(w http.ResponseWriter, r *http.Request) {
	req := new(req.SonarrTitleQueryReq)

	req.Current = utils.ParseInt(r.URL.Query().Get("current"))
	req.PageSize = utils.ParseInt(r.URL.Query().Get("pageSize"))
	req.Title = strings.Trim(r.URL.Query().Get("title"), " ")
	req.TvdbID = strings.Trim(r.URL.Query().Get("tvdbId"), " ")

	v, err := services.Sonarr.ApiQuery(req)
	if err != nil {
		render.JSON(w, r, ErrInternalServer(err))
		return
	}
	render.JSON(w, r, v)
}

func sonarrTitleRemove(w http.ResponseWriter, r *http.Request) {
	req := new(req.SonarrTitleDeleteReq)
	if err := render.Bind(r, req); err != nil {
		render.JSON(w, r, ErrInvalidRequest(err))
		return
	}

	v, err := services.Sonarr.ApiBulkDelete(*req)
	if err != nil {
		render.JSON(w, r, ErrInternalServer(err))
		return
	}
	render.JSON(w, r, v)
}

func sonarrTitleSync(w http.ResponseWriter, r *http.Request) {
	v, err := services.Sonarr.ApiSync()
	if err != nil {
		render.JSON(w, r, ErrInternalServer(err))
		return
	}
	render.JSON(w, r, v)
}

func sonarrRuleSync(w http.ResponseWriter, r *http.Request) {
	v, err := services.Sonarr.ApiRuleSync()
	if err != nil {
		render.JSON(w, r, ErrInternalServer(err))
		return
	}
	render.JSON(w, r, v)
}

func sonarrRuleQuery(w http.ResponseWriter, r *http.Request) {
	req := new(req.SonarrRuleQueryReq)

	req.Current = utils.ParseInt(r.URL.Query().Get("current"))
	req.PageSize = utils.ParseInt(r.URL.Query().Get("pageSize"))
	req.Token = strings.Trim(r.URL.Query().Get("token"), " ")
	req.Remark = strings.Trim(r.URL.Query().Get("remark"), " ")
	v, err := services.Sonarr.ApiRuleQuery(req)
	if err != nil {
		render.JSON(w, r, ErrInternalServer(err))
		return
	}
	render.JSON(w, r, v)
}

func sonarrRuleEnable(w http.ResponseWriter, r *http.Request) {
	req := new(req.SonarrRuleSwitchValidStatusReq)
	if err := render.Bind(r, req); err != nil {
		render.JSON(w, r, ErrInvalidRequest(err))
		return
	}

	v, err := services.Sonarr.ApiSwitchValidStatus(*req, consts.VALID_STATUS)
	if err != nil {
		render.JSON(w, r, ErrInternalServer(err))
		return
	}
	render.JSON(w, r, v)
}

func sonarrRuleDisable(w http.ResponseWriter, r *http.Request) {
	req := new(req.SonarrRuleSwitchValidStatusReq)
	if err := render.Bind(r, req); err != nil {
		render.JSON(w, r, ErrInvalidRequest(err))
		return
	}
	for _, v := range []string(*req) {
		if v == consts.RULE_MAIN_ID {
			render.JSON(w, r, ErrInvalidRequest(errors.New("main rule can't be disabled")))
			return
		}
	}

	v, err := services.Sonarr.ApiSwitchValidStatus(*req, consts.INVALID_STATUS)
	if err != nil {
		render.JSON(w, r, ErrInternalServer(err))
		return
	}
	render.JSON(w, r, v)
}

func sonarrRuleTokenQuery(w http.ResponseWriter, r *http.Request) {
	v, err := services.Sonarr.ApiTokenList()
	if err != nil {
		render.JSON(w, r, ErrInternalServer(err))
		return
	}
	render.JSON(w, r, v)
}
