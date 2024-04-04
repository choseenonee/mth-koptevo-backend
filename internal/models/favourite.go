package models

type Like struct {
	UserID   int `json:"user_id"`
	EntityID int `json:"entity_id"`
}
