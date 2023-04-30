package model

import (
	"github.com/google/uuid"
	"time"
)

type Session struct {
	ID           uuid.UUID `db:"id"`
	UserID       uuid.UUID `db:"user_id"`
	RefreshToken uuid.UUID `db:"refresh_token"`
	IP           string    `db:"ip"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}
