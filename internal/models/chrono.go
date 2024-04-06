package models

import "time"

type ChronoEntity struct {
	ID        int       `json:"id"`
	TimeStamp time.Time `json:"timestamp"`
	TripID    int       `json:"trip_id"`
}

type Chrono struct {
	LikedRoutes  []ChronoEntity `json:"liked_routes"`
	LikedPlaces  []ChronoEntity `json:"liked_places"`
	PlaceReviews []ChronoEntity `json:"place_reviews"`
	RouteReviews []ChronoEntity `json:"route_reviews"`
	CheckIns     []ChronoEntity `json:"check_ins"`
	RouteLogs    []ChronoEntity `json:"route_logs"`
}
