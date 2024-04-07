package models

type RouteBase struct {
	CityID     int         `json:"city_id"`
	Price      int         `json:"price"`
	Name       string      `json:"name"`
	Properties interface{} `json:"properties"`
}

type PlaceIDWithPosition struct {
	PlaceID  int `json:"place_id"`
	Position int `json:"position"`
}

type PlaceWithPosition struct {
	Place    Place `json:"place"`
	Position int   `json:"position"`
}

type RouteCreate struct {
	TagIDs               []int                 `json:"tag_ids"`
	PlaceIDsWithPosition []PlaceIDWithPosition `json:"places"`
	RouteBase
}

type RouteRaw struct {
	ID                   int                   `json:"id"`
	Tags                 []Tag                 `json:"tags"`
	PlaceIDsWithPosition []PlaceIDWithPosition `json:"place_ids"`
	RouteBase
}

type RouteDisplay struct {
	ID             int `json:"id"`
	NextPlaceID    int `json:"next_place_id"`
	CompletedPlace int `json:"completed_place"`
}

type Route struct {
	ID     int                 `json:"id"`
	Tags   []Tag               `json:"tags"`
	Places []PlaceWithPosition `json:"places"`
	RouteBase
}
