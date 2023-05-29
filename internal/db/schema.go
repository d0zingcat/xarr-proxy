package db

import (
	"time"

	"xarr-proxy/internal/model"
)

type (
	BaseModel struct {
		Id         int        `json:"id"`
		CreateTime *time.Time `json:"create_time"`
		UpdateTime *time.Time `json:"update_time"`
	}
	SystemUser struct {
		model.SystemUser
		Password   string     `json:"password"`
		CreateTime *time.Time `json:"create_time"`
		UpdateTime *time.Time `json:"update_time"`
	}

	SystemConfig struct {
		model.SystemConfig
	}
)
