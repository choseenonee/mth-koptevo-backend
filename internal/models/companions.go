package models

import "time"

type CompanionsFilters struct {
	Page     int
	ID       int
	DateFrom time.Time
	DateTo   time.Time
}

type CompanionsPlaceCreate struct {
	UserID   int
	PlaceID  int
	DateFrom time.Time
	DateTo   time.Time
}

type CompanionsRouteCreate struct {
	UserID   int
	RouteID  int
	DateFrom time.Time
	DateTo   time.Time
}

type CompanionsPlace struct {
	UserProperties  interface{}
	PlaceName       string
	CityName        string
	PlaceProperties interface{}
	DateFrom        time.Time
	DateTo          time.Time
}

type CompanionsRoute struct {
	UserProperties  interface{}
	RouteName       string
	Price           int
	RouteProperties interface{}
	CityName        string
	DateFrom        time.Time
	DateTo          time.Time
}
