package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"mth/internal/models"
	"mth/pkg/config"
	"mth/pkg/customerr"
)

type companionsRepo struct {
	db *sqlx.DB
}

func InitCompanionsRepo(db *sqlx.DB) Companions {
	return companionsRepo{db: db}
}

func (c companionsRepo) CreatePlaceCompanions(ctx context.Context, companion models.CompanionsPlaceCreate) error {
	tx, err := c.db.Beginx()
	if err != nil {
		return customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.TransactionErr, Err: err})
	}

	createQuery := `INSERT INTO companions_places (user_id, place_id, date_from, date_to)
					VALUES ($1, $2, $3, $4);`

	res, err := tx.ExecContext(ctx, createQuery, companion.UserID, companion.PlaceID, companion.DateFrom, companion.DateTo)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return customerr.ErrNormalizer(
				customerr.ErrorPair{Message: customerr.ExecErr, Err: err},
				customerr.ErrorPair{Message: customerr.RollbackErr, Err: rbErr},
			)
		}
		return customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ExecErr, Err: err})
	}
	count, err := res.RowsAffected()
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return customerr.ErrNormalizer(
				customerr.ErrorPair{Message: customerr.RowsErr, Err: err},
				customerr.ErrorPair{Message: customerr.RollbackErr, Err: rbErr},
			)
		}
		return customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.RowsErr, Err: err})
	}
	if count != 1 {
		if rbErr := tx.Rollback(); rbErr != nil {
			return customerr.ErrNormalizer(
				customerr.ErrorPair{Message: customerr.RowsErr, Err: fmt.Errorf(customerr.CountErr, count)},
				customerr.ErrorPair{Message: customerr.RollbackErr, Err: rbErr},
			)
		}
		return customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.RowsErr, Err: fmt.Errorf(customerr.CountErr, count)})
	}

	if err = tx.Commit(); err != nil {
		return customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.CommitErr, Err: err})
	}

	return nil
}

func (c companionsRepo) CreateRouteCompanions(ctx context.Context, companion models.CompanionsRouteCreate) error {
	tx, err := c.db.Beginx()
	if err != nil {
		return customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.TransactionErr, Err: err})
	}

	createQuery := `INSERT INTO companions_routes (user_id, route_id, date_from, date_to)
					VALUES ($1, $2, $3, $4);`

	res, err := tx.ExecContext(ctx, createQuery, companion.UserID, companion.RouteID, companion.DateFrom, companion.DateTo)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return customerr.ErrNormalizer(
				customerr.ErrorPair{Message: customerr.ExecErr, Err: err},
				customerr.ErrorPair{Message: customerr.RollbackErr, Err: rbErr},
			)
		}
		return customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ExecErr, Err: err})
	}
	count, err := res.RowsAffected()
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return customerr.ErrNormalizer(
				customerr.ErrorPair{Message: customerr.RowsErr, Err: err},
				customerr.ErrorPair{Message: customerr.RollbackErr, Err: rbErr},
			)
		}
		return customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.RowsErr, Err: err})
	}
	if count != 1 {
		if rbErr := tx.Rollback(); rbErr != nil {
			return customerr.ErrNormalizer(
				customerr.ErrorPair{Message: customerr.RowsErr, Err: fmt.Errorf(customerr.CountErr, count)},
				customerr.ErrorPair{Message: customerr.RollbackErr, Err: rbErr},
			)
		}
		return customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.RowsErr, Err: fmt.Errorf(customerr.CountErr, count)})
	}

	if err = tx.Commit(); err != nil {
		return customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.CommitErr, Err: err})
	}

	return nil
}

func (c companionsRepo) GetByUser(ctx context.Context, userID int) ([]models.CompanionsPlace, []models.CompanionsRoute, error) {
	placeQuery := `SELECT c.date_from, c.date_to, p.name, p.properties, city.name, u.id, p.id, c.id
					FROM companions_places c
					LEFT JOIN users u ON c.user_id = u.id
					LEFT JOIN places p ON c.place_id = p.id
					LEFT JOIN city ON p.city_id = city.id
					WHERE u.id = $1`

	placeRows, err := c.db.QueryxContext(ctx, placeQuery, userID)
	if err != nil {
		return []models.CompanionsPlace{}, []models.CompanionsRoute{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.QueryErr, Err: err})
	}

	var places []models.CompanionsPlace

	for placeRows.Next() {
		var place models.CompanionsPlace
		var propertiesRaw []byte

		err := placeRows.Scan(&place.DateFrom, &place.DateTo, &place.PlaceName,
			&propertiesRaw, &place.CityName, &place.UserID, &place.PlaceID, &place.ID)
		if err != nil {
			return []models.CompanionsPlace{}, []models.CompanionsRoute{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ScanErr, Err: err})
		}

		err = json.Unmarshal(propertiesRaw, &place.PlaceProperties)
		if err != nil {
			return []models.CompanionsPlace{}, []models.CompanionsRoute{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.BindErr, Err: err})
		}

		places = append(places, place)
	}

	err = placeRows.Err()
	if err != nil {
		return []models.CompanionsPlace{}, []models.CompanionsRoute{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ScanErr, Err: err})
	}

	routeQuery := `SELECT c.date_from, c.date_to, p.name, p.properties, city.name, u.id, p.id, c.id
					FROM companions_routes c
					LEFT JOIN users u ON c.user_id = u.id
					LEFT JOIN places p ON c.route_id = p.id
					LEFT JOIN city ON p.city_id = city.id
					WHERE u.id = $1`

	routeRows, err := c.db.QueryxContext(ctx, routeQuery, userID)
	if err != nil {
		return []models.CompanionsPlace{}, []models.CompanionsRoute{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.QueryErr, Err: err})
	}

	var routes []models.CompanionsRoute

	for routeRows.Next() {
		var route models.CompanionsRoute
		var propertiesRaw []byte

		err := routeRows.Scan(&route.DateFrom, &route.DateTo, &route.RouteName,
			&propertiesRaw, &route.CityName, &route.UserID, &route.RouteID, &route.ID)
		if err != nil {
			return []models.CompanionsPlace{}, []models.CompanionsRoute{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ScanErr, Err: err})
		}

		err = json.Unmarshal(propertiesRaw, &route.RouteProperties)
		if err != nil {
			return []models.CompanionsPlace{}, []models.CompanionsRoute{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ScanErr, Err: err})
		}

		routes = append(routes, route)
	}

	err = placeRows.Err()
	if err != nil {
		return []models.CompanionsPlace{}, []models.CompanionsRoute{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ScanErr, Err: err})
	}

	return places, routes, nil
}

func (c companionsRepo) GetCompanionsPlace(ctx context.Context, filters models.CompanionsFilters) ([]models.CompanionsPlace, error) {
	var companions []models.CompanionsPlace

	selectQuery := `SELECT c.date_from, c.date_to, u.properties, p.name, p.properties, city.name, u.id, p.id, c.id
					FROM companions_places c
					LEFT JOIN users u ON c.user_id = u.id
					LEFT JOIN places p ON c.place_id = p.id
					LEFT JOIN city ON p.city_id = city.id
					WHERE NOT (c.date_from > $2 OR c.date_to < $1) AND c.place_id = $3
					OFFSET $4 LIMIT $5;`

	rows, err := c.db.QueryxContext(ctx, selectQuery, filters.DateFrom, filters.DateTo, filters.EntityID,
		filters.Page*viper.GetInt(config.CompanionsOnPage), viper.GetInt(config.CompanionsOnPage))
	if err != nil {
		return []models.CompanionsPlace{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.QueryErr, Err: err})
	}

	for rows.Next() {
		var companion models.CompanionsPlace
		var propertiesRaw []byte
		var userPropertiesRaw []byte

		err := rows.Scan(&companion.DateFrom, &companion.DateTo, &userPropertiesRaw, &companion.PlaceName,
			&propertiesRaw, &companion.CityName, &companion.UserID, &companion.PlaceID, &companion.ID)
		if err != nil {
			return []models.CompanionsPlace{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ScanErr, Err: err})
		}

		err = json.Unmarshal(userPropertiesRaw, &companion.UserProperties)
		if err != nil {
			return []models.CompanionsPlace{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.BindErr, Err: err})
		}

		err = json.Unmarshal(propertiesRaw, &companion.PlaceProperties)
		if err != nil {
			return []models.CompanionsPlace{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.BindErr, Err: err})
		}

		companions = append(companions, companion)
	}

	return companions, nil
}

func (c companionsRepo) GetCompanionsRoute(ctx context.Context, filters models.CompanionsFilters) ([]models.CompanionsRoute, error) {
	var companions []models.CompanionsRoute

	selectQuery := `SELECT c.date_from, c.date_to, u.properties, r.name, r.price, r.properties, city.name, u.id, r.id, c.id
					FROM companions_routes c
					LEFT JOIN users u ON c.user_id = u.id
					LEFT JOIN routes r ON c.route_id = r.id
					LEFT JOIN city ON r.city_id = city.id
					WHERE NOT (c.date_from > $2 OR c.date_to < $1) AND c.route_id = $3
					OFFSET $4 LIMIT $5;`

	rows, err := c.db.QueryxContext(ctx, selectQuery, filters.DateFrom, filters.DateTo, filters.EntityID,
		filters.Page*viper.GetInt(config.CompanionsOnPage), viper.GetInt(config.CompanionsOnPage))
	if err != nil {
		return []models.CompanionsRoute{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.QueryErr, Err: err})
	}

	for rows.Next() {
		var companion models.CompanionsRoute
		var userPropertiesRaw []byte
		var routePropertiesRaw []byte

		err := rows.Scan(&companion.DateFrom, &companion.DateTo, &userPropertiesRaw, &companion.RouteName,
			&companion.Price, &routePropertiesRaw, &companion.CityName, &companion.UserID, &companion.RouteID, &companion.ID)
		if err != nil {
			return []models.CompanionsRoute{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ScanErr, Err: err})
		}

		err = json.Unmarshal(userPropertiesRaw, &companion.UserProperties)
		if err != nil {
			return []models.CompanionsRoute{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.BindErr, Err: err})
		}

		err = json.Unmarshal(routePropertiesRaw, &companion.RouteProperties)
		if err != nil {
			return []models.CompanionsRoute{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.BindErr, Err: err})
		}

		companions = append(companions, companion)
	}

	return companions, nil
}

func (c companionsRepo) DeleteCompanionsPlace(ctx context.Context, id int) error {
	tx, err := c.db.Beginx()
	if err != nil {
		return customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.TransactionErr, Err: err})
	}

	deleteQuery := `DELETE FROM companions_places WHERE id = $1`

	res, err := tx.ExecContext(ctx, deleteQuery, id)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return customerr.ErrNormalizer(
				customerr.ErrorPair{Message: customerr.ExecErr, Err: err},
				customerr.ErrorPair{Message: customerr.RollbackErr, Err: rbErr},
			)
		}
		return customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ExecErr, Err: err})
	}
	count, err := res.RowsAffected()
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return customerr.ErrNormalizer(
				customerr.ErrorPair{Message: customerr.RowsErr, Err: err},
				customerr.ErrorPair{Message: customerr.RollbackErr, Err: rbErr},
			)
		}
		return customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.RowsErr, Err: err})
	}
	if count != 1 {
		if rbErr := tx.Rollback(); rbErr != nil {
			return customerr.ErrNormalizer(
				customerr.ErrorPair{Message: customerr.RowsErr, Err: fmt.Errorf(customerr.CountErr, count)},
				customerr.ErrorPair{Message: customerr.RollbackErr, Err: rbErr},
			)
		}
		return customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.RowsErr, Err: fmt.Errorf(customerr.CountErr, count)})
	}

	if err = tx.Commit(); err != nil {
		return customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.CommitErr, Err: err})
	}

	return nil
}

func (c companionsRepo) DeleteCompanionsRoute(ctx context.Context, id int) error {
	tx, err := c.db.Beginx()
	if err != nil {
		return customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.TransactionErr, Err: err})
	}

	deleteQuery := `DELETE FROM companions_routes WHERE id = $1`

	res, err := tx.ExecContext(ctx, deleteQuery, id)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return customerr.ErrNormalizer(
				customerr.ErrorPair{Message: customerr.ExecErr, Err: err},
				customerr.ErrorPair{Message: customerr.RollbackErr, Err: rbErr},
			)
		}
		return customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ExecErr, Err: err})
	}
	count, err := res.RowsAffected()
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return customerr.ErrNormalizer(
				customerr.ErrorPair{Message: customerr.RowsErr, Err: err},
				customerr.ErrorPair{Message: customerr.RollbackErr, Err: rbErr},
			)
		}
		return customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.RowsErr, Err: err})
	}
	if count != 1 {
		if rbErr := tx.Rollback(); rbErr != nil {
			return customerr.ErrNormalizer(
				customerr.ErrorPair{Message: customerr.RowsErr, Err: fmt.Errorf(customerr.CountErr, count)},
				customerr.ErrorPair{Message: customerr.RollbackErr, Err: rbErr},
			)
		}
		return customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.RowsErr, Err: fmt.Errorf(customerr.CountErr, count)})
	}

	if err = tx.Commit(); err != nil {
		return customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.CommitErr, Err: err})
	}

	return nil
}
