package service

import (
	"cc/internal/domain"
	"cc/internal/dto"
	"cc/internal/model"
	"cc/internal/storage"
	"cc/pkg/apperror"
	"context"
	"errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	SignIn(ctx context.Context, request dto.Credentials) (domain.User, error)
	SignUp(ctx context.Context, request dto.Credentials) (domain.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (domain.User, error)
	GetUserByName(ctx context.Context, name string) (domain.User, error)
}

type userService struct {
	storage storage.UserStorage
}

func NewUserService(storage storage.UserStorage) UserService {
	return &userService{storage: storage}
}

func (service *userService) SignIn(ctx context.Context, request dto.Credentials) (user domain.User, err error) {
	var usr model.User
	usr, err = service.storage.GetByName(ctx, request.Name)
	if err != nil {
		if apperr, ok := apperror.Is(err, apperror.Internal); ok {
			return user, apperr.WithScope("sign in")
		}

		if errors.Is(err, apperror.NotFound) {
			return user, apperror.BadRequest.WithMessage("invalid name or password")
		}

		return
	}

	err = bcrypt.CompareHashAndPassword(usr.Password, []byte(request.Password))
	if err != nil {
		return user, apperror.BadRequest.WithMessage("invalid name or password")
	}

	return usr.Domain(), nil
}

func (service *userService) SignUp(ctx context.Context, request dto.Credentials) (user domain.User, err error) {
	var exists bool
	exists, err = service.storage.ExistsUserByName(ctx, request.Name)
	if err != nil {
		if apperr, ok := apperror.Is(err, apperror.Internal); ok {
			return user, apperr.WithScope("sign up")
		}

		return
	} else if exists {
		return user, apperror.AlreadyExists.WithMessage("name is already exists")
	}

	var password []byte
	password, err = bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return
	}

	usr := model.User{
		ID:       uuid.New(),
		Name:     request.Name,
		Password: password,
	}
	err = service.storage.CreateUser(ctx, usr)
	if err != nil {
		if apperr, ok := apperror.Is(err, apperror.Internal); ok {
			return user, apperr.WithScope("sign up")
		}

		return
	}

	return usr.Domain(), nil
}

func (service *userService) GetUserByID(ctx context.Context, id uuid.UUID) (user domain.User, err error) {
	var usr model.User
	usr, err = service.storage.GetByID(ctx, id)
	if err != nil {
		if apperr, ok := apperror.Is(err, apperror.Internal); ok {
			return user, apperr.WithScope("get user by id")
		}

		return
	}

	return usr.Domain(), nil
}

func (service *userService) GetUserByName(ctx context.Context, name string) (user domain.User, err error) {
	var usr model.User
	usr, err = service.storage.GetByName(ctx, name)
	if err != nil {
		if apperr, ok := apperror.Is(err, apperror.Internal); ok {
			return user, apperr.WithScope("get user by name")
		}

		return
	}

	return usr.Domain(), nil
}
