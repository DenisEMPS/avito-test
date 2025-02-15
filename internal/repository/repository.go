package repository

import (
	"errors"
	"log/slog"

	"github.com/jmoiron/sqlx"
)

var (
	ErrUserExists   = errors.New("user allready exists")
	ErrUserNotFound = errors.New("user not found")
)

type Repository struct {
	Authorization
	Coins
}

func NewRepositry(db *sqlx.DB, log *slog.Logger) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db, log),
		Coins:         NewCoinsPostgres(db, log),
	}
}
