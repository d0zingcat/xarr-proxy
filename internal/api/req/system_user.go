package req

import (
	"net/http"
)

type (
	SystemUserLoginReq struct {
		Username string
		Password string
	}
)

func (*SystemUserLoginReq) Bind(r *http.Request) error {
	return nil
}
