package model

import "time"

type (
	SystemUser struct {
		ID          int    `json:"id"`
		Username    string `json:"username"`
		Role        string `json:"role"`
		ValidStatus int    `json:"validStatus"`
	}

	SystemConfig struct {
		ID          int        `json:"id"`
		Key         string     `json:"key"`
		Value       string     `json:"value"`
		ValidStatus int        `json:"validStatus"`
		CreateTime  *time.Time `json:"createTime"`
		UpdateTime  *time.Time `json:"updateTime"`
	}

	Rule struct {
		ID          string `toml:"id"`
		Token       string `toml:"token"`
		Priority    int    `toml:"priority"`
		Regex       string `toml:"regex"`
		Replacement string `toml:"replacement"`
		Offset      int    `toml:"offset"`
		Example     string `toml:"example"`
		Remark      string `toml:"remark"`
		Author      string `toml:"author"`
	}
)
