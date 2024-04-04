package models

import "time"

type CompanionsFilters struct {
	Page     int       `json:"page"`
	EntityID int       `json:"entity_id"`
	DateFrom time.Time `json:"date_from"`
	DateTo   time.Time `json:"date_to"`
}

type CompanionCreateBase struct {
	UserID   int       `json:"user_id"`
	DateFrom time.Time `json:"date_from"`
	DateTo   time.Time `json:"date_to"`
}

type CompanionsPlaceCreate struct {
	PlaceID int `json:"place_id"`
	CompanionCreateBase
}

type CompanionsRouteCreate struct {
	RouteID int `json:"route_id"`
	CompanionCreateBase
}

type CompanionBase struct {
	ID             int         `json:"id"`
	UserID         int         `json:"user_id"`
	UserProperties interface{} `json:"user_properties"`
	CityName       string      `json:"city_name"`
	DateFrom       time.Time   `json:"date_from"`
	DateTo         time.Time   `json:"date_to"`
}

type CompanionsPlace struct {
	PlaceID         int         `json:"place_id"`
	PlaceName       string      `json:"place_name"`
	PlaceProperties interface{} `json:"place_properties"`
	CompanionBase
}

type CompanionsRoute struct {
	RouteID         int         `json:"route_id"`
	RouteName       string      `json:"route_name"`
	Price           int         `json:"price"`
	RouteProperties interface{} `json:"route_properties"`
	CompanionBase
}
