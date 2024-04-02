package service

import (
	"context"
	"mth/internal/models"
	"mth/internal/repository"
	"mth/pkg/log"
)

type reviewService struct {
	reviewRepo repository.Review
	logger     *log.Logs
}

func InitReviewService(reviewRepo repository.Review, logger *log.Logs) Review {
	return reviewService{
		reviewRepo: reviewRepo,
		logger:     logger,
	}
}

func (r reviewService) CreateOnRoute(ctx context.Context, routeReview models.RouteReviewCreate) (int, error) {
	id, err := r.reviewRepo.CreateOnRoute(ctx, routeReview)
	if err != nil {
		r.logger.Error(err.Error())
		return 0, err
	}

	return id, nil
}

func (r reviewService) CreateOnPlace(ctx context.Context, placeReview models.PlaceReviewCreate) (int, error) {
	id, err := r.reviewRepo.CreateOnPlace(ctx, placeReview)
	if err != nil {
		r.logger.Error(err.Error())
		return 0, err
	}

	return id, nil
}

func (r reviewService) GetByAuthor(ctx context.Context, authorID int) ([]models.PlaceReview, []models.RouteReview, error) {
	placeReviews, routeReviews, err := r.reviewRepo.GetByAuthor(ctx, authorID)
	if err != nil {
		r.logger.Error(err.Error())
		return []models.PlaceReview{}, []models.RouteReview{}, err
	}

	return placeReviews, routeReviews, nil
}

func (r reviewService) GetByRoute(ctx context.Context, routeID int) ([]models.RouteReview, error) {
	routeReviews, err := r.reviewRepo.GetByRoute(ctx, routeID)
	if err != nil {
		r.logger.Error(err.Error())
		return []models.RouteReview{}, err
	}

	return routeReviews, nil
}

func (r reviewService) GetByPlace(ctx context.Context, placeID int) ([]models.PlaceReview, error) {
	placeReviews, err := r.reviewRepo.GetByPlace(ctx, placeID)
	if err != nil {
		r.logger.Error(err.Error())
		return []models.PlaceReview{}, err
	}

	return placeReviews, nil
}

func (r reviewService) UpdateOnPlace(ctx context.Context, reviewUpd models.ReviewUpdate) error {
	err := r.reviewRepo.UpdateOnPlace(ctx, reviewUpd)
	if err != nil {
		r.logger.Error(err.Error())
		return err
	}

	return nil
}

func (r reviewService) UpdateOnRoute(ctx context.Context, reviewUpd models.ReviewUpdate) error {
	err := r.reviewRepo.UpdateOnRoute(ctx, reviewUpd)
	if err != nil {
		r.logger.Error(err.Error())
		return err
	}

	return nil
}
