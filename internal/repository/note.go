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

	createNoteQuery := `INSERT INTO notes (user_id, place_id, properties) VALUES ($1, $2, $3) RETURNING id;`

	var createdID int
	err = tx.QueryRowxContext(ctx, createNoteQuery, noteCreate.UserID, noteCreate.PlaceID, jsonProperties).Scan(&createdID)
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

func (n noteRepo) GetByID(ctx context.Context, noteID int) (models.Note, error) {
	query := `SELECT n.id, n.user_id, n.place_id, n.properties, 
       			CASE WHEN upc.place_id IS NOT NULL THEN true ELSE false END as joined
			FROM notes n
			LEFT JOIN users_place_checkin upc ON n.place_id = upc.place_id AND n.user_id = upc.user_id
            WHERE id = $1;`

	rows, err := n.db.QueryContext(ctx, query, noteID)
	if err != nil {
		return models.Note{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ExecErr, Err: err})
	}

	var note models.Note
	var propertiesRow []byte
	for rows.Next() {
		err = rows.Scan(&note.ID, &note.UserID, &note.PlaceID, &propertiesRow, &note.IsCheckIn)
		if err != nil {
			return models.Note{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ScanErr, Err: err})
		}

		err = json.Unmarshal(propertiesRow, &note.Properties)
		if err != nil {
			return models.Note{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.BindErr, Err: err})
		}

	}

	err = rows.Err()
	if err != nil {
		return models.Note{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.RowsErr, Err: err})
	}

	return note, nil
}

func (n noteRepo) GetByUser(ctx context.Context, userID int) ([]models.Note, error) {
	query := `SELECT n.id
			FROM notes n
            WHERE n.user_id = $1;`

	rows, err := n.db.QueryContext(ctx, query, userID)
	if err != nil {
		return []models.Note{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ExecErr, Err: err})
	}

	var noteIDs []int
	var noteID int
	for rows.Next() {
		err = rows.Scan(&noteID)
		if err != nil {
			return []models.Note{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ScanErr, Err: err})
		}

		noteIDs = append(noteIDs, noteID)
	}

	var notes []models.Note
	for _, noteID := range noteIDs {
		note, err := n.GetByID(ctx, noteID)
		if err != nil {
			return []models.Note{}, err
		}

		notes = append(notes, note)
	}

	err = rows.Err()
	if err != nil {
		return []models.Note{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.RowsErr, Err: err})
	}

	return notes, nil
}

func (n noteRepo) Update(ctx context.Context, noteUpd models.NoteUpdate) error {
	tx, err := n.db.Beginx()
	if err != nil {
		return customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.TransactionErr, Err: err})
	}

	jsonProperties, err := json.Marshal(noteUpd.Properties)
	if err != nil {
		return customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.BindErr, Err: err})
	}

	query := `UPDATE notes SET properties = $2 WHERE id = $1;`

	_, err = tx.ExecContext(ctx, query, jsonProperties)
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
