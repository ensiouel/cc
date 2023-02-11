package storage

import (
	"cc/app/internal/apperror"
	"cc/app/internal/model"
	"cc/app/pkg/postgres"
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type AuthStorage interface {
	CreateSession(ctx context.Context, session model.Session) error
	UpdateSession(ctx context.Context, session model.Session) error
	GetSessionByRefreshToken(ctx context.Context, refreshToken uuid.UUID) (model.Session, error)
}

type authStorage struct {
	client postgres.Client
}

func NewAuthStorage(client postgres.Client) AuthStorage {
	return &authStorage{client: client}
}

func (storage *authStorage) CreateSession(ctx context.Context, session model.Session) (err error) {
	q := `
INSERT INTO 
    sessions (id, user_id, refresh_token, ip, created_at, updated_at) 
VALUES 
    ($1, $2, $3, $4, $5, $6)
`

	_, err = storage.client.Exec(ctx, q,
		session.ID,
		session.UserID,
		session.RefreshToken,
		session.IP,
		session.CreatedAt,
		session.UpdatedAt,
	)
	if err != nil {
		return apperror.Internal.WithError(err)
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

	_, err = storage.client.Exec(ctx, q,
		session.ID,
		session.UserID,
		session.RefreshToken,
		session.IP,
		session.CreatedAt,
		session.UpdatedAt,
	)
	if err != nil {
		return apperror.Internal.WithError(err)
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

	err = storage.client.Get(ctx, &session, q, refreshToken)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return session, apperror.NotExists
		}

		return session, apperror.Internal.WithError(err)
	}

	return
}
