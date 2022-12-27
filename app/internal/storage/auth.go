package storage

import (
	"cc/app/internal/apperror"
	"cc/app/internal/model"
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type AuthStorage interface {
	CreateSession(ctx context.Context, session model.Session) error
	UpdateSession(ctx context.Context, session model.Session) error
	GetSessionByRefreshToken(ctx context.Context, refreshToken uuid.UUID) (model.Session, error)
}

type authStorage struct {
	db *sqlx.DB
}

func NewAuthStorage(db *sqlx.DB) AuthStorage {
	return &authStorage{db: db}
}

func (storage *authStorage) CreateSession(ctx context.Context, session model.Session) (err error) {
	q := `
INSERT INTO 
    sessions (id, user_id, refresh_token, ip, created_at, updated_at) 
VALUES 
    ($1, $2, $3, $4, $5, $6)
`

	_, err = storage.db.ExecContext(ctx, q,
		session.ID,
		session.UserID,
		session.RefreshToken,
		session.IP,
		session.CreatedAt,
		session.UpdatedAt,
	)
	if err != nil {
		return apperror.ErrInternalError.SetError(err)
	}

	return
}

func (storage *authStorage) UpdateSession(ctx context.Context, session model.Session) (err error) {
	q := `
UPDATE 
    sessions 
SET 
    user_id = $2,
	refresh_token = $3,
	ip = $4,
	created_at = $5,
	updated_at = $6
WHERE 
    id = $1
`

	_, err = storage.db.ExecContext(ctx, q,
		session.ID,
		session.UserID,
		session.RefreshToken,
		session.IP,
		session.CreatedAt,
		session.UpdatedAt,
	)
	if err != nil {
		return apperror.ErrInternalError.SetError(err)
	}

	return
}

func (storage *authStorage) GetSessionByRefreshToken(ctx context.Context, refreshToken uuid.UUID) (session model.Session, err error) {
	q := `
SELECT 
    id, user_id, refresh_token, ip, created_at, updated_at 
FROM 
    sessions
WHERE
	refresh_token = $1
`

	err = storage.db.GetContext(ctx, &session, q, refreshToken)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return session, apperror.ErrNotExists
		}

		return session, apperror.ErrInternalError.SetError(err)
	}

	return
}
