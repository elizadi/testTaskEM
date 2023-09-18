package domain

type GetUsersResponse struct {
	Users    []User
	RespInfo RespInfo
}

type RespInfo struct {
	Total     int64
	PageCount int64
}
