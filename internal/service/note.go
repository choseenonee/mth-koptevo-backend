package service

import (
	"context"
	"mth/internal/models"
	"mth/internal/repository"
	"mth/pkg/log"
)

type noteService struct {
	noteRepo repository.Note
	logger   *log.Logs
}

func InitNoteService(noteRepo repository.Note, logger *log.Logs) Note {
	return noteService{
		noteRepo: noteRepo,
		logger:   logger,
	}
}

func (n noteService) Create(ctx context.Context, noteCreate models.NoteCreate) (int, error) {
	id, err := n.noteRepo.Create(ctx, noteCreate)
	if err != nil {
		n.logger.Error(err.Error())
		return 0, err
	}

	return id, nil
}

func (n noteService) GetByIDs(ctx context.Context, userID int, placeID int) (models.Note, error) {
	note, err := n.noteRepo.GetByIDs(ctx, userID, placeID)
	if err != nil {
		n.logger.Error(err.Error())
		return models.Note{}, err
	}

	return note, nil
}

func (n noteService) GetByID(ctx context.Context, noteID int) (models.Note, error) {
	note, err := n.noteRepo.GetByID(ctx, noteID)
	if err != nil {
		n.logger.Error(err.Error())
		return models.Note{}, err
	}

	return note, nil
}

func (n noteService) GetByUser(ctx context.Context, userID int) ([]models.Note, error) {
	notes, err := n.noteRepo.GetByUser(ctx, userID)
	if err != nil {
		n.logger.Error(err.Error())
		return []models.Note{}, err
	}

	return notes, nil
}

func (n noteService) Update(ctx context.Context, noteUpd models.NoteCreate) error {
	err := n.noteRepo.Update(ctx, noteUpd)
	if err != nil {
		n.logger.Error(err.Error())
		return err
	}

	return nil
}
