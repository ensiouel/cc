package service

import (
	"cc/app/internal/apperror"
	"cc/app/internal/storage"
	"context"
	"github.com/google/uuid"
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
	if err != nil {
		if apperr, ok := apperror.Is(err, apperror.TypeInternal); ok {
			return tags, apperr.WithScope("TagService.SelectByUser")
		}

		return
	}

	return
}
