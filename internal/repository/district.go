package repository

import (
	"context"
	"encoding/json"
	"github.com/jmoiron/sqlx"
	"mth/internal/models"
	"mth/pkg/customerr"
)

type districtRepo struct {
	db *sqlx.DB
}

func InitDistrictRepo(db *sqlx.DB) District {
	return districtRepo{
		db: db,
	}
}

func (d districtRepo) GetByCityID(ctx context.Context, cityID int) ([]models.District, error) {
	query := `SELECT id, name, city_id, properties FROM district WHERE id = $1;`

	rows, err := d.db.QueryContext(ctx, query, cityID)
	if err != nil {
		return []models.District{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ExecErr, Err: err})
	}

	var districts []models.District
	for rows.Next() {
		var district models.District
		var propertiesRaw []byte
		err = rows.Scan(&district.ID, &district.Name, &district.CityID, &propertiesRaw)
		if err != nil {
			return []models.District{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ScanErr, Err: err})
		}

		err = json.Unmarshal(propertiesRaw, &district.Properties)
		if err != nil {
			return []models.District{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.BindErr, Err: err})
		}

		districts = append(districts, district)
	}

	err = rows.Err()
	if err != nil {
		return []models.District{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.RowsErr, Err: err})
	}

	return districts, nil
}
