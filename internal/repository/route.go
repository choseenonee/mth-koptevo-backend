package repository

import (
	"context"
	"encoding/json"
	"github.com/guregu/null/v5"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"mth/internal/models"
	"mth/pkg/config"
	"mth/pkg/customerr"
)

type routeRepo struct {
	db *sqlx.DB
}

func InitRouteRepo(db *sqlx.DB) Route {
	return routeRepo{
		db: db,
	}
}

func (r routeRepo) Create(ctx context.Context, route models.RouteCreate) (int, error) {
	propertiesRaw, err := json.Marshal(route.Properties)
	if err != nil {
		return 0, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.BindErr, Err: err})
	}

	tx, err := r.db.Beginx()
	if err != nil {
		return 0, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.TransactionErr, Err: err})
	}

	createRouteQuery := `INSERT INTO routes (city_id, price, name, properties) VALUES ($1, $2, $3, $4) RETURNING id;`

	var createdRouteID int
	err = tx.QueryRowxContext(ctx, createRouteQuery, route.CityID, route.Price, route.Name, propertiesRaw).Scan(&createdRouteID)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return 0, customerr.ErrNormalizer(
				customerr.ErrorPair{Message: customerr.ScanErr, Err: err},
				customerr.ErrorPair{Message: customerr.RollbackErr, Err: rbErr},
			)
		}

		return 0, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ScanErr, Err: err})
	}

	for _, tagID := range route.TagIDs {
		createRouteTagRelationQuery := `INSERT INTO routes_tags (route_id, tag_id) VALUES ($1, $2);`
		_, err = tx.ExecContext(ctx, createRouteTagRelationQuery, createdRouteID, tagID)
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

	for _, placeID := range route.PlaceIDs {
		createRoutePlaceRelationQuery := `INSERT INTO routes_places (route_id, place_id) VALUES ($1, $2);`
		_, err = tx.ExecContext(ctx, createRoutePlaceRelationQuery, createdRouteID, placeID)
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

	return createdRouteID, nil
}

func (r routeRepo) GetByID(ctx context.Context, routeID int) (models.RouteRaw, error) {
	query := `SELECT r.id, r.city_id, r.price, r.name, r.properties, t.id, t.name, rp.place_id  FROM routes r
				LEFT JOIN routes_places rp on r.id = rp.route_id
    			LEFT JOIN routes_tags rt on r.id = rt.route_id
				LEFT JOIN tags t on rt.tag_id = t.id
				WHERE r.id = $1;`

	rows, err := r.db.QueryContext(ctx, query, routeID)
	if err != nil {
		return models.RouteRaw{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ExecErr, Err: err})
	}

	var route models.RouteRaw
	var propertiesRow []byte
	var tagID null.Int
	var tagName null.String
	var placeID null.Int
	for rows.Next() {
		err = rows.Scan(&route.ID, &route.CityID, &route.Price, &route.Name, &propertiesRow, &tagID, &tagName, &placeID)
		if err != nil {
			return models.RouteRaw{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ScanErr, Err: err})
		}

		err = json.Unmarshal(propertiesRow, &route.Properties)
		if err != nil {
			return models.RouteRaw{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.BindErr, Err: err})
		}

		if tagName.Valid {
			var tag models.Tag
			tag.ID = int(tagID.Int64)
			tag.Name = tagName.String
			route.Tags = append(route.Tags, tag)
		}
		if placeID.Valid {
			route.PlaceIDs = append(route.PlaceIDs, int(placeID.Int64))
		}
	}

	err = rows.Err()
	if err != nil {
		return models.RouteRaw{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.RowsErr, Err: err})
	}

	return route, nil
}

func (r routeRepo) GetAll(ctx context.Context, page int) ([]models.RouteRaw, error) {
	query := `SELECT r.id FROM routes r
					ORDER BY r.id
					LIMIT $1 OFFSET $2`

	rows, err := r.db.QueryContext(ctx, query, viper.GetInt(config.PlacesOnPage), page)
	if err != nil {
		return []models.RouteRaw{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ExecErr, Err: err})
	}

	var routeIDs []int
	for rows.Next() {
		var routeID null.Int

		err = rows.Scan(&routeID)
		if err != nil {
			return []models.RouteRaw{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ScanErr, Err: err})
		}

		if routeID.Valid {
			routeIDs = append(routeIDs, int(routeID.Int64))
		}
	}

	var routes []models.RouteRaw
	if len(routeIDs) > 0 {
		for _, routeID := range routeIDs {
			route, err := r.GetByID(ctx, routeID)
			if err != nil {
				return []models.RouteRaw{}, err
			}

			routes = append(routes, route)
		}
	}

	return routes, nil
}
