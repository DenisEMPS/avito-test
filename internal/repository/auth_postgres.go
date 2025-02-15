package repository

import (
	"avito/internal/types"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type Authorization interface {
	RegisterNewUser(user types.UserCreate, passHash []byte) (int64, error)
	LoginUser(username string) (types.UserDAO, error)
}

type AuthPostgres struct {
	db  *sqlx.DB
	log *slog.Logger
}

func NewAuthPostgres(db *sqlx.DB, log *slog.Logger) *AuthPostgres {
	return &AuthPostgres{db: db, log: log}
}

func (r *AuthPostgres) RegisterNewUser(user types.UserCreate, passHash []byte) (int64, error) {
	const op = "auth_postgres.RegisterNewUser"

	query := "INSERT INTO users (username, pass_hash, name, surname, birthdate) VALUES($1, $2, $3, $4, $5) RETURNING id"

	var id int64
	err := r.db.QueryRow(query, user.Username, passHash, user.Name, user.Surname, user.Birthdate).Scan(&id)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return 0, ErrUserExists
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (r *AuthPostgres) LoginUser(username string) (types.UserDAO, error) {
	const op = "auth_postgres.LoginUser"

	var user types.UserDAO

	query := "SELECT id, username, pass_hash FROM users WHERE username = $1"

	if err := r.db.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Password); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return types.UserDAO{}, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}
		return types.UserDAO{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}
