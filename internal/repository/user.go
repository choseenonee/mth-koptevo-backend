package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/guregu/null/v5"
	"github.com/jmoiron/sqlx"
	"mth/internal/models"
	"mth/pkg/customerr"
	"time"
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

func (u userRepo) createRouteLog(ctx context.Context, query string, routeLog models.RouteLogWithOneTime) error {
	tx, err := u.db.Beginx()
	if err != nil {
		return customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.TransactionErr, Err: err})
	}

	_, err = tx.ExecContext(ctx, query, routeLog.UserID, routeLog.RouteID, routeLog.TimeStamp)
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

func (u userRepo) StartRoute(ctx context.Context, routeLog models.RouteLogWithOneTime) error {
	query := `INSERT INTO users_route_logs (user_id, route_id, start_time) VALUES ($1, $2, $3);`
	return u.createRouteLog(ctx, query, routeLog)
}

func (u userRepo) EndRoute(ctx context.Context, routeLog models.RouteLogWithOneTime) error {
	query := `UPDATE users_route_logs SET end_time = $3 WHERE user_id = $1 AND route_id = $2;`
	return u.createRouteLog(ctx, query, routeLog)
}

func (u userRepo) GetCheckedInPlaceIDs(ctx context.Context, userID int) ([]int, error) {
	query := `SELECT place_id FROM users_place_checkin WHERE user_id = $1`

	rows, err := u.db.QueryContext(ctx, query, userID)
	if err != nil {
		return []int{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ExecErr, Err: err})
	}

	var placeIDs []int
	for rows.Next() {
		var placeID int

		err = rows.Scan(&placeID)
		if err != nil {
			return []int{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ScanErr, Err: err})
		}

		placeIDs = append(placeIDs, placeID)
	}

	err = rows.Err()
	if err != nil {
		return []int{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.RowsErr, Err: err})
	}

	return placeIDs, nil
}

func (u userRepo) GetCheckInTimeStamp(ctx context.Context, userID, placeID int) (time.Time, error) {
	query := `SELECT timestamp FROM users_place_checkin WHERE user_id = $1 AND place_id = $2`

	var timeStamp time.Time
	err := u.db.QueryRowContext(ctx, query, userID, placeID).Scan(&timeStamp)
	if err != nil {
		return time.Time{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ExecErr, Err: err})
	}

	return timeStamp, nil
}

func (u userRepo) GetRouteLogs(ctx context.Context, userID int) ([]models.RouteLog, error) {
	query := `SELECT user_id, route_id, start_time, end_time FROM users_route_logs WHERE user_id = $1`

	rows, err := u.db.QueryContext(ctx, query, userID)
	if err != nil {
		return []models.RouteLog{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ExecErr, Err: err})
	}

	var routeLogs []models.RouteLog
	for rows.Next() {
		var routeLog models.RouteLog
		var endTime null.Time

		err = rows.Scan(&routeLog.UserID, &routeLog.RouteId, &routeLog.StartTime, &endTime)
		if err != nil {
			return []models.RouteLog{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ScanErr, Err: err})
		}

		if endTime.Valid {
			routeLog.EndTime = endTime.Time
		}

		routeLogs = append(routeLogs, routeLog)
	}

	err = rows.Err()
	if err != nil {
		return []models.RouteLog{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.RowsErr, Err: err})
	}

	return routeLogs, nil
}

func (u userRepo) GetProperties(ctx context.Context, userID int) (string, interface{}, error) {
	query := `SELECT login, properties FROM users WHERE id = $1`

	var propertiesRaw []byte
	var properties interface{}
	var login string
	err := u.db.QueryRowContext(ctx, query, userID).Scan(&login, &propertiesRaw)
	if err != nil {
		return "", nil, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ScanErr, Err: err})
	}

	err = json.Unmarshal(propertiesRaw, &properties)
	if err != nil {
		return "", nil, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.BindErr, Err: err})
	}

	return login, properties, nil
}

func (u userRepo) UpdateProperties(ctx context.Context, userID int, properties interface{}) error {
	query := `UPDATE users SET properties = $2 WHERE id = $1`

	tx, err := u.db.Beginx()
	if err != nil {
		return customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.TransactionErr, Err: err})
	}

	propertiesRaw, err := json.Marshal(properties)
	if err != nil {
		return customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.BindErr, Err: err})
	}

	res, err := tx.ExecContext(ctx, query, userID, propertiesRaw)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return customerr.ErrNormalizer(
				customerr.ErrorPair{Message: customerr.ScanErr, Err: err},
				customerr.ErrorPair{Message: customerr.RollbackErr, Err: rbErr},
			)
		}

		return customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ExecErr, Err: err})
	}

	count, err := res.RowsAffected()
	if count != 1 {
		if rbErr := tx.Rollback(); rbErr != nil {
			return customerr.ErrNormalizer(
				customerr.ErrorPair{Message: customerr.CountErr, Err: err},
				customerr.ErrorPair{Message: customerr.RollbackErr, Err: rbErr},
			)
		}

		return customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.CountErr, Err: fmt.Errorf("%v, user not found", count)})
	}

	err = tx.Commit()
	if err != nil {
		return customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.CommitErr, Err: err})
	}

	return nil
}
