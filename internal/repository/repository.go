package repository

import (
	"context"
	"mth/internal/models"
	"time"
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
	Update(ctx context.Context, noteUpd models.NoteCreate) error
}

type Companions interface {
	CreatePlaceCompanions(ctx context.Context, companion models.CompanionsPlaceCreate) error
	CreateRouteCompanions(ctx context.Context, companion models.CompanionsRouteCreate) error
	// GetByUser сначала places, затем routes
	GetByUser(ctx context.Context, userID int) ([]models.CompanionsPlace, []models.CompanionsRoute, error)
	GetCompanionsPlace(ctx context.Context, filters models.CompanionsFilters) ([]models.CompanionsPlace, error)
	GetCompanionsRoute(ctx context.Context, filters models.CompanionsFilters) ([]models.CompanionsRoute, error)
	DeleteCompanionsPlace(ctx context.Context, id int) error
	DeleteCompanionsRoute(ctx context.Context, id int) error
}

type Favourite interface {
	LikePlace(ctx context.Context, like models.Like) error
	LikeRoute(ctx context.Context, like models.Like) error
	// GetLikedByUser PlaceIDs then RouteIDs
	GetLikedByUser(ctx context.Context, userID int) ([]int, []int, error)
	GetPlaceTimestamp(ctx context.Context, userID, placeID int) (time.Time, error)
	GetRouteTimestamp(ctx context.Context, userID, routeID int) (time.Time, error)
	DeleteOnPlace(ctx context.Context, like models.Like) error
	DeleteOnRoute(ctx context.Context, like models.Like) error
}

// User TODO: из get route logs сделать позже get_chrono service
type User interface {
	GetUser(ctx context.Context, login string) (int, string, error)
	GetProperties(ctx context.Context, userID int) (string, interface{}, error)
	UpdateProperties(ctx context.Context, userID int, properties interface{}) error
	CreateUser(ctx context.Context, userCreate models.UserCreate) (int, error)
	CheckInPlace(ctx context.Context, userID, placeID int) error
	GetCheckedInPlaceIDs(ctx context.Context, userID int) ([]int, error)
	GetRouteLogs(ctx context.Context, userID int) ([]models.RouteLog, error)
	StartRoute(ctx context.Context, routeLog models.RouteLogWithOneTime) error
	EndRoute(ctx context.Context, routeLog models.RouteLogWithOneTime) error
	GetCheckInTimeStamp(ctx context.Context, userID, placeID int) (time.Time, error)
}

type Trip interface {
	Create(ctx context.Context, tripCreate models.TripCreate) (int, error)
	GetTripByID(ctx context.Context, tripID int) (models.Trip, error)
	GetTripsByUser(ctx context.Context, userID int) ([]models.Trip, error)
	AddRoute(ctx context.Context, tripID, routeID, day, position int) error
	AddPlace(ctx context.Context, tripID, placeID, day, position int) error
	ChangeRouteDay(ctx context.Context, tripID, routeID, day int) error
	ChangePlaceDay(ctx context.Context, tripID, placeID, day int) error
	ChangeRoutePosition(ctx context.Context, tripID, routeID, position int) error
	ChangePlacePosition(ctx context.Context, tripID, placeID, position int) error
	DeleteRoute(ctx context.Context, tripID, routeID int) error
	DeletePlace(ctx context.Context, tripID, placeID int) error
}
