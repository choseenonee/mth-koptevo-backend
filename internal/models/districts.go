package models

type District struct {
	ID         int         `json:"id"`
	Name       string      `json:"name"`
	Properties interface{} `json:"properties"`
}
