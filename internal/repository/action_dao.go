package repository

type ActionDao struct {
	Id          int    `json:"id"`
	Key         string `json:"key"`
	Channel     string `json:"channel"`
	Description string `json:"description"`
	IsDeleted   bool   `json:"is_deleted"`
}
