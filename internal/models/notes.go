package models

type NoteCreate struct {
	UserID     int         `json:"user_id"`
	PlaceID    int         `json:"place_id"`
	Properties interface{} `json:"properties"`
}

type Note struct {
	ID int `json:"id"`
	NoteCreate
}

type NoteUpdate struct {
	ID         int         `json:"id"`
	Properties interface{} `json:"properties"`
}
