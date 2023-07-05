package resp

type Pagination struct {
	Current  int `json:"current"`
	PageSize int `json:"pageSize"`
	Total    int `json:"total"`
}

type RowsResponse[T any] struct {
	Pagination
	List []T `json:"list"`
}
