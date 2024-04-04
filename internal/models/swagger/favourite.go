package swagger

import "mth/internal/models"

type Favourite struct {
	Places    []models.Place    `json:"places"`
	RoutesRaw []models.RouteRaw `json:"routes"`
}

func InitFavourite(places []models.Place, routesRaw []models.RouteRaw) Favourite {
	return Favourite{
		Places:    places,
		RoutesRaw: routesRaw,
	}
}
