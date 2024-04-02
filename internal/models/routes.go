package models

type RouteBase struct {
	CityID     int         `json:"city_id"`
	Price      int         `json:"price"`
	Name       string      `json:"name"`
	Properties interface{} `json:"properties"`
}

type RouteCreate struct {
	TagIDs   []int `json:"tag_ids"`
	PlaceIDs []int `json:"place_ids"`
	RouteBase
}

type RouteRaw struct {
	ID       int   `json:"id"`
	Tags     []Tag `json:"tags"`
	PlaceIDs []int `json:"place_ids"`
	RouteBase
}

type Route struct {
	ID     int     `json:"id"`
	Tags   []Tag   `json:"tags"`
	Places []Place `json:"places"`
	RouteBase
}
