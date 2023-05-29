package api

import (
	"net/http"

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
