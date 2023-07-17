package db

import (
	"time"

	"xarr-proxy/internal/model"
)

type (
	BaseModel struct {
		Id         int `json:"id"`
		CreateTime *time.Time
		UpdateTime *time.Time
	}
	SystemUser struct {
		model.SystemUser
		Password   string `json:"-"`
		CreateTime *time.Time
		UpdateTime *time.Time
	}

	SystemConfig struct {
		model.SystemConfig
	}

	SonarrTitle struct {
		ID           int    `json:"id"`
		TvdbID       int    `json:"tvdbId"`
		Sno          int    `json:"sno"`
		MainTitle    string `json:"mainTitle"`
		Title        string `json:"title"`
		CleanTitle   string `json:"cleanTitle"`
		SeasonNumber int    `json:"seasonNumber"`
		Monitored    int    `json:"monitored"`
		ValidStatus  int    `json:"validStatus"`
		SeriesID     int    `json:"seriesId"`
		CreateTime   *time.Time
		UpdateTime   *time.Time
	}

	TMDBTitle struct {
		ID          int        `json:"id"`
		TvdbID      int        `json:"tvdbId"`
		TmdbID      int        `json:"tmdbId"`
		Language    string     `json:"language"`
		Title       string     `json:"title"`
		ValidStatus int        `json:"validStatus"`
		CreateTime  *time.Time `json:"createTime"`
		UpdateTime  *time.Time `json:"updateTime"`
	}

	SonarrRule struct {
		ID          string    `gorm:"column:id;primaryKey" json:"id"`
		Token       string    `gorm:"column:token" json:"token"`
		Priority    int       `gorm:"column:priority" json:"priority"`
		Regex       string    `gorm:"column:regex" json:"regex"`
		Replacement string    `gorm:"column:replacement" json:"replacement"`
		Offset      int       `gorm:"column:offset" json:"offset"`
		Example     string    `gorm:"column:example" json:"example"`
		Remark      string    `gorm:"column:remark" json:"remark"`
		Author      string    `gorm:"column:author" json:"author"`
		ValidStatus int       `gorm:"column:valid_status" json:"validStatus"`
		CreateTime  time.Time `gorm:"column:create_time" json:"createTime"`
		UpdateTime  time.Time `gorm:"column:update_time" json:"updateTime"`
	}

	RadarrRule struct {
		SonarrRule
	}
)
