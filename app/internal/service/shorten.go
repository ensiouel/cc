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
	CreateShorten(ctx context.Context, userID uuid.UUID, request dto.CreateShorten) (domain.Shorten, error)
	UpdateShorten(ctx context.Context, userID uuid.UUID, shortenID uint64, request dto.UpdateShorten) (domain.Shorten, error)
	DeleteShorten(ctx context.Context, userID uuid.UUID, shortenID uint64) error
	GetShortenByID(ctx context.Context, id uint64) (domain.Shorten, error)
	GetShortenByURL(ctx context.Context, url string) (domain.Shorten, error)
	GetShortenURL(ctx context.Context, shortenID uint64) (string, error)
	SelectShortensByUserID(ctx context.Context, id uuid.UUID) ([]domain.Shorten, error)
}

type shortenService struct {
	storage storage.ShortenStorage
	host    string
}

func NewShortenService(storage storage.ShortenStorage, host string) ShortenService {
	return &shortenService{storage: storage, host: host}
}

func (service *shortenService) CreateShorten(ctx context.Context, userID uuid.UUID, request dto.CreateShorten) (shorten domain.Shorten, err error) {
	var id uint64
	if request.Key != "" {
		id, err = base62.Decode(request.Key)
		if err != nil {
			return
		}

		var exists bool
		exists, err = service.storage.ExistsShortenByID(ctx, userID, id)
		if err != nil {
			if apperr, ok := apperror.Internal(err); ok {
				return shorten, apperr.SetScope("create shorten")
			}

			return
		} else if exists {
			return shorten, apperror.ErrAlreadyExists.SetMessage("shorten with this id already exist")
		}
	} else {
		for exists := true; exists; {
			id = uint64(rand.Uint32())
			exists, err = service.storage.ExistsShortenByID(ctx, userID, id)
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
	exists, err = service.storage.ExistsShortenByURL(ctx, userID, request.URL)
	if err != nil {
		if apperr, ok := apperror.Internal(err); ok {
			return shorten, apperr.SetScope("create shorten")
		}

		return
	} else if exists {
		return shorten, apperror.ErrAlreadyExists.SetMessage("shorten with this url already exist")
	}

	if request.Title == "" {
		url, _ := urlx.Parse(request.URL)
		request.Title = url.Host
	}

	shrtn := model.Shorten{
		ID:        id,
		UserID:    userID,
		Title:     request.Title,
		URL:       request.URL,
		CreatedAt: time.Now(),
	}
	err = service.storage.CreateShorten(ctx, shrtn)
	if err != nil {
		if apperr, ok := apperror.Internal(err); ok {
			return shorten, apperr.SetScope("create shorten")
		}

		return
	}

	return shrtn.Domain(service.host), nil
}

func (service *shortenService) UpdateShorten(ctx context.Context, userID uuid.UUID, shortenID uint64, request dto.UpdateShorten) (shorten domain.Shorten, err error) {
	var shrtn model.Shorten
	shrtn, err = service.storage.GetShortenByID(ctx, shortenID)
	if err != nil {
		if apperr, ok := apperror.Internal(err); ok {
			return shorten, apperr.SetScope("update shorten")
		}

		return
	}

	if shrtn.UserID != userID {
		//TODO возможно у storage сделать методы Owner (userID, shortenID)
		return shorten, apperror.ErrInvalidCredentials.SetMessage("you don't have access to this shorten")
	}

	shrtn.Title = request.Title
	shrtn.URL = request.URL

	err = service.storage.UpdateShorten(ctx, shrtn)
	if err != nil {
		if apperr, ok := apperror.Internal(err); ok {
			return shorten, apperr.SetScope("update shorten")
		}

		return
	}

	return shrtn.Domain(service.host), nil
}

func (service *shortenService) DeleteShorten(ctx context.Context, userID uuid.UUID, shortenID uint64) (err error) {
	var exists bool
	exists, err = service.storage.ExistsShortenByID(ctx, userID, shortenID)
	if err != nil {
		if apperr, ok := apperror.Internal(err); ok {
			return apperr.SetScope("delete shorten")
		}

		return
	} else if !exists {
		return apperror.ErrNotExists.SetMessage("shorten with this id does not exist")
	}

	return service.storage.DeleteShorten(ctx, userID, shortenID)
}

func (service *shortenService) GetShortenByID(ctx context.Context, id uint64) (shorten domain.Shorten, err error) {
	var shrtn model.Shorten
	shrtn, err = service.storage.GetShortenByID(ctx, id)
	if err != nil {
		if apperr, ok := apperror.Internal(err); ok {
			return shorten, apperr.SetScope("get shorten by id")
		}

		return
	}

	return shrtn.Domain(service.host), nil
}

func (service *shortenService) GetShortenByURL(ctx context.Context, url string) (shorten domain.Shorten, err error) {
	var shrtn model.Shorten
	shrtn, err = service.storage.GetShortenByURL(ctx, url)
	if err != nil {
		if apperr, ok := apperror.Internal(err); ok {
			return shorten, apperr.SetScope("get shorten by url")
		}

		return
	}

	return shrtn.Domain(service.host), nil
}

func (service *shortenService) GetShortenURL(ctx context.Context, shortenID uint64) (url string, err error) {
	url, err = service.storage.GetShortenURL(ctx, shortenID)
	if err != nil {
		if apperr, ok := apperror.Internal(err); ok {
			return url, apperr.SetScope("get Shorten url")
		}

		return
	}

	return
}

func (service *shortenService) SelectShortensByUserID(ctx context.Context, id uuid.UUID) (shortens []domain.Shorten, err error) {
	var shrtns model.Shortens
	shrtns, err = service.storage.SelectShortensByUserID(ctx, id)
	if err != nil {
		if apperr, ok := apperror.Internal(err); ok {
			return shortens, apperr.SetScope("select shortens by user id")
		}

		return
	}

	return shrtns.Domain(service.host), nil
}
