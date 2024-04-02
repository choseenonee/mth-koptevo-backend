package repository

import (
	"context"
	"encoding/json"
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

func (p placeRepo) GetAllWithFilter(ctx context.Context, districtID int, cityID int, tagIDs []int, page int) ([]models.Place, error) {
	//TODO implement me
	panic("implement me")
}

func (p placeRepo) GetByID(ctx context.Context, placeID int) {
	//TODO implement me
	panic("implement me")
}
