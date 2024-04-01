package models

type TagCreate struct {
	Name string `json:"name"`
}

type Tag struct {
	ID int `json:"id"`
	TagCreate
}
