package repository

import (
	"context"
	"encoding/json"
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"mth/internal/models"
	"mth/pkg/customerr"
)

type placeRepo struct {
	db *sqlx.DB
}

func InitPlaceRepo(db *sqlx.DB) Place {
	return placeRepo{
		db: db,
	}
}

func (p placeRepo) Create(ctx context.Context, placeCreate models.PlaceCreate) (int, error) {
	tx, err := p.db.Beginx()
	if err != nil {
		return 0, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.TransactionErr, Err: err})
	}

	jsonProperties, err := json.Marshal(placeCreate.Properties)
	if err != nil {
		return 0, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.BindErr, Err: err})
	}

	createPlaceQuery := `INSERT INTO places (city_id, district_id, properties) VALUES ($1, $2, $3) RETURNING id;`

	var createdID int
	err = tx.QueryRowxContext(ctx, createPlaceQuery, placeCreate.CityID, placeCreate.DistrictID, jsonProperties).Scan(&createdID)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return 0, customerr.ErrNormalizer(
				customerr.ErrorPair{Message: customerr.ScanErr, Err: err},
				customerr.ErrorPair{Message: customerr.RollbackErr, Err: rbErr},
			)
		}

		return 0, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ScanErr, Err: err})
	}

	createPlaceTagQuery := `INSERT INTO places_tags (place_id, tag_id) VALUES ($1, $2);`
	for _, tagID := range placeCreate.TagIDs {
		_, err := tx.ExecContext(ctx, createPlaceTagQuery, createdID, tagID)
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

// GetAllWithFilter todo: implement tagIDs and pagination
func (p placeRepo) GetAllWithFilter(ctx context.Context, districtID int, cityID int, tagIDs []int, page int) ([]models.Place, error) {
	queryBuilder := squirrel.Select("id", "city_id", "district_id", "properties").
		From("places")

	if districtID != 0 {
		queryBuilder.Where(squirrel.Eq{"district_id": districtID})
	}
	if cityID != 0 {
		queryBuilder.Where(squirrel.Eq{"city_id": cityID})
	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return []models.Place{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.QueryBuild, Err: err})
	}

	rows, err := p.db.QueryContext(ctx, query, args...)
	if err != nil {
		return []models.Place{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ExecErr, Err: err})
	}

	var places []models.Place

	for rows.Next() {
		var place models.Place
		var propertiesRaw []byte

		err = rows.Scan(&place.ID, &place.CityID, &place.DistrictID, &propertiesRaw)
		if err != nil {
			return []models.Place{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ScanErr, Err: err})
		}

		err = json.Unmarshal(propertiesRaw, &place.Properties)
		if err != nil {
			return []models.Place{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.BindErr, Err: err})
		}

		places = append(places, place)
	}

	err = rows.Err()
	if err != nil {
		return []models.Place{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.RowsErr, Err: err})
	}

	return places, nil
}

func (p placeRepo) GetByID(ctx context.Context, placeID int) {
	//TODO implement me
	panic("implement me")
}
