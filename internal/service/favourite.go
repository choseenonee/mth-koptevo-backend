package service

import (
	"context"
	"mth/internal/models"
	"mth/internal/repository"
	"mth/pkg/log"
)

type favouriteService struct {
	favouriteRepo repository.Favourite
	placeRepo     repository.Place
	routeRepo     repository.Route
	logger        *log.Logs
}

func InitFavouriteService(favouriteRepo repository.Favourite, placeRepo repository.Place, routeRepo repository.Route, logger *log.Logs) Favourite {
	return favouriteService{
		favouriteRepo: favouriteRepo,
		placeRepo:     placeRepo,
		routeRepo:     routeRepo,
		logger:        logger,
	}
}

func (f favouriteService) LikePlace(ctx context.Context, like models.Like) error {
	err := f.favouriteRepo.LikePlace(ctx, like)
	if err != nil {
		f.logger.Error(err.Error())
		return err
	}

	return nil
}

func (f favouriteService) LikeRoute(ctx context.Context, like models.Like) error {
	err := f.favouriteRepo.LikeRoute(ctx, like)
	if err != nil {
		f.logger.Error(err.Error())
		return err
	}

	return nil
}

// GetLikedByUser places and RAW routes
func (f favouriteService) GetLikedByUser(ctx context.Context, userID int) ([]models.Place, []models.RouteRaw, error) {
	placeIDs, routeIDs, err := f.favouriteRepo.GetLikedByUser(ctx, userID)
	if err != nil {
		f.logger.Error(err.Error())
		return []models.Place{}, []models.RouteRaw{}, err
	}

	places := make([]models.Place, 0, len(placeIDs))

	for _, placeID := range placeIDs {
		place, err := f.placeRepo.GetByID(ctx, placeID)
		if err != nil {
			f.logger.Error(err.Error())
			return []models.Place{}, []models.RouteRaw{}, err
		}

		places = append(places, place)
	}

	routesRaw := make([]models.RouteRaw, 0, len(placeIDs))
	for _, routeID := range routeIDs {
		routeRaw, err := f.routeRepo.GetByID(ctx, routeID)
		if err != nil {
			f.logger.Error(err.Error())
			return []models.Place{}, []models.RouteRaw{}, err
		}

		routesRaw = append(routesRaw, routeRaw)
	}

	return places, routesRaw, nil
}
