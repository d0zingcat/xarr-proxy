package api

import (
	"net/http"
	"strings"

	"xarr-proxy/internal/api/req"
	"xarr-proxy/internal/services"
	"xarr-proxy/internal/utils"

	"github.com/go-chi/render"
)

func tmdbTitleQuery(w http.ResponseWriter, r *http.Request) {
	req := new(req.TMDBTitleQueryReq)

	req.Current = utils.ParseInt(r.URL.Query().Get("current"))
	req.PageSize = utils.ParseInt(r.URL.Query().Get("pageSize"))
	req.Title = strings.Trim(r.URL.Query().Get("title"), " ")
	req.TvdbID = strings.Trim(r.URL.Query().Get("tvdbId"), " ")

	v, err := services.TMDB.ApiQuery(req)
	if err != nil {
		render.JSON(w, r, ErrInternalServer(err))
		return
	}
	render.JSON(w, r, v)
}

func tmdbTitleRemove(w http.ResponseWriter, r *http.Request) {
	req := new(req.TMDBTitleDeleteReq)
	if err := render.Bind(r, req); err != nil {
		render.JSON(w, r, ErrInvalidRequest(err))
		return
	}
	v, err := services.TMDB.ApiBulkDelete(*req)
	if err != nil {
		render.JSON(w, r, ErrInternalServer(err))
		return
	}
	render.JSON(w, r, v)
}

func tmdbTitleSync(w http.ResponseWriter, r *http.Request) {
	v, err := services.TMDB.ApiSync()
	if err != nil {
		render.JSON(w, r, ErrInternalServer(err))
		return
	}
	render.JSON(w, r, v)
}
