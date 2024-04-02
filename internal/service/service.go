package service

import (
	"context"
	"mth/internal/models"
	"mth/internal/models/swagger"
)

type Tag interface {
	Create(ctx context.Context, tag models.TagCreate) (int, error)
	GetAll(ctx context.Context) ([]models.Tag, error)
}

type Review interface {
	CreateOnRoute(ctx context.Context, routeReview models.RouteReviewCreate) (int, error)
	CreateOnPlace(ctx context.Context, placeReview models.PlaceReviewCreate) (int, error)
	GetByAuthor(ctx context.Context, authorID int) ([]models.PlaceReview, []models.RouteReview, error)
	GetByRoute(ctx context.Context, routeID int) ([]models.RouteReview, error)
	GetByPlace(ctx context.Context, placeID int) ([]models.PlaceReview, error)
	UpdateOnPlace(ctx context.Context, reviewUpd models.ReviewUpdate) error
	UpdateOnRoute(ctx context.Context, reviewUpd models.ReviewUpdate) error
}

type Place interface {
	Create(ctx context.Context, placeCreate models.PlaceCreate) (int, error)
	GetAllWithFilter(ctx context.Context, filters swagger.Filters) ([]models.Place, error)
	GetByID(ctx context.Context, placeID int) (models.Place, error)
}

type District interface {
	GetByID(ctx context.Context, cityID int) ([]models.District, error)
}
