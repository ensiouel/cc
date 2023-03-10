package domain

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type Session struct {
	UserID       uuid.UUID `json:"user_id"`
	AccessToken  string    `json:"access_token"`
	RefreshToken uuid.UUID `json:"refresh_token"`
}

type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	jwt.RegisteredClaims
}
