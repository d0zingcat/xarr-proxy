package model

type (
	SystemUser struct {
		Id          int    `json:"id"`
		Username    string `json:"username"`
		Role        string `json:"role"`
		ValidStatus int    `json:"valid_status"`
	}
)
