package service

import (
	"context"
	"mth/internal/models"
)

type Tag interface {
	Create(ctx context.Context, tag models.TagCreate) (int, error)
	GetAll(ctx context.Context) ([]models.Tag, error)
}
