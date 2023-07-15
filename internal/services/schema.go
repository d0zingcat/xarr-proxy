package services

import (
	"net/http"
	"net/url"
	"time"
)

type (
	// filter
	RequestWrapper struct {
		request      *http.Request
		parameterMap map[string][]string
	}
	ResponseWrapper struct{}
	// sonarr
	SonarrHistory []struct {
		EpisodeID   int    `json:"episodeId"`
		SeriesID    int    `json:"seriesId"`
		SourceTitle string `json:"sourceTitle"`
		Languages   []struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"languages"`
		Quality struct {
			Quality struct {
				ID         int    `json:"id"`
				Name       string `json:"name"`
				Source     string `json:"source"`
				Resolution int    `json:"resolution"`
			} `json:"quality"`
			Revision struct {
				Version  int  `json:"version"`
				Real     int  `json:"real"`
				IsRepack bool `json:"isRepack"`
			} `json:"revision"`
		} `json:"quality"`
		CustomFormats []struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"customFormats"`
		CustomFormatScore   int       `json:"customFormatScore"`
		QualityCutoffNotMet bool      `json:"qualityCutoffNotMet"`
		Date                time.Time `json:"date"`
		DownloadID          string    `json:"downloadId"`
		EventType           string    `json:"eventType"`
		Data                struct {
			Indexer            string    `json:"indexer"`
			NzbInfoURL         string    `json:"nzbInfoUrl"`
			ReleaseGroup       string    `json:"releaseGroup"`
			Age                string    `json:"age"`
			AgeHours           string    `json:"ageHours"`
			AgeMinutes         string    `json:"ageMinutes"`
			PublishedDate      time.Time `json:"publishedDate"`
			DownloadClient     string    `json:"downloadClient"`
			DownloadClientName string    `json:"downloadClientName"`
			Size               string    `json:"size"`
			DownloadURL        string    `json:"downloadUrl"`
			GUID               string    `json:"guid"`
			TvdbID             string    `json:"tvdbId"`
			TvRageID           string    `json:"tvRageId"`
			Protocol           string    `json:"protocol"`
			CustomFormatScore  string    `json:"customFormatScore"`
			SeriesMatchType    string    `json:"seriesMatchType"`
			ReleaseSource      string    `json:"releaseSource"`
			TorrentInfoHash    string    `json:"torrentInfoHash"`
		} `json:"data"`
		ID int `json:"id"`
	}
	SonarrSeries []struct {
		Title             string                  `json:"title"`
		AlternateTitles   []SonarrAlternateTitles `json:"alternateTitles"`
		SortTitle         string                  `json:"sortTitle"`
		Status            string                  `json:"status"`
		Ended             bool                    `json:"ended"`
		Overview          string                  `json:"overview"`
		PreviousAiring    time.Time               `json:"previousAiring"`
		Network           string                  `json:"network,omitempty"`
		AirTime           string                  `json:"airTime,omitempty"`
		Images            []SonarrImages          `json:"images"`
		OriginalLanguage  SonarrOriginalLanguage  `json:"originalLanguage"`
		Seasons           []SonarrSeasons         `json:"seasons"`
		Year              int                     `json:"year"`
		Path              string                  `json:"path"`
		QualityProfileID  int                     `json:"qualityProfileId"`
		SeasonFolder      bool                    `json:"seasonFolder"`
		Monitored         bool                    `json:"monitored"`
		UseSceneNumbering bool                    `json:"useSceneNumbering"`
		Runtime           int                     `json:"runtime"`
		TvdbID            int                     `json:"tvdbId"`
		TvRageID          int                     `json:"tvRageId"`
		TvMazeID          int                     `json:"tvMazeId"`
		FirstAired        time.Time               `json:"firstAired"`
		SeriesType        string                  `json:"seriesType"`
		CleanTitle        string                  `json:"cleanTitle"`
		ImdbID            string                  `json:"imdbId,omitempty"`
		TitleSlug         string                  `json:"titleSlug"`
		RootFolderPath    string                  `json:"rootFolderPath"`
		Certification     string                  `json:"certification,omitempty"`
		Genres            []string                `json:"genres"`
		Tags              []any                   `json:"tags"`
		Added             time.Time               `json:"added"`
		Ratings           SonarrRatings           `json:"ratings"`
		Statistics        SonarrStatistics        `json:"statistics"`
		LanguageProfileID int                     `json:"languageProfileId"`
		ID                int                     `json:"id"`
		NextAiring        time.Time               `json:"nextAiring,omitempty"`
	}
	SonarrAlternateTitles struct {
		Title             string `json:"title"`
		SceneSeasonNumber int    `json:"sceneSeasonNumber"`
	}
	SonarrImages struct {
		CoverType string `json:"coverType"`
		URL       string `json:"url"`
		RemoteURL string `json:"remoteUrl"`
	}
	SonarrOriginalLanguage struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
	SonarrStatistics struct {
		PreviousAiring    time.Time `json:"previousAiring"`
		EpisodeFileCount  int       `json:"episodeFileCount"`
		EpisodeCount      int       `json:"episodeCount"`
		TotalEpisodeCount int       `json:"totalEpisodeCount"`
		SizeOnDisk        int64     `json:"sizeOnDisk"`
		ReleaseGroups     []string  `json:"releaseGroups"`
		PercentOfEpisodes float64   `json:"percentOfEpisodes"`
		SeasonCount       int       `json:"seasonCount"`
	}
	SonarrSeasons struct {
		SeasonNumber int              `json:"seasonNumber"`
		Monitored    bool             `json:"monitored"`
		Statistics   SonarrStatistics `json:"statistics"`
	}
	SonarrRatings struct {
		Votes int `json:"votes"`
		Value int `json:"value"`
	}

	// tmdb
	TMDBSearchResult struct {
		MovieResults     []any                  `json:"movie_results"`
		PersonResults    []any                  `json:"person_results"`
		TvResults        []TMDBTvResults        `json:"tv_results"`
		TvEpisodeResults []TMDBTvEpisodeResults `json:"tv_episode_results"`
		TvSeasonResults  []any                  `json:"tv_season_results"`
	}

	TMDBTvResults struct {
		Adult            bool     `json:"adult"`
		BackdropPath     string   `json:"backdrop_path"`
		ID               int      `json:"id"`
		Name             string   `json:"name"`
		OriginalLanguage string   `json:"original_language"`
		OriginalName     string   `json:"original_name"`
		Overview         string   `json:"overview"`
		PosterPath       string   `json:"poster_path"`
		MediaType        string   `json:"media_type"`
		GenreIds         []int    `json:"genre_ids"`
		Popularity       float64  `json:"popularity"`
		FirstAirDate     string   `json:"first_air_date"`
		VoteAverage      float64  `json:"vote_average"`
		VoteCount        int      `json:"vote_count"`
		OriginCountry    []string `json:"origin_country"`
	}

	TMDBTvEpisodeResults struct {
		ID             int     `json:"id"`
		Name           string  `json:"name"`
		Overview       string  `json:"overview"`
		MediaType      string  `json:"media_type"`
		VoteAverage    float64 `json:"vote_average"`
		VoteCount      int     `json:"vote_count"`
		AirDate        string  `json:"air_date"`
		EpisodeNumber  int     `json:"episode_number"`
		ProductionCode string  `json:"production_code"`
		Runtime        int     `json:"runtime"`
		SeasonNumber   int     `json:"season_number"`
		ShowID         int     `json:"show_id"`
		StillPath      string  `json:"still_path"`
	}
)

// 创建新的请求包装器
func NewRequestWrapper(request *http.Request) *RequestWrapper {
	parameterMap := make(map[string][]string, len(request.URL.Query()))
	for key, values := range request.URL.Query() {
		parameterMap[key] = values
	}
	return &RequestWrapper{
		request:      request,
		parameterMap: parameterMap,
	}
}

// 设置参数
func (r *RequestWrapper) SetParameter(name, value string) {
	r.parameterMap[name] = []string{value}
}

// 获取参数映射
func (r *RequestWrapper) GetParameterMap() map[string][]string {
	return r.parameterMap
}

// 获取参数名称列表
func (r *RequestWrapper) GetParameterNames() []string {
	names := make([]string, 0, len(r.parameterMap))
	for name := range r.parameterMap {
		names = append(names, name)
	}
	return names
}

// 获取参数值
func (r *RequestWrapper) GetParameterValues(name string) []string {
	return r.parameterMap[name]
}

// 获取单个参数值
func (r *RequestWrapper) GetParameter(name string) string {
	values := r.parameterMap[name]
	if len(values) > 0 {
		return values[0]
	}
	return ""
}

// 获取查询字符串
func (r *RequestWrapper) GetQueryString() string {
	queryString := url.Values{}
	for key, values := range r.parameterMap {
		for _, value := range values {
			queryString.Add(key, value)
		}
	}
	return queryString.Encode()
}

func (r *RequestWrapper) GetPath() string {
	return r.request.URL.Path
}
