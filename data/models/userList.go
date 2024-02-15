package models

type UserList struct {
	Uid   string   `json:"uid"`
	Users []string `json:"users"`
}
