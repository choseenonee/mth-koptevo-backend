package repository

import (
	"context"
	"encoding/json"
	"github.com/jmoiron/sqlx"
	"mth/internal/models"
	"mth/pkg/customerr"
)

type noteRepo struct {
	db *sqlx.DB
}

func InitNoteRepo(db *sqlx.DB) Note {
	return noteRepo{
		db: db,
	}
}

func (n noteRepo) Create(ctx context.Context, noteCreate models.NoteCreate) (int, error) {
	tx, err := n.db.Beginx()
	if err != nil {
		return 0, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.TransactionErr, Err: err})
	}

	jsonProperties, err := json.Marshal(noteCreate.Properties)
	if err != nil {
		return 0, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.BindErr, Err: err})
	}

	createPlaceQuery := `INSERT INTO notes (city_id, district_id, properties, name, variety) VALUES ($1, $2, $3, $4, $5) RETURNING id;`

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

func (n noteRepo) GetByID(ctx context.Context, noteID int) (models.Note, error) {
	//TODO implement me
	panic("implement me")
}

func (n noteRepo) GetByUser(ctx context.Context, userID int) ([]models.Note, error) {
	//TODO implement me
	panic("implement me")
}

func (n noteRepo) Update(ctx context.Context, noteUpd models.NoteUpdate) error {
	//TODO implement me
	panic("implement me")
}
