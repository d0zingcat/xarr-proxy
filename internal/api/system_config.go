package api

import (
	"errors"
	"net/http"

	"xarr-proxy/internal/api/req"
	"xarr-proxy/internal/consts"
	"xarr-proxy/internal/model"
	"xarr-proxy/internal/services"

	"github.com/go-chi/render"
)

func systemVersion(w http.ResponseWriter, r *http.Request) {
	v := services.SystemConfig.Version()
	render.JSON(w, r, v)
}

func authorList(w http.ResponseWriter, r *http.Request) {
	v := services.SystemConfig.AuthorList()
	render.JSON(w, r, v)
}

func configQuery(w http.ResponseWriter, r *http.Request) {
	// TODO: sonar rename task
	v := services.SystemConfig.ConfigQuery()
	render.JSON(w, r, v)
}

func configUpdate(w http.ResponseWriter, r *http.Request) {
	userInfo := r.Context().Value(consts.USER_INFO_CTX_KEY)
	var req req.SystemConfigUpdateReq
	if err := render.Bind(r, &req); err != nil {
		render.JSON(w, r, ErrInvalidRequest(err))
		return
	}
	if userInfo == nil {
		render.JSON(w, r, ErrInvalidRequest(errors.New("invalid user info")))
		return
	}
	v := services.SystemConfig.ConfigUpdate(userInfo.(model.SystemUser), []model.SystemConfig(req))
	// TODO: clear title sync cache
	render.JSON(w, r, v)
}
