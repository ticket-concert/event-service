package constants

type MetaData struct {
	Page      int64 `json:"page"`
	Count     int64 `json:"count"`
	TotalPage int64 `json:"totalPage"`
	TotalData int64 `json:"totalData"`
}

const (
	Online = "Online"
	Gold   = "Gold"
	Silver = "Silver"
	Bronze = "Bronze"
	Wood   = "Wood"
)
