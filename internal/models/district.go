package models

type District struct {
	ID         int         `json:"id"`
	Name       string      `json:"name"`
	CityID     int         `json:"city_id"`
	Properties interface{} `json:"properties"`
}
