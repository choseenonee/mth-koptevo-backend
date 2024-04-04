package models

import "time"

type CompanionsFilters struct {
	Page     int
	EntityID int
	DateFrom time.Time
	DateTo   time.Time
}

type CompanionCreateBase struct {
	UserID   int
	DateFrom time.Time
	DateTo   time.Time
}

type CompanionsPlaceCreate struct {
	PlaceID int
	CompanionCreateBase
}

type CompanionsRouteCreate struct {
	RouteID int
	CompanionCreateBase
}

type CompanionsPlace struct {
	UserID          int
	UserProperties  interface{}
	PlaceName       string
	CityName        string
	PlaceProperties interface{}
	DateFrom        time.Time
	DateTo          time.Time
}

type CompanionsRoute struct {
	UserID          int
	UserProperties  interface{}
	RouteName       string
	Price           int
	RouteProperties interface{}
	CityName        string
	DateFrom        time.Time
	DateTo          time.Time
}
