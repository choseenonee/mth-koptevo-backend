package repository

import (
	"context"
	"mth/internal/models"
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
	GetAllWithFilter(ctx context.Context, districtID int, cityID int, tagIDs []int, page int, name string, variety string) ([]models.Place, error)
	GetByID(ctx context.Context, placeID int) (models.Place, error)
}

type District interface {
	GetByCityID(ctx context.Context, cityID int) ([]models.District, error)
}

type Route interface {
	Create(ctx context.Context, route models.RouteCreate) (int, error)
	GetByID(ctx context.Context, routeID int) (models.RouteRaw, error)
	GetAll(ctx context.Context, page int) ([]models.RouteRaw, error)
}

type Note interface {
	Create(ctx context.Context, noteCreate models.NoteCreate) (int, error)
	GetByIDs(ctx context.Context, userID int, placeID int) (models.Note, error)
	GetByID(ctx context.Context, userID int) (models.Note, error)
	GetByUser(ctx context.Context, userID int) ([]models.Note, error)
	Update(ctx context.Context, noteUpd models.NoteUpdate) error
}

type Favourite interface {
	LikePlace(ctx context.Context, like models.Like) error
	LikeRoute(ctx context.Context, like models.Like) error
	// GetLikedByUser PlaceIDs then RouteIDs
	GetLikedByUser(ctx context.Context, userID int) ([]int, []int, error)
	DeleteOnPlace(ctx context.Context, like models.Like) error
	DeleteOnRoute(ctx context.Context, like models.Like) error
}
