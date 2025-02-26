package models

type User struct {
	Id       int    `json:"id"`
	Username string `json:"usernmae"`
	Password string `json:"password"`
}
