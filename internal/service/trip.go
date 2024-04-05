package service

import (
	"context"
	"mth/internal/models"
	"mth/internal/repository"
	"mth/pkg/log"
)

type tripService struct {
	tripRepo repository.Trip
	logger   *log.Logs
}

func InitTripService(tripRepo repository.Trip, logger *log.Logs) Trip {
	return tripService{
		tripRepo: tripRepo,
		logger:   logger,
	}
}

func (t tripService) Create(ctx context.Context, tripCreate models.TripCreate) (int, error) {
	id, err := t.tripRepo.Create(ctx, tripCreate)
	if err != nil {
		t.logger.Error(err.Error())
		return 0, err
	}

	return id, err
}

func (t tripService) GetTripByID(ctx context.Context, tripID int) (models.Trip, error) {
	trip, err := t.tripRepo.GetTripByID(ctx, tripID)
	if err != nil {
		t.logger.Error(err.Error())
		return models.Trip{}, err
	}

	return trip, nil
}

func (t tripService) GetTripsByUser(ctx context.Context, userID int) ([]models.Trip, error) {
	trips, err := t.tripRepo.GetTripsByUser(ctx, userID)
	if err != nil {
		t.logger.Error(err.Error())
		return []models.Trip{}, err
	}

	return trips, nil
}

func (t tripService) AddRoute(ctx context.Context, tripID, routeID, day, position int) error {
	if err := t.tripRepo.AddRoute(ctx, tripID, routeID, day, position); err != nil {
		t.logger.Error(err.Error())
		return err
	}

	return nil
}

func (t tripService) AddPlace(ctx context.Context, tripID, placeID, day, position int) error {
	if err := t.tripRepo.AddPlace(ctx, tripID, placeID, day, position); err != nil {
		t.logger.Error(err.Error())
		return err
	}

	return nil
}

func (t tripService) ChangeRouteDay(ctx context.Context, tripID, routeID, day int) error {
	if err := t.tripRepo.ChangeRouteDay(ctx, tripID, routeID, day); err != nil {
		t.logger.Error(err.Error())
		return err
	}

	return nil
}

func (t tripService) ChangePlaceDay(ctx context.Context, tripID, placeID, day int) error {
	if err := t.tripRepo.ChangePlaceDay(ctx, tripID, placeID, day); err != nil {
		t.logger.Error(err.Error())
		return err
	}

	return nil
}

func (t tripService) ChangeRoutePosition(ctx context.Context, tripID, routeID, position int) error {
	if err := t.tripRepo.ChangeRoutePosition(ctx, tripID, routeID, position); err != nil {
		t.logger.Error(err.Error())
		return err
	}

	return nil
}

func (t tripService) ChangePlacePosition(ctx context.Context, tripID, placeID, position int) error {
	if err := t.tripRepo.ChangePlacePosition(ctx, tripID, placeID, position); err != nil {
		t.logger.Error(err.Error())
		return err
	}

	return nil
}

func (t tripService) DeleteRoute(ctx context.Context, tripID, routeID int) error {
	if err := t.tripRepo.DeleteRoute(ctx, tripID, routeID); err != nil {
		t.logger.Error(err.Error())
		return err
	}

	return nil
}

func (t tripService) DeletePlace(ctx context.Context, tripID, placeID int) error {
	if err := t.tripRepo.DeletePlace(ctx, tripID, placeID); err != nil {
		t.logger.Error(err.Error())
		return err
	}

	return nil
}
