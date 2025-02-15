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

type Product interface {
}

type Coins interface {
}

type Order interface {
}

type Repository struct {
	Authorization
	Product
	Coins
	Order
}

func NewRepositry(db *sqlx.DB, log *slog.Logger) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db, log),
	}
}
