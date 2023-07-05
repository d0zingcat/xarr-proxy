package req

type (
	IndexerReq struct {
		SearchKey     string `json:"searchKey"`
		SearchType    string `json:"searchType"`
		SeasonNumber  string `json:"seasonNumber"`
		EpisodeNumber string `json:"episodeNumber"`
		Offset        int    `json:"offset"`
		Limit         int    `json:"limit"`
	}
)
