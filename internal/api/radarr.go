package api

import (
	"net/http"

	"xarr-proxy/internal/services"

	"github.com/go-chi/render"
)

func radarrTitleQuery(w http.ResponseWriter, r *http.Request) {
	v := services.SystemConfig.ApiVersion()
	render.JSON(w, r, v)
}
