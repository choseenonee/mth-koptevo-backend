package service

import (
	"context"
	"mth/internal/models"
	"mth/internal/repository"
	"mth/pkg/log"
)

type companionsService struct {
	companionRepo repository.Companions
	logger        *log.Logs
}

func InitCompanionsService(companionsRepo repository.Companions, logger *log.Logs) Companions {
	return companionsService{
		companionRepo: companionsRepo,
		logger:        logger,
	}
}

func (c companionsService) CreatePlaceCompanions(ctx context.Context, companion models.CompanionsPlaceCreate) error {
	err := c.companionRepo.CreatePlaceCompanions(ctx, companion)
	if err != nil {
		c.logger.Error(err.Error())
		return err
	}

	return nil
}

func (c companionsService) CreateRouteCompanions(ctx context.Context, companion models.CompanionsRouteCreate) error {
	err := c.companionRepo.CreateRouteCompanions(ctx, companion)
	if err != nil {
		c.logger.Error(err.Error())
		return err
	}

	return nil
}

func (c companionsService) GetByUser(ctx context.Context, userID int) ([]models.CompanionsPlace, []models.CompanionsRoute, error) {
	places, routes, err := c.companionRepo.GetByUser(ctx, userID)
	if err != nil {
		c.logger.Error(err.Error())
		return places, routes, err
	}

	return places, routes, nil
}

func (c companionsService) GetCompanionsPlace(ctx context.Context, filters models.CompanionsFilters) ([]models.CompanionsPlace, error) {
	places, err := c.companionRepo.GetCompanionsPlace(ctx, filters)
	if err != nil {
		c.logger.Error(err.Error())
		return places, err
	}

	return places, nil
}

func (c companionsService) GetCompanionsRoute(ctx context.Context, filters models.CompanionsFilters) ([]models.CompanionsRoute, error) {
	routes, err := c.companionRepo.GetCompanionsRoute(ctx, filters)
	if err != nil {
		c.logger.Error(err.Error())
		return routes, err
	}

	return routes, nil
}

func (c companionsService) DeleteCompanionsPlace(ctx context.Context, id int) error {
	err := c.companionRepo.DeleteCompanionsPlace(ctx, id)
	if err != nil {
		c.logger.Error(err.Error())
		return err
	}

	return nil
}

func (c companionsService) DeleteCompanionsRoute(ctx context.Context, id int) error {
	err := c.companionRepo.DeleteCompanionsRoute(ctx, id)
	if err != nil {
		c.logger.Error(err.Error())
		return err
	}

	return nil
}
