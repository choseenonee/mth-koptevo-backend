package service

import (
	"context"
	"mth/internal/models"
	"mth/internal/models/swagger"
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
	GetAllWithFilter(ctx context.Context, filters swagger.Filters) ([]models.Place, error)
	GetByID(ctx context.Context, placeID int) (models.Place, error)
}

type District interface {
	GetByID(ctx context.Context, cityID int) ([]models.District, error)
}

type Route interface {
	Create(ctx context.Context, route models.RouteCreate) (int, error)
	GetByID(ctx context.Context, routeID int) (models.Route, error)
	GetAll(ctx context.Context, page int) ([]models.Route, error)
}

type Note interface {
	Create(ctx context.Context, noteCreate models.NoteCreate) (int, error)
	GetByIDs(ctx context.Context, userID int, placeID int) (models.Note, error)
	GetByID(ctx context.Context, noteID int) (models.Note, error)
	GetByUser(ctx context.Context, userID int) ([]models.Note, error)
	Update(ctx context.Context, noteUpd models.NoteUpdate) error
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
	// GetLikedByUser Places then Routes RAW
	GetLikedByUser(ctx context.Context, userID int) ([]models.Place, []models.RouteRaw, error)
	DeleteOnPlace(ctx context.Context, like models.Like) error
	DeleteOnRoute(ctx context.Context, like models.Like) error
}

type User interface {
	GetUser(ctx context.Context, login, password string) (int, error)
	GetProperties(ctx context.Context, userID int) (string, time.Time, interface{}, error)
	UpdateProperties(ctx context.Context, userID int, properties interface{}) error
	CreateUser(ctx context.Context, userCreate models.UserCreate) (int, error)
	CheckIn(ctx context.Context, cipher string, userID int) (string, error)
	ValidateHash(ctx context.Context, hash string) bool
	GetCheckedPlaces(ctx context.Context, userID int) ([]models.Place, error)
	GetChrono(ctx context.Context, userID int) (models.Chrono, error)
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
