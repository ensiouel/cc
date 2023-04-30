package service

import (
	"cc/internal/storage"
	"cc/pkg/apperror"
	"context"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type TagService interface {
	SelectByUser(ctx context.Context, userID uuid.UUID) ([]string, error)
}

type tagService struct {
	storage storage.TagStorage
}

func NewTagService(storage storage.TagStorage) TagService {
	return &tagService{storage: storage}
}

func (service *tagService) SelectByUser(ctx context.Context, userID uuid.UUID) (tags []string, err error) {
	tags, err = service.storage.SelectByUser(ctx, userID)
	if err != nil && !errors.Is(err, apperror.NotFound) {
		if apperr, ok := apperror.Is(err, apperror.Internal); ok {
			return tags, apperr.WithScope("TagService.SelectByUser")
		}

		return
	}

	return
}
