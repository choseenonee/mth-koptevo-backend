package repository

import (
	"context"
	"encoding/json"
	"github.com/jmoiron/sqlx"
	"mth/internal/models"
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

func (u userRepo) GetUser(ctx context.Context, login string) (int, string, error) {
	query := `SELECT id, password FROM users WHERE login = $1`

	var userID int
	var password string

	err := u.db.QueryRowContext(ctx, query, login).Scan(&userID, &password)
	if err != nil {
		return 0, "", customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ScanErr, Err: err})
	}

	return userID, password, nil
}

func (u userRepo) CreateUser(ctx context.Context, userCreate models.UserCreate) (int, error) {
	propertiesRaw, err := json.Marshal(userCreate.Properties)
	if err != nil {
		return 0, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.BindErr, Err: err})
	}

	tx, err := u.db.Beginx()
	if err != nil {
		return 0, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.TransactionErr, Err: err})
	}

	createRouteQuery := `INSERT INTO users (login, password, properties) VALUES ($1, $2, $3) RETURNING id;`

	var createdID int
	err = tx.QueryRowxContext(ctx, createRouteQuery, userCreate.Login, userCreate.Password, propertiesRaw).Scan(&createdID)
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
