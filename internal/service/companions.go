package service

import (
	"context"
	"mth/internal/models"
	"mth/internal/repository"
	"mth/pkg/log"
)

type companionsService struct {
	placeRepo repository.Companions
	logger    *log.Logs
}

func InitCompanionsService(companionsRepo repository.Companions, logger *log.Logs) Companions {
	return companionsService{
		placeRepo: companionsRepo,
		logger:    logger,
	}
}

func (c companionsService) CreatePlaceCompanions(ctx context.Context, companion models.CompanionsPlaceCreate) error {
	err := c.placeRepo.CreatePlaceCompanions(ctx, companion)
	if err != nil {
		c.logger.Error(err.Error())
		return err
	}

	return nil
}

func (c companionsService) CreateRouteCompanions(ctx context.Context, companion models.CompanionsRouteCreate) error {
	err := c.placeRepo.CreateRouteCompanions(ctx, companion)
	if err != nil {
		c.logger.Error(err.Error())
		return err
	}

	return nil
}
