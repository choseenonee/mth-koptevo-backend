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
