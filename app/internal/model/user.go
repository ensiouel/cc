package model

import (
	"cc/app/internal/domain"
	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID `db:"id"`
	Name     string    `db:"name"`
	Password []byte    `db:"password"`
}

func (user User) Domain() domain.User {
	return domain.User{
		ID:   user.ID,
		Name: user.Name,
	}
}
