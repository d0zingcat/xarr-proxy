package req

import "net/http"

type (
	TMDBTitleQueryReq struct {
		ReqPage
		TvdbID string `json:"tvdbId"`
		Title  string `json:"title"`
	}
	TMDBTitleDeleteReq []int
)

func (s *TMDBTitleDeleteReq) Bind(r *http.Request) error {
	return nil
}
