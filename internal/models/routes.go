package models

type RouteBase struct {
	Name       string      `json:"name"`
	CityID     int         `json:"city_id"`
	Price      int         `json:"price"`
	Properties interface{} `json:"properties"`
}

type RouteCreate struct {
	RouteBase
}

type Route struct {
	ID                 int `json:"id"`
	Name               string
	Price              int
	Pos                int
	City               string
	District           string
	RouteProperties    interface{}
	DistrictProperties interface{}
	PlaceProperties    interface{}
}
