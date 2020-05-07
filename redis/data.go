package redis

type UserInfoSimple struct{
	Phone string `json:"mobile"`
	Expire int64 `json:"expire"`
}