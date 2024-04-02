package service

import (
	"context"
	"mth/internal/models"
	"mth/internal/repository"
	"mth/pkg/log"
)

type districtService struct {
	districtRepo repository.District
	logger       *log.Logs
}

func InitDistrictService(districtRepo repository.District, logger *log.Logs) District {
	return districtService{
		districtRepo: districtRepo,
		logger:       logger,
	}
}

func (d districtService) GetByID(ctx context.Context, cityID int) ([]models.District, error) {
	districts, err := d.districtRepo.GetByCityID(ctx, cityID)
	if err != nil {
		d.logger.Error(err.Error())
		return []models.District{}, err
	}

	return districts, nil
}
