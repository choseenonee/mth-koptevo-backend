package models

type UserBase struct {
	Login      string      `json:"login"`
	Properties interface{} `json:"properties"`
}

type UserCreate struct {
	Password string `json:"password"`
	UserBase
}

type User struct {
	ID int `json:"id"`
	UserBase
}
