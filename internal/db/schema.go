package db

import "time"

type (
	BaseModel struct {
		Id         int        `json:"id"`
		CreateTime *time.Time `json:"create_time"`
		UpdateTime *time.Time `json:"update_time"`
	}
	SystemUser struct {
		BaseModel
		Username    string `json:"username"`
		Password    string `json:"password"`
		Role        string `json:"role"`
		ValidStatus int    `json:"valid_status"`
	}
)
