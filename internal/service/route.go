package service

import (
	"context"
	"mth/internal/models"
	"mth/internal/repository"
	"mth/pkg/log"
)

type routeService struct {
	routeRepo repository.Route
	placeRepo repository.Place
	logger    *log.Logs
}

func InitRouteService(routeRepo repository.Route, placeRepo repository.Place, logger *log.Logs) Route {
	return routeService{
		routeRepo: routeRepo,
		placeRepo: placeRepo,
		logger:    logger,
	}
}

func (r routeService) Create(ctx context.Context, route models.RouteCreate) (int, error) {
	id, err := r.routeRepo.Create(ctx, route)
	if err != nil {
		r.logger.Error(err.Error())
		return 0, err
	}

	return id, nil
}

func (r routeService) GetByID(ctx context.Context, routeID int) (models.Route, error) {
	routeRaw, err := r.routeRepo.GetByID(ctx, routeID)
	if err != nil {
		r.logger.Error(err.Error())
		return models.Route{}, err
	}

	var route models.Route

	route.RouteBase = routeRaw.RouteBase
	route.ID = routeRaw.ID
	route.Tags = routeRaw.Tags

	for _, placeID := range routeRaw.PlaceIDs {
		place, err := r.placeRepo.GetByID(ctx, placeID)
		if err != nil {
			r.logger.Error(err.Error())
			return models.Route{}, err
		}

		route.Places = append(route.Places, place)
	}

	return route, nil
}

func (r routeService) GetAll(ctx context.Context, page int) ([]models.Route, error) {
	//TODO implement me
	panic("implement me")
}
