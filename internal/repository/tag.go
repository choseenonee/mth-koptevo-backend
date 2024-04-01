package repository

import (
	"context"
	"github.com/jmoiron/sqlx"
	"mth/internal/models"
	"mth/pkg/customerr"
)

type tagRepo struct {
	db *sqlx.DB
}

func InitTagRepo(db *sqlx.DB) Tag {
	return tagRepo{
		db: db,
	}
}

func (t tagRepo) Create(ctx context.Context, tag models.TagCreate) (int, error) {
	tx, err := t.db.Beginx()
	if err != nil {
		return 0, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.TransactionErr, Err: err})
	}

	createTagQuery := `INSERT INTO tags (name) VALUES ($1) RETURNING id;`

	var createdTagID int
	err = tx.QueryRowxContext(ctx, createTagQuery, tag.Name).Scan(&createdTagID)
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

	return createdTagID, nil
}

func (t tagRepo) GetAll(ctx context.Context) ([]models.Tag, error) {
	getAllTagsQuery := `SELECT id, name FROM tags;`

	rows, err := t.db.QueryContext(ctx, getAllTagsQuery)
	if err != nil {
		return []models.Tag{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ExecErr, Err: err})
	}

	var tags []models.Tag

	for rows.Next() {
		var tag models.Tag
		err := rows.Scan(&tag.ID, &tag.Name)
		if err != nil {
			return []models.Tag{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.ScanErr, Err: nil})
		}
		tags = append(tags, tag)
	}

	err = rows.Err()
	if err != nil {
		return []models.Tag{}, customerr.ErrNormalizer(customerr.ErrorPair{Message: customerr.RowsErr, Err: nil})
	}
	
	return tags, nil
}
