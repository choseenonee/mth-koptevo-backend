package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/jmoiron/sqlx"
	"mth/internal/models"
	"mth/pkg/customerr"
	"time"
)

type reviewRepo struct {
	db *sqlx.DB
}

func InitReviewRepo(db *sqlx.DB) Review {
	return reviewRepo{
		db: db,
	}
}

type reviewCreate struct {
	AuthorID   int
	EntityID   int
	Properties interface{}
	Mark       float32
	TimeStamp  time.Time
	_          struct{}
}

type reviewGet struct {
	ID int
	reviewCreate
}

func (r reviewRepo) create(ctx context.Context, query string, review reviewCreate) (int, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return 0, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.TransactionErr, Err: err})
	}

	jsonProperties, err := json.Marshal(review.Properties)
	if err != nil {
		return 0, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.BindErr, Err: err})
	}

	var createdID int
	err = tx.QueryRowxContext(ctx, query, review.EntityID, review.AuthorID, jsonProperties, review.Mark).Scan(&createdID)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return 0, customerr.ErrNormalizer(
				customerr.ErrorPair{Message: customerr.ScanErr, Err: err},
				customerr.ErrorPair{Message: customerr.RollbackErr, Err: rbErr},
			)
		}

		return 0, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ScanErr, Err: err})
	}

	if err = tx.Commit(); err != nil {
		return 0, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.CommitErr, Err: err})
	}

	return createdID, nil
}

func (r reviewRepo) get(ctx context.Context, query string, authorID int, entityID int) ([]reviewGet, error) {
	var rows *sql.Rows
	var err error

	switch authorID {
	case 0:
		rows, err = r.db.QueryContext(ctx, query, entityID)
	default:
		rows, err = r.db.QueryContext(ctx, query, authorID)
	}

	if err != nil {
		return []reviewGet{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ExecErr, Err: err})
	}

	var reviews []reviewGet

	for rows.Next() {
		var propertiesRaw []byte
		var review reviewGet
		err := rows.Scan(&review.ID, &review.EntityID, &review.AuthorID, &propertiesRaw, &review.Mark, &review.TimeStamp)
		if err != nil {
			return []reviewGet{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ScanErr, Err: nil})
		}

		err = json.Unmarshal(propertiesRaw, &review.Properties)
		if err != nil {
			return []reviewGet{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.BindErr, Err: err})
		}

		reviews = append(reviews, review)
	}

	err = rows.Err()
	if err != nil {
		return []reviewGet{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.RowsErr, Err: nil})
	}

	return reviews, nil
}

func (r reviewRepo) update(ctx context.Context, query string, reviewUpd models.ReviewUpdate) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.TransactionErr, Err: err})
	}

	jsonProperties, err := json.Marshal(reviewUpd.Properties)
	if err != nil {
		return customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.BindErr, Err: err})
	}

	res, err := tx.ExecContext(ctx, query, reviewUpd.ID, jsonProperties, reviewUpd.Mark)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return customerr.ErrNormalizer(
				customerr.ErrorPair{Message: customerr.ScanErr, Err: err},
				customerr.ErrorPair{Message: customerr.RollbackErr, Err: rbErr},
			)
		}

		return customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ScanErr, Err: err})
	}
	count, err := res.RowsAffected()
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return customerr.ErrNormalizer(
				customerr.ErrorPair{Message: customerr.ScanErr, Err: err},
				customerr.ErrorPair{Message: customerr.RollbackErr, Err: rbErr},
			)
		}

		return customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ScanErr, Err: err})
	}
	if count != 1 {
		if rbErr := tx.Rollback(); rbErr != nil {
			return customerr.ErrNormalizer(
				customerr.ErrorPair{Message: customerr.CountErr, Err: err},
				customerr.ErrorPair{Message: customerr.RollbackErr, Err: rbErr},
			)
		}
		return customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.CountErr, Err: err})
	}

	err = tx.Commit()
	if err != nil {
		return customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.CommitErr, Err: err})
	}

	return nil
}

func (r reviewRepo) CreateOnRoute(ctx context.Context, routeReview models.RouteReviewCreate) (int, error) {
	createRouteReviewQuery := `INSERT INTO route_reviews (route_id, author_id, properties, mark, timestamp) VALUES ($1, $2, $3, $4, current_timestamp) RETURNING id`
	review := reviewCreate{
		AuthorID:   routeReview.AuthorID,
		EntityID:   routeReview.RouteID,
		Properties: routeReview.Properties,
		Mark:       routeReview.Mark,
	}
	return r.create(ctx, createRouteReviewQuery, review)
}

func (r reviewRepo) CreateOnPlace(ctx context.Context, placeReview models.PlaceReviewCreate) (int, error) {
	createRouteReviewQuery := `INSERT INTO places_reviews (place_id, author_id, properties, mark, timestamp) VALUES ($1, $2, $3, $4, current_timestamp) RETURNING id;`
	review := reviewCreate{
		AuthorID:   placeReview.AuthorID,
		EntityID:   placeReview.PlaceID,
		Properties: placeReview.Properties,
		Mark:       placeReview.Mark,
	}
	return r.create(ctx, createRouteReviewQuery, review)
}

func (r reviewRepo) GetByAuthor(ctx context.Context, authorID int) ([]models.PlaceReview, []models.RouteReview, error) {
	getRouteReviewByUser := `SELECT id, route_id, author_id, properties, mark, timestamp FROM route_reviews WHERE author_id = $1;`
	reviews, err := r.get(ctx, getRouteReviewByUser, authorID, 0)
	if err != nil {
		return []models.PlaceReview{}, []models.RouteReview{}, err
	}

	routeReviews := make([]models.RouteReview, len(reviews))
	for i := range reviews {
		routeReviews[i] = models.RouteReview{
			ID: reviews[i].ID,
			RouteReviewCreate: models.RouteReviewCreate{
				RouteID: reviews[i].EntityID,
				ReviewBase: models.ReviewBase{
					AuthorID:   reviews[i].AuthorID,
					Properties: reviews[i].Properties,
					Mark:       reviews[i].Mark,
					TimeStamp:  reviews[i].TimeStamp,
				},
			},
		}
	}

	getPlaceReviewsByAuthorQuery := `SELECT id, place_id, author_id, properties, mark, timestamp FROM places_reviews WHERE author_id = $1;`
	reviews, err = r.get(ctx, getPlaceReviewsByAuthorQuery, authorID, 0)
	if err != nil {
		return []models.PlaceReview{}, []models.RouteReview{}, err
	}

	placeReviews := make([]models.PlaceReview, len(reviews))
	for i := range reviews {
		placeReviews[i] = models.PlaceReview{
			ID: reviews[i].ID,
			PlaceReviewCreate: models.PlaceReviewCreate{
				PlaceID: reviews[i].EntityID,
				ReviewBase: models.ReviewBase{
					AuthorID:   reviews[i].AuthorID,
					Properties: reviews[i].Properties,
					Mark:       reviews[i].Mark,
					TimeStamp:  reviews[i].TimeStamp,
				},
			},
		}
	}

	return placeReviews, routeReviews, nil
}

func (r reviewRepo) GetByRoute(ctx context.Context, routeID int) ([]models.RouteReview, error) {
	getByRouteQuery := `SELECT id, route_id, author_id, properties, mark, timestamp FROM route_reviews WHERE route_id = $1;`
	reviews, err := r.get(ctx, getByRouteQuery, 0, routeID)
	if err != nil {
		return []models.RouteReview{}, err
	}

	routeReviews := make([]models.RouteReview, len(reviews))
	for i := range reviews {
		routeReviews[i] = models.RouteReview{
			ID: reviews[i].ID,
			RouteReviewCreate: models.RouteReviewCreate{
				RouteID: reviews[i].EntityID,
				ReviewBase: models.ReviewBase{
					AuthorID:   reviews[i].AuthorID,
					Properties: reviews[i].Properties,
					Mark:       reviews[i].Mark,
					TimeStamp:  reviews[i].TimeStamp,
				},
			},
		}
	}

	return routeReviews, nil
}

func (r reviewRepo) GetByPlace(ctx context.Context, placeID int) ([]models.PlaceReview, error) {
	getByPlaceQuery := `SELECT id, place_id, author_id, properties, mark, timestamp FROM places_reviews WHERE place_id = $1;`
	reviews, err := r.get(ctx, getByPlaceQuery, 0, placeID)
	if err != nil {
		return []models.PlaceReview{}, err
	}

	placeReviews := make([]models.PlaceReview, len(reviews))
	for i := range reviews {
		placeReviews[i] = models.PlaceReview{
			ID: reviews[i].ID,
			PlaceReviewCreate: models.PlaceReviewCreate{
				PlaceID: reviews[i].EntityID,
				ReviewBase: models.ReviewBase{
					AuthorID:   reviews[i].AuthorID,
					Properties: reviews[i].Properties,
					Mark:       reviews[i].Mark,
					TimeStamp:  reviews[i].TimeStamp,
				},
			},
		}
	}

	return placeReviews, nil
}

func (r reviewRepo) UpdateOnPlace(ctx context.Context, reviewUpd models.ReviewUpdate) error {
	updatePlaceReviewQuery := `UPDATE places_reviews SET properties = $2, mark = $3 WHERE id = $1;`
	return r.update(ctx, updatePlaceReviewQuery, reviewUpd)
}

func (r reviewRepo) UpdateOnRoute(ctx context.Context, reviewUpd models.ReviewUpdate) error {
	updateRouteReviewQuery := `UPDATE route_reviews SET properties = $2, mark = $3 WHERE id = $1;`
	return r.update(ctx, updateRouteReviewQuery, reviewUpd)
}
