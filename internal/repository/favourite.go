package repository

import (
	"context"
	"github.com/jmoiron/sqlx"
	"mth/internal/models"
	"mth/pkg/customerr"
)

type favouriteRepo struct {
	db *sqlx.DB
}

func InitFavouriteRepo(db *sqlx.DB) Favourite {
	return favouriteRepo{
		db: db,
	}
}

func (f favouriteRepo) LikePlace(ctx context.Context, like models.Like) error {
	tx, err := f.db.Beginx()
	if err != nil {
		return customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.TransactionErr, Err: err})
	}

	query := `INSERT INTO users_favourite_places (user_id, place_id, timestamp) VALUES ($1, $2, current_timestamp);`

	_, err = tx.ExecContext(ctx, query, like.UserID, like.EntityID)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return customerr.ErrNormalizer(
				customerr.ErrorPair{Message: customerr.ScanErr, Err: err},
				customerr.ErrorPair{Message: customerr.RollbackErr, Err: rbErr},
			)
		}

		return customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ScanErr, Err: err})
	}

	if err = tx.Commit(); err != nil {
		return customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.CommitErr, Err: err})
	}

	return nil
}

func (f favouriteRepo) LikeRoute(ctx context.Context, like models.Like) error {
	tx, err := f.db.Beginx()
	if err != nil {
		return customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.TransactionErr, Err: err})
	}

	query := `INSERT INTO users_favourite_routes (user_id, route_id, timestamp) VALUES ($1, $2, current_timestamp);`

	_, err = tx.ExecContext(ctx, query, like.UserID, like.EntityID)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return customerr.ErrNormalizer(
				customerr.ErrorPair{Message: customerr.ScanErr, Err: err},
				customerr.ErrorPair{Message: customerr.RollbackErr, Err: rbErr},
			)
		}

		return customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ScanErr, Err: err})
	}

	if err = tx.Commit(); err != nil {
		return customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.CommitErr, Err: err})
	}

	return nil
}

func (f favouriteRepo) GetLikedByUser(ctx context.Context, userID int) ([]int, []int, error) {
	routesQuery := `SELECT r.id FROM users_favourite_routes ufr
					LEFT JOIN routes r on ufr.route_id = r.id
					WHERE ufr.user_id = $1;`

	routesRows, err := f.db.QueryContext(ctx, routesQuery, userID)
	if err != nil {
		return []int{}, []int{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ExecErr, Err: err})
	}

	var routeIDs []int

	for routesRows.Next() {
		var routeID int

		err = routesRows.Scan(&routeID)
		if err != nil {
			return []int{}, []int{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ScanErr, Err: err})
		}

		routeIDs = append(routeIDs, routeID)
	}

	err = routesRows.Err()
	if err != nil {
		return []int{}, []int{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.RowsErr, Err: err})
	}

	placesQuery := `SELECT p.id FROM users_favourite_places ufp
					LEFT JOIN places p on ufp.place_id = p.id
					WHERE ufp.user_id = $1;`

	placeRows, err := f.db.QueryContext(ctx, placesQuery, userID)
	if err != nil {
		return []int{}, []int{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ExecErr, Err: err})
	}

	var placeIDs []int

	for placeRows.Next() {
		var placeID int

		err = placeRows.Scan(&placeID)
		if err != nil {
			return []int{}, []int{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ScanErr, Err: err})
		}

		placeIDs = append(placeIDs, placeID)
	}

	err = placeRows.Err()
	if err != nil {
		return []int{}, []int{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.RowsErr, Err: err})
	}

	return placeIDs, routeIDs, nil
}

func (f favouriteRepo) delete(ctx context.Context, query string, userID, entityID int) error {
	tx, err := f.db.Beginx()
	if err != nil {
		return customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.TransactionErr, Err: err})
	}

	_, err = tx.ExecContext(ctx, query, userID, entityID)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return customerr.ErrNormalizer(
				customerr.ErrorPair{Message: customerr.ScanErr, Err: err},
				customerr.ErrorPair{Message: customerr.RollbackErr, Err: rbErr},
			)
		}

		return customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ScanErr, Err: err})
	}

	if err = tx.Commit(); err != nil {
		return customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.CommitErr, Err: err})
	}

	return nil
}

func (f favouriteRepo) DeleteOnPlace(ctx context.Context, like models.Like) error {
	query := `DELETE FROM users_favourite_places WHERE user_id = $1 AND place_id = $2;`
	return f.delete(ctx, query, like.UserID, like.EntityID)
}

func (f favouriteRepo) DeleteOnRoute(ctx context.Context, like models.Like) error {
	query := `DELETE FROM users_favourite_routes WHERE user_id = $1 AND route_id = $2;`
	return f.delete(ctx, query, like.UserID, like.EntityID)
}
