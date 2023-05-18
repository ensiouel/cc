package service

import (
	"cc/internal/config"
	"cc/internal/domain"
	"cc/internal/dto"
	"cc/internal/model"
	"cc/internal/storage"
	"cc/pkg/apperror"
	"context"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"time"
)

type AuthService interface {
	CreateSession(ctx context.Context, userID uuid.UUID, ip string) (domain.Session, error)
	UpdateSession(ctx context.Context, request dto.Refresh) (domain.Session, error)
	ParseToken(token string) (*jwt.Token, error)
}

type authService struct {
	storage storage.AuthStorage
	config  config.Auth
}

func NewAuthService(storage storage.AuthStorage, config config.Auth) AuthService {
	return &authService{storage: storage, config: config}
}

func (service *authService) CreateSession(ctx context.Context, userID uuid.UUID, ip string) (session domain.Session, err error) {
	now := time.Now()

	var accessToken string
	accessToken, err = createToken(userID, service.config.SigningKey, now.Add(service.config.ExpirationAt))

	sssn := model.Session{
		ID:           uuid.New(),
		UserID:       userID,
		RefreshToken: uuid.New(),
		IP:           ip,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	err = service.storage.CreateSession(ctx, sssn)
	if err != nil {
		if apperr, ok := apperror.Is(err, apperror.Internal); ok {
			return session, apperr.WithScope("create session")
		}

		return
	}

	session = domain.Session{
		UserID:       sssn.UserID,
		AccessToken:  accessToken,
		RefreshToken: sssn.RefreshToken,
	}

	return
}

func (service *authService) UpdateSession(ctx context.Context, request dto.Refresh) (session domain.Session, err error) {
	var sssn model.Session
	sssn, err = service.storage.GetSessionByRefreshToken(ctx, request.RefreshToken)
	if err != nil {
		if apperr, ok := apperror.Is(err, apperror.Internal); ok {
			return session, apperr.WithScope("update session")
		}

		return
	}

	now := time.Now()

	var accessToken string
	accessToken, err = createToken(sssn.UserID, service.config.SigningKey, now.Add(service.config.ExpirationAt))
	if err != nil {
		return
	}

	sssn.RefreshToken = uuid.New()
	sssn.UpdatedAt = now

	err = service.storage.UpdateSession(ctx, sssn)
	if err != nil {
		if apperr, ok := apperror.Is(err, apperror.Internal); ok {
			return session, apperr.WithScope("update session")
		}

		return
	}

	session = domain.Session{
		UserID:       sssn.UserID,
		AccessToken:  accessToken,
		RefreshToken: sssn.RefreshToken,
	}

	return
}

func (service *authService) ParseToken(payload string) (token *jwt.Token, err error) {
	token, err = jwt.ParseWithClaims(payload, &domain.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(service.config.SigningKey), nil
	})
	if err != nil {
		if apperr, ok := apperror.Is(err, apperror.Internal); ok {
			return token, apperr.WithScope("parse token")
		}

		return
	}

	return
}

func createToken(userID uuid.UUID, signingKey string, expirationTime time.Time) (accessToken string, err error) {
	claims := domain.Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	accessToken, err = token.SignedString([]byte(signingKey))
	if err != nil {
		if apperr, ok := apperror.Is(err, apperror.Internal); ok {
			return accessToken, apperr.WithScope("create token")
		}

		return
	}

	return
}
