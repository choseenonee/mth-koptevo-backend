package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/guregu/null/v5"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"mth/internal/models"
	"mth/pkg/config"
	"mth/pkg/customerr"

	_ "github.com/lib/pq"
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

	createPlaceQuery := `INSERT INTO places (city_id, district_id, properties, name, variety) VALUES ($1, $2, $3, $4, $5) RETURNING id;`

	var createdID int
	err = tx.QueryRowxContext(ctx, createPlaceQuery, placeCreate.CityID, placeCreate.DistrictID, jsonProperties,
		placeCreate.Name, placeCreate.Variety).Scan(&createdID)
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

func (p placeRepo) getPlaceTags(ctx context.Context, place *models.Place) error {
	if place == nil {
		return errors.New("you passing nil pointer to the getPlaceTags!")
	}
	query := `SELECT tags.id, name FROM tags
				RIGHT JOIN places_tags pt on tags.id = pt.tag_id
				WHERE pt.place_id = $1;`

	rows, err := p.db.QueryContext(ctx, query, place.ID)
	if err != nil {
		return customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ExecErr, Err: err})
	}

	for rows.Next() {
		var tag models.Tag

		err = rows.Scan(&tag.ID, &tag.Name)
		if err != nil {
			return customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ScanErr, Err: err})
		}

		place.Tags = append(place.Tags, tag)
	}

	err = rows.Err()
	if err != nil {
		return customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.RowsErr, Err: err})
	}

	return nil
}

// GetAllWithFilter todo: implement tagIDs and pagination
func (p placeRepo) GetAllWithFilter(ctx context.Context, districtID int, cityID int, tagIDs []int, page int, name string, variety string) ([]models.Place, error) {
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	queryBuilder := psql.Select("places.id", "city_id", "district_id", "properties", "places.name", "places.variety").
		From("places")

	if len(tagIDs) > 0 {
		queryBuilder = queryBuilder.
			Join("places_tags ON places.id = places_tags.place_id").
			Join("tags ON places_tags.tag_id = tags.id").
			Where(squirrel.Eq{"places_tags.tag_id": tagIDs}).
			GroupBy("places.id").
			Having("COUNT(DISTINCT places_tags.tag_id) >= ?", len(tagIDs))
	}
	if districtID != 0 {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"district_id": districtID})
	}
	if cityID != 0 {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"city_id": cityID})
	}
	if name != "" {
		queryBuilder = queryBuilder.Where(squirrel.Like{"places.name": "%" + name + "%"})
	}
	if variety != "" {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"places.variety": variety})
	}

	// OFFSET с 0 нада бээмс
	queryBuilder = queryBuilder.Limit(uint64(viper.GetInt(config.PlacesOnPage))).Offset(uint64(viper.GetInt(config.PlacesOnPage) * page))

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

		err = rows.Scan(&place.ID, &place.CityID, &place.DistrictID, &propertiesRaw, &place.Name, &place.Variety)
		if err != nil {
			return []models.Place{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ScanErr, Err: err})
		}

		if len(places) > 0 {
			if places[len(places)-1].ID == place.ID {
				continue
			}
		}

		err = json.Unmarshal(propertiesRaw, &place.Properties)
		if err != nil {
			return []models.Place{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.BindErr, Err: err})
		}

		err = p.getPlaceTags(ctx, &place)
		if err != nil {
			return []models.Place{}, fmt.Errorf("get place tags, err: %v", err)
		}

		places = append(places, place)
	}

	err = rows.Err()
	if err != nil {
		return []models.Place{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.RowsErr, Err: err})
	}

	return places, nil
}

func (p placeRepo) GetByID(ctx context.Context, placeID int) (models.Place, error) {
	query := `SELECT places.id, city_id, district_id, properties, places.name, variety, t.id, t.name FROM places
				LEFT JOIN places_tags pt on places.id = pt.place_id
				LEFT JOIN tags t on pt.tag_id = t.id
				WHERE places.id = $1;`

	rows, err := p.db.QueryContext(ctx, query, placeID)
	if err != nil {
		return models.Place{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ExecErr, Err: err})
	}

	var place models.Place
	var propertiesRow []byte
	var tagID null.Int
	var tagName null.String
	for rows.Next() {
		err = rows.Scan(&place.ID, &place.CityID, &place.DistrictID, &propertiesRow, &place.Name, &place.Variety, &tagID, &tagName)
		if err != nil {
			return models.Place{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ScanErr, Err: err})
		}

		err = json.Unmarshal(propertiesRow, &place.Properties)
		if err != nil {
			return models.Place{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.BindErr, Err: err})
		}

		if tagName.Valid {
			var tag models.Tag
			tag.ID = int(tagID.Int64)
			tag.Name = tagName.String
			place.Tags = append(place.Tags, tag)
		}
	}

	err = rows.Err()
	if err != nil {
		return models.Place{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.RowsErr, Err: err})
	}

	return place, nil
}
