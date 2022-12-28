package service

import (
	"cc/app/internal/apperror"
	"cc/app/internal/domain"
	"cc/app/internal/model"
	"cc/app/internal/storage"
	"context"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"time"
)

type AuthService interface {
	CreateSession(ctx context.Context, userID uuid.UUID, ip string) (domain.Session, error)
	UpdateSession(ctx context.Context, refreshToken uuid.UUID) (domain.Session, error)
	ParseToken(token string) (*jwt.Token, error)
}

type authService struct {
	storage        storage.AuthStorage
	signingKey     string
	expirationTime time.Duration
}

func NewAuthService(storage storage.AuthStorage, signingKey string, expirationTime time.Duration) AuthService {
	return &authService{storage: storage, signingKey: signingKey, expirationTime: expirationTime}
}

func (service *authService) CreateSession(ctx context.Context, userID uuid.UUID, ip string) (session domain.Session, err error) {
	now := time.Now()

	var accessToken string
	accessToken, err = createToken(userID, service.signingKey, now.Add(service.expirationTime))

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
		if apperr, ok := apperror.Internal(err); ok {
			return session, apperr.SetScope("create session")
		}

		return
	}

	session = domain.Session{
		UserID:       userID,
		AccessToken:  accessToken,
		RefreshToken: session.RefreshToken,
	}

	return
}

func (service *authService) UpdateSession(ctx context.Context, refreshToken uuid.UUID) (session domain.Session, err error) {
	var sssn model.Session
	sssn, err = service.storage.GetSessionByRefreshToken(ctx, refreshToken)
	if err != nil {
		if apperr, ok := apperror.Internal(err); ok {
			return session, apperr.SetScope("update session")
		}

		return
	}

	now := time.Now()

	var accessToken string
	accessToken, err = createToken(sssn.UserID, service.signingKey, now.Add(service.expirationTime))
	if err != nil {
		return
	}

	sssn.RefreshToken = uuid.New()
	sssn.UpdatedAt = now

	err = service.storage.UpdateSession(ctx, sssn)
	if err != nil {
		if apperr, ok := apperror.Internal(err); ok {
			return session, apperr.SetScope("update session")
		}

		return
	}

	session = domain.Session{
		UserID:       sssn.UserID,
		AccessToken:  accessToken,
		RefreshToken: session.RefreshToken,
	}

	return
}

func (service *authService) ParseToken(payload string) (token *jwt.Token, err error) {
	token, err = jwt.ParseWithClaims(payload, &domain.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(service.signingKey), nil
	})
	if err != nil {
		if apperr, ok := apperror.Internal(err); ok {
			return token, apperr.SetScope("parse token")
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
		if apperr, ok := apperror.Internal(err); ok {
			return accessToken, apperr.SetScope("create token")
		}

		return
	}

	return
}
