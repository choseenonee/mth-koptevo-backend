package service

import (
	"context"
	"mth/internal/models"
	"mth/internal/repository"
	"mth/pkg/log"
)

type placeService struct {
	placeRepo repository.Place
	logger    *log.Logs
}

func InitPlaceService(placeRepo repository.Place, logger *log.Logs) Place {
	return placeService{
		placeRepo: placeRepo,
		logger:    logger,
	}
}

func (p placeService) Create(ctx context.Context, placeCreate models.PlaceCreate) (int, error) {
	id, err := p.placeRepo.Create(ctx, placeCreate)
	if err != nil {
		p.logger.Error(err.Error())
		return 0, err
	}

	return id, nil
}

func (p placeService) GetAllWithFilter(ctx context.Context, districtID int, cityID int, tagIDs []int, page int, name string) ([]models.Place, error) {
	places, err := p.placeRepo.GetAllWithFilter(ctx, districtID, cityID, tagIDs, page, name)
	if err != nil {
		p.logger.Error(err.Error())
		return []models.Place{}, err
	}

	return places, nil
}

func (p placeService) GetByID(ctx context.Context, placeID int) (models.Place, error) {
	place, err := p.placeRepo.GetByID(ctx, placeID)
	if err != nil {
		p.logger.Error(err.Error())
		return models.Place{}, err
	}

	return place, nil
}
