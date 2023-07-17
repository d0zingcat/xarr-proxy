package req

import (
	"net/http"

	"xarr-proxy/internal/model"
)

type (
	SystemConfigUpdateReq []model.SystemConfig
)

func (*SystemConfigUpdateReq) Bind(r *http.Request) error {
	return nil
}
