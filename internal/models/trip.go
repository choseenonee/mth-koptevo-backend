package models

import "time"

type EntityWithDayAndPosition struct {
	EntityID int `json:"id"`
	Day      int `json:"day"`
	Position int `json:"position"`
}

type TripBase struct {
	UserID     int         `json:"user_id"`
	DateStart  time.Time   `json:"date_start"`
	DateEnd    time.Time   `json:"date_end"`
	Properties interface{} `json:"properties"`
}

type TripCreate struct {
	Places []EntityWithDayAndPosition `json:"places"`
	Routes []EntityWithDayAndPosition `json:"Routes"`
	TripBase
}

type TripRaw struct {
	PlaceIDs []int
	RouteIDs []int
	TripBase
}

type Trip struct {
	ID int `json:"id"`
	TripCreate
}
