package model

import "time"

type (
	SystemUser struct {
		Id          int    `json:"id"`
		Username    string `json:"username"`
		Role        string `json:"role"`
		ValidStatus int    `json:"valid_status"`
	}

	SystemConfig struct {
		Id          int        `json:"id"`
		Key         string     `json:"key"`
		Value       string     `json:"value"`
		ValidStatus int        `json:"validStatus"`
		CreateTime  *time.Time `json:"createTime"`
		UpdateTime  *time.Time `json:"updateTime"`
	}
)
