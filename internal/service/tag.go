package service

import (
	"context"
	"mth/internal/models"
	"mth/internal/repository"
	"mth/pkg/log"
)

type tagService struct {
	tagRepo repository.Tag
	logger  *log.Logs
}

func InitTagService(tagRepo repository.Tag, logger *log.Logs) Tag {
	return tagService{
		tagRepo: tagRepo,
		logger:  logger,
	}
}

func (t tagService) Create(ctx context.Context, tag models.TagCreate) (int, error) {
	id, err := t.tagRepo.Create(ctx, tag)
	if err != nil {
		t.logger.Error(err.Error())
		return 0, err
	}

	return id, nil
}

func (t tagService) GetAll(ctx context.Context) ([]models.Tag, error) {
	tags, err := t.tagRepo.GetAll(ctx)
	if err != nil {
		t.logger.Error(err.Error())
		return tags, err
	}

	return tags, nil
}
