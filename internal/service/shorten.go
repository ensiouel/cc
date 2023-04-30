package service

import (
	"cc/internal/domain"
	"cc/internal/dto"
	"cc/internal/model"
	"cc/internal/storage"
	"cc/pkg/apperror"
	"cc/pkg/base62"
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
	SelectByUser(ctx context.Context, userID uuid.UUID) (domain.Shortens, error)
	SelectByTags(ctx context.Context, userID uuid.UUID, tags []string) (domain.Shortens, error)
	GetURL(ctx context.Context, shortenID uint64) (string, error)
}

type shortenService struct {
	storage   storage.ShortenStorage
	domainURL string
}

func NewShortenService(storage storage.ShortenStorage, domainURL string) ShortenService {
	return &shortenService{storage: storage, domainURL: domainURL}
}

func (service *shortenService) Create(ctx context.Context, userID uuid.UUID, request dto.CreateShorten) (shorten domain.Shorten, err error) {
	// TODO fix it
	url1, _ := urlx.Parse(service.domainURL)

	url2, err := urlx.Parse(request.URL)
	if err != nil {
		return shorten, apperror.BadRequest.WithMessage("invalid url")
	}

	if url1.Hostname() == url2.Hostname() {
		return shorten, apperror.BadRequest.WithMessage("invalid url")
	}

	var id uint64
	if request.Key != "" {
		id, err = base62.Decode(request.Key)
		if err != nil {
			return
		}

		var exists bool
		exists, err = service.storage.ExistsByID(ctx, userID, id)
		if err != nil {
			if apperr, ok := apperror.Is(err, apperror.Internal); ok {
				return shorten, apperr.WithScope("shortenService.Create")
			}

			return
		}
		if exists {
			return shorten, apperror.AlreadyExists.WithMessage("key already exist")
		}
	} else {
		for exists := true; exists; {
			id = uint64(rand.Uint32())
			exists, err = service.storage.ExistsByID(ctx, userID, id)
			if err != nil {
				if apperr, ok := apperror.Is(err, apperror.Internal); ok {
					return shorten, apperr.WithScope("shortenService.ExistsByID")
				}

				return
			}
		}
	}

	var exists bool
	exists, err = service.storage.ExistsByURL(ctx, userID, request.URL)
	if err != nil {
		if apperr, ok := apperror.Is(err, apperror.Internal); ok {
			return shorten, apperr.WithScope("shortenService.Create")
		}

		return
	}
	if exists {
		return shorten, apperror.AlreadyExists.WithMessage("url already exist")
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
		Tags:      []string{},
		CreatedAt: now,
		UpdatedAt: now,
	}
	err = service.storage.Create(ctx, shrtn)
	if err != nil {
		if apperr, ok := apperror.Is(err, apperror.Internal); ok {
			return shorten, apperr.WithScope("shortenService.Create")
		}

		return
	}

	return shrtn.Domain(service.domainURL), nil
}

func (service *shortenService) Update(ctx context.Context, userID uuid.UUID, shortenID uint64, request dto.UpdateShorten) (shorten domain.Shorten, err error) {
	var shrtn model.Shorten
	shrtn, err = service.storage.GetByID(ctx, shortenID)
	if err != nil {
		if apperr, ok := apperror.Is(err, apperror.Internal); ok {
			return shorten, apperr.WithScope("shortenService.Update.GetByID")
		}

		return
	}

	if shrtn.UserID != userID {
		return shorten, apperror.BadRequest.WithMessage("you don't have access to this shorten")
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
		if apperr, ok := apperror.Is(err, apperror.Internal); ok {
			return shorten, apperr.WithScope("shortenService.Update")
		}

		return
	}

	return shrtn.Domain(service.domainURL), nil
}

func (service *shortenService) Delete(ctx context.Context, userID uuid.UUID, shortenID uint64) (err error) {
	var exists bool
	exists, err = service.storage.ExistsByID(ctx, userID, shortenID)
	if err != nil {
		if apperr, ok := apperror.Is(err, apperror.Internal); ok {
			return apperr.WithScope("shortenService.Delete")
		}

		return
	} else if !exists {
		return apperror.NotFound.WithMessage("shorten with this id does not exist")
	}

	return service.storage.Delete(ctx, userID, shortenID)
}

func (service *shortenService) GetByID(ctx context.Context, id uint64) (shorten domain.Shorten, err error) {
	var shrtn model.Shorten
	shrtn, err = service.storage.GetByID(ctx, id)
	if err != nil {
		if apperr, ok := apperror.Is(err, apperror.Internal); ok {
			return shorten, apperr.WithScope("shortenService.Get")
		}

		return
	}

	return shrtn.Domain(service.domainURL), nil
}

func (service *shortenService) GetURL(ctx context.Context, shortenID uint64) (url string, err error) {
	url, err = service.storage.GetURL(ctx, shortenID)
	if err != nil {
		if apperr, ok := apperror.Is(err, apperror.Internal); ok {
			return url, apperr.WithScope("shortenService.GetURL")
		}

		return
	}

	return
}

func (service *shortenService) SelectByUser(ctx context.Context, userID uuid.UUID) (shortens domain.Shortens, err error) {
	var shrtns model.Shortens
	shrtns, err = service.storage.SelectByUser(ctx, userID)
	if err != nil {
		if apperr, ok := apperror.Is(err, apperror.Internal); ok {
			return shortens, apperr.WithScope("shortenService.SelectByUser")
		}

		return
	}

	return shrtns.Domain(service.domainURL), nil
}

func (service *shortenService) SelectByTags(ctx context.Context, userID uuid.UUID, tags []string) (shortens domain.Shortens, err error) {
	var entities model.Shortens
	entities, err = service.storage.SelectByTags(ctx, userID, tags)
	if err != nil {
		if apperr, ok := apperror.Is(err, apperror.Internal); ok {
			return shortens, apperr.WithScope("shortenService.SelectByTags")
		}

		return
	}

	return entities.Domain(service.domainURL), nil
}
