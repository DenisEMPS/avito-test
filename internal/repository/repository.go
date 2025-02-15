package repository

import (
	"log/slog"

	"github.com/jmoiron/sqlx"
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
