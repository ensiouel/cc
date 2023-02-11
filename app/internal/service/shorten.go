package service

import (
	"cc/app/internal/apperror"
	"cc/app/internal/domain"
	"cc/app/internal/dto"
	"cc/app/internal/model"
	"cc/app/internal/storage"
	"cc/app/pkg/base62"
	"cc/app/pkg/urlutils"
	"context"
	"github.com/google/uuid"
	"github.com/goware/urlx"
	"math/rand"
	"time"
)

type ShortenService interface {
	Create(ctx context.Context, userID uuid.UUID, request dto.CreateShorten) (domain.Shorten, error)
	Delete(ctx context.Context, userID uuid.UUID, shortenID uint64) error

	Update(ctx context.Context, userID uuid.UUID, shortenID uint64, request dto.UpdateShorten) (domain.Shorten, error)

	GetByID(ctx context.Context, shortenID uint64) (domain.Shorten, error)
	GetByURL(ctx context.Context, url string) (domain.Shorten, error)

	SelectByUser(ctx context.Context, userID uuid.UUID) (domain.Shortens, error)
	SelectByTags(ctx context.Context, userID uuid.UUID, tags []string) (domain.Shortens, error)

	GetURL(ctx context.Context, shortenID uint64) (string, error)
}

type shortenService struct {
	storage storage.ShortenStorage
	host    string
}

func NewShortenService(storage storage.ShortenStorage, host string) ShortenService {
	return &shortenService{storage: storage, host: host}
}

func (service *shortenService) Create(ctx context.Context, userID uuid.UUID, request dto.CreateShorten) (shorten domain.Shorten, err error) {
	var id uint64
	if request.Key != "" {
		id, err = base62.Decode(request.Key)
		if err != nil {
			return
		}

		var exists bool
		exists, err = service.storage.ExistsByID(ctx, userID, id)
		if err != nil {
			if apperr, ok := apperror.Is(err, apperror.TypeInternal); ok {
				return shorten, apperr.WithScope("ShortenService.Create")
			}

			return
		} else if exists {
			return shorten, apperror.AlreadyExists.WithMessage("shorten with this id already exist")
		}
	} else {
		for exists := true; exists; {
			id = uint64(rand.Uint32())
			exists, err = service.storage.ExistsByID(ctx, userID, id)
			if err != nil {
				return
			}
		}
	}

	request.URL, err = urlutils.Normalize(request.URL)
	if err != nil {
		return
	}

	var exists bool
	exists, err = service.storage.ExistsByURL(ctx, userID, request.URL)
	if err != nil {
		if apperr, ok := apperror.Is(err, apperror.TypeInternal); ok {
			return shorten, apperr.WithScope("ShortenService.Create")
		}

		return
	} else if exists {
		return shorten, apperror.AlreadyExists.WithMessage("shorten with this url already exist")
	}

	if request.Title == "" {
		url, _ := urlx.Parse(request.URL)
		request.Title = url.Host
	}

	now := time.Now()

	shrtn := model.Shorten{
		ID:        id,
		UserID:    userID,
		Title:     request.Title,
		URL:       request.URL,
		CreatedAt: now,
		UpdatedAt: now,
	}
	err = service.storage.Create(ctx, shrtn)
	if err != nil {
		if apperr, ok := apperror.Is(err, apperror.TypeInternal); ok {
			return shorten, apperr.WithScope("ShortenService.Create")
		}

		return
	}

	return shrtn.Domain(service.host), nil
}

func (service *shortenService) Update(ctx context.Context, userID uuid.UUID, shortenID uint64, request dto.UpdateShorten) (shorten domain.Shorten, err error) {
	var shrtn model.Shorten
	shrtn, err = service.storage.GetByID(ctx, shortenID)
	if err != nil {
		if apperr, ok := apperror.Is(err, apperror.TypeInternal); ok {
			return shorten, apperr.WithScope("ShortenService.Update.GetByID")
		}

		return
	}

	if shrtn.UserID != userID {
		return shorten, apperror.InvalidCredentials.WithMessage("you don't have access to this shorten")
	}

	if request.Title != "" {
		shrtn.Title = request.Title
	}

	if request.URL != "" {
		shrtn.URL = request.URL
	}

	if len(request.Tags) != 0 {
		shrtn.Tags = request.Tags
	}

	shrtn.UpdatedAt = time.Now()

	err = service.storage.Update(ctx, shrtn)
	if err != nil {
		if apperr, ok := apperror.Is(err, apperror.TypeInternal); ok {
			return shorten, apperr.WithScope("ShortenService.Update")
		}

		return
	}

	return shrtn.Domain(service.host), nil
}

func (service *shortenService) Delete(ctx context.Context, userID uuid.UUID, shortenID uint64) (err error) {
	var exists bool
	exists, err = service.storage.ExistsByID(ctx, userID, shortenID)
	if err != nil {
		if apperr, ok := apperror.Is(err, apperror.TypeInternal); ok {
			return apperr.WithScope("ShortenService.Delete")
		}

		return
	} else if !exists {
		return apperror.NotExists.WithMessage("shorten with this id does not exist")
	}

	return service.storage.Delete(ctx, userID, shortenID)
}

func (service *shortenService) GetByID(ctx context.Context, id uint64) (shorten domain.Shorten, err error) {
	var shrtn model.Shorten
	shrtn, err = service.storage.GetByID(ctx, id)
	if err != nil {
		if apperr, ok := apperror.Is(err, apperror.TypeInternal); ok {
			return shorten, apperr.WithScope("ShortenService.Get")
		}

		return
	}

	return shrtn.Domain(service.host), nil
}

func (service *shortenService) GetByURL(ctx context.Context, url string) (shorten domain.Shorten, err error) {
	var shrtn model.Shorten
	shrtn, err = service.storage.GetByURL(ctx, url)
	if err != nil {
		if apperr, ok := apperror.Is(err, apperror.TypeInternal); ok {
			return shorten, apperr.WithScope("ShortenService.GetByURL")
		}

		return
	}

	return shrtn.Domain(service.host), nil
}

func (service *shortenService) GetURL(ctx context.Context, shortenID uint64) (url string, err error) {
	url, err = service.storage.GetURL(ctx, shortenID)
	if err != nil {
		if apperr, ok := apperror.Is(err, apperror.TypeInternal); ok {
			return url, apperr.WithScope("ShortenService.GetURL")
		}

		return
	}

	return
}

func (service *shortenService) SelectByUser(ctx context.Context, userID uuid.UUID) (shortens domain.Shortens, err error) {
	var shrtns model.Shortens
	shrtns, err = service.storage.SelectByUser(ctx, userID)
	if err != nil {
		if apperr, ok := apperror.Is(err, apperror.TypeInternal); ok {
			return shortens, apperr.WithScope("ShortenService.SelectByUser")
		}

		return
	}

	return shrtns.Domain(service.host), nil
}

func (service *shortenService) SelectByTags(ctx context.Context, userID uuid.UUID, tags []string) (shortens domain.Shortens, err error) {
	var entities model.Shortens
	entities, err = service.storage.SelectByTags(ctx, userID, tags)
	if err != nil {
		if apperr, ok := apperror.Is(err, apperror.TypeInternal); ok {
			return shortens, apperr.WithScope("ShortenService.SelectByTags")
		}

		return
	}

	return entities.Domain(service.host), nil
}
