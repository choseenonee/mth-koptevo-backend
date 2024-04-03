package models

import "time"

type ReviewBase struct {
	AuthorID   int         `json:"author_id"`
	Properties interface{} `json:"properties"`
	Mark       float32     `json:"mark"`
	TimeStamp  time.Time
}

type ReviewUpdate struct {
	ID         int         `json:"id"`
	Properties interface{} `json:"properties"`
	Mark       float32     `json:"mark"`
}

type PlaceReviewCreate struct {
	PlaceID int `json:"place_id"`
	ReviewBase
}

type RouteReviewCreate struct {
	RouteID int `json:"route_id"`
	ReviewBase
}

type PlaceReview struct {
	ID int `json:"id"`
	PlaceReviewCreate
}

type RouteReview struct {
	ID int `json:"id"`
	RouteReviewCreate
}
