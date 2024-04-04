package repository

import (
	"context"
	"github.com/jmoiron/sqlx"
	"mth/pkg/customerr"
)

type userRepo struct {
	db *sqlx.DB
}

func InitUserRepo(db *sqlx.DB) User {
	return userRepo{
		db: db,
	}
}

func (u userRepo) CheckInPlace(ctx context.Context, userID, placeID int) error {
	tx, err := u.db.Beginx()
	if err != nil {
		return customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.TransactionErr, Err: err})
	}

	createCheckInQuery := `INSERT INTO users_place_checkin (user_id, place_id, timestamp) VALUES ($1, $2, current_timestamp);`

	_, err = tx.ExecContext(ctx, createCheckInQuery, userID, placeID)
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
