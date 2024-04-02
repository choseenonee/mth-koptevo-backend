package repository

import (
	"context"
	"github.com/jmoiron/sqlx"
	"mth/internal/models"
	"mth/pkg/customerr"
)

type districtRepo struct {
	db *sqlx.DB
}

func InitDistrictRepo(db *sqlx.DB) Districts {
	return districtRepo{db: db}
}

func (d districtRepo) GetByCity(ctx context.Context, cityID int) ([]models.District, error) {
	var districts []models.District

	districtsQuery := `SELECT id, name, properties FROM district WHERE city_id = $1;`

	rows, err := d.db.QueryxContext(ctx, districtsQuery, cityID)
	defer rows.Close()

	if err != nil {
		return []models.District{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.QueryErr, Err: err})
	}

	for rows.Next() {
		var district models.District

		err := rows.Scan(&district.ID, &district.Name, &district.Properties)
		if err != nil {
			return []models.District{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ScanErr, Err: err})
		}
	}

	if err := rows.Err(); err != nil {
		return []models.District{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.RowsErr, Err: err})
	}

	return districts, nil
}
