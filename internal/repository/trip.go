package repository

import (
	"context"
	"encoding/json"
	"github.com/jmoiron/sqlx"
	"mth/internal/models"
	"mth/pkg/customerr"
)

type tripRepo struct {
	db *sqlx.DB
}

func InitTripRepo(db *sqlx.DB) Trip {
	return tripRepo{
		db: db,
	}
}

func (t tripRepo) Create(ctx context.Context, tripCreate models.TripCreate) (int, error) {
	tx, err := t.db.Beginx()
	if err != nil {
		return 0, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.TransactionErr, Err: err})
	}

	propertiesRaw, err := json.Marshal(tripCreate.Properties)
	if err != nil {
		return 0, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.BindErr, Err: err})
	}

	createTripQuery := `INSERT INTO trips (user_id, date_start, date_end, properties) VALUES ($1, $2, $3, $4) RETURNING id;`

	var createdID int
	err = tx.QueryRowxContext(ctx, createTripQuery, tripCreate.UserID, tripCreate.DateStart, tripCreate.DateEnd, propertiesRaw).
		Scan(&createdID)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return 0, customerr.ErrNormalizer(
				customerr.ErrorPair{Message: customerr.ScanErr, Err: err},
				customerr.ErrorPair{Message: customerr.RollbackErr, Err: rbErr},
			)
		}

		return 0, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ScanErr, Err: err})
	}

	createTripPlaceQuery := `INSERT INTO trip_places (trip_id, day, position, place_id) VALUES ($1, $2, $3, $4);`
	for _, place := range tripCreate.Places {
		_, err := tx.ExecContext(ctx, createTripPlaceQuery, createdID, place.Day, place.Position, place.EntityID)
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				return 0, customerr.ErrNormalizer(
					customerr.ErrorPair{Message: customerr.ScanErr, Err: err},
					customerr.ErrorPair{Message: customerr.RollbackErr, Err: rbErr},
				)
			}

			return 0, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ScanErr, Err: err})
		}
	}

	createTripRouteQuery := `INSERT INTO trip_routes (trip_id, day, position, route_id) VALUES ($1, $2, $3, $4);`
	for _, route := range tripCreate.Routes {
		_, err := tx.ExecContext(ctx, createTripRouteQuery, createdID, route.Day, route.Position, route.EntityID)
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				return 0, customerr.ErrNormalizer(
					customerr.ErrorPair{Message: customerr.ScanErr, Err: err},
					customerr.ErrorPair{Message: customerr.RollbackErr, Err: rbErr},
				)
			}

			return 0, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ScanErr, Err: err})
		}
	}

	if err = tx.Commit(); err != nil {
		return 0, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.CommitErr, Err: err})
	}

	return createdID, nil
}

func (t tripRepo) GetTripByID(ctx context.Context, tripID int) (models.Trip, error) {
	query := `SELECT t.id, t.user_id, t.date_start, t.date_end, t.properties
				FROM trips t
				WHERE t.id = $1`

	rows, err := t.db.QueryContext(ctx, query, tripID)
	if err != nil {
		return models.Trip{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ExecErr, Err: err})
	}

	var trip models.Trip
	var propertiesRaw []byte
	for rows.Next() {
		err = rows.Scan(&trip.ID, &trip.UserID, &trip.DateStart, &trip.DateEnd, &propertiesRaw)
		if err != nil {
			return models.Trip{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ScanErr, Err: err})
		}

		err = json.Unmarshal(propertiesRaw, &trip.Properties)
		if err != nil {
			return models.Trip{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.BindErr, Err: err})
		}
	}

	err = rows.Err()
	if err != nil {
		return models.Trip{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.RowsErr, Err: err})
	}

	placeIDsQuery := `SELECT tp.place_id, tp.day, tp.position
				FROM trip_places tp
				WHERE tp.trip_id = $1`

	rows, err = t.db.QueryContext(ctx, placeIDsQuery, tripID)
	if err != nil {
		return models.Trip{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ExecErr, Err: err})
	}

	var place models.EntityWithDayAndPosition
	for rows.Next() {
		err = rows.Scan(&place.EntityID, &place.Day, &place.Position)
		if err != nil {
			return models.Trip{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ScanErr, Err: err})
		}

		trip.Places = append(trip.Places, place)
	}

	err = rows.Err()
	if err != nil {
		return models.Trip{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.RowsErr, Err: err})
	}

	routeIDsQuery := `SELECT tr.route_id, tr.day, tr.position
				FROM trip_routes tr
				WHERE tr.trip_id = $1`

	rows, err = t.db.QueryContext(ctx, routeIDsQuery, tripID)
	if err != nil {
		return models.Trip{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ExecErr, Err: err})
	}

	var route models.EntityWithDayAndPosition
	for rows.Next() {
		err = rows.Scan(&route.EntityID, &route.Day, &route.Position)
		if err != nil {
			return models.Trip{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ScanErr, Err: err})
		}

		trip.Routes = append(trip.Routes, route)
	}

	err = rows.Err()
	if err != nil {
		return models.Trip{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.RowsErr, Err: err})
	}

	return trip, nil
}

func (t tripRepo) GetTripsByUser(ctx context.Context, userID int) ([]models.Trip, error) {
	query := `SELECT t.id FROM trips t WHERE t.user_id = $1`

	rows, err := t.db.QueryContext(ctx, query, userID)
	if err != nil {
		return []models.Trip{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ExecErr, Err: err})
	}

	var tripIDs []int
	for rows.Next() {
		var tripID int
		err = rows.Scan(&tripID)
		if err != nil {
			return []models.Trip{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ScanErr, Err: err})
		}

		tripIDs = append(tripIDs, tripID)
	}

	err = rows.Err()
	if err != nil {
		return []models.Trip{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.RowsErr, Err: err})
	}

	var trips []models.Trip
	for _, tripID := range tripIDs {
		trip, err := t.GetTripByID(ctx, tripID)
		if err != nil {
			return []models.Trip{}, err
		}

		trips = append(trips, trip)
	}

	return trips, nil
}

func (t tripRepo) AddRoute(ctx context.Context, tripID, routeID, day, position int) error {
	//TODO implement me
	panic("implement me")
}

func (t tripRepo) AddPlace(ctx context.Context, tripID, placeID, day, position int) error {
	//TODO implement me
	panic("implement me")
}

func (t tripRepo) ChangeRouteDay(ctx context.Context, tripID, placeID, day int) error {
	//TODO implement me
	panic("implement me")
}

func (t tripRepo) ChangePlaceDay(ctx context.Context, tripID, placeID, day int) error {
	//TODO implement me
	panic("implement me")
}

func (t tripRepo) ChangeRoutePosition(ctx context.Context, tripID, routeID int) error {
	//TODO implement me
	panic("implement me")
}

func (t tripRepo) ChangePlacePosition(ctx context.Context, tripID, placeID int) error {
	//TODO implement me
	panic("implement me")
}

func (t tripRepo) DeleteRoute(ctx context.Context, tripID, routeID int) error {
	//TODO implement me
	panic("implement me")
}

func (t tripRepo) DeletePlace(ctx context.Context, tripID, placeID int) error {
	//TODO implement me
	panic("implement me")
}
