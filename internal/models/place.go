package models

type PlaceBase struct {
	Properties interface{} `json:"properties"`
	CityID     int         `json:"city_id"`
	DistrictID int         `json:"district_id"`
	Name       string      `json:"name"`
}

type PlaceCreate struct {
	TagIDs []int `json:"tag_ids"`
	PlaceBase
}

type Place struct {
	ID   int   `json:"id"`
	Tags []Tag `json:"tags"`
	PlaceBase
}
