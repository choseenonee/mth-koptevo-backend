package swagger

type TripBase struct {
	EntityID int `json:"entity_id"`
	TripID   int `json:"trip_id"`
}

type TripChangePosition struct {
	Position int `json:"position"`
	TripBase
}

type TripChangeDay struct {
	Day int `json:"day"`
	TripBase
}

type TripAdd struct {
	Day      int `json:"day"`
	Position int `json:"position"`
	TripBase
}
