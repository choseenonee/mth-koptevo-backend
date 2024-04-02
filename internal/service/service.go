package service

import (
	"context"
	"mth/internal/models"
)

type Tag interface {
	Create(ctx context.Context, tag models.TagCreate) (int, error)
	GetAll(ctx context.Context) ([]models.Tag, error)
}

type Review interface {
	CreateRoute(ctx context.Context, routeReview models.RouteReviewCreate) (int, error)
	CreatePlace(ctx context.Context, placeReview models.PlaceReviewCreate) (int, error)
	GetByAuthor(ctx context.Context, authorID int) ([]models.PlaceReview, []models.RouteReview, error)
	GetByRoute(ctx context.Context, routeID int) ([]models.RouteReview, error)
	GetByPlace(ctx context.Context, placeID int) ([]models.PlaceReview, error)
	UpdatePlace(ctx context.Context, reviewUpd models.ReviewUpdate) error
	UpdateRoute(ctx context.Context, reviewUpd models.ReviewUpdate) error
}
