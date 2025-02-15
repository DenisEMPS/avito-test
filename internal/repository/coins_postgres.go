package repository

import (
	"avito/internal/types"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jmoiron/sqlx"
)

var (
	ErrNotEnougthFunds    = errors.New("not enought funds")
	ErrRecieverNotFounded = errors.New("reciever not founded")
	ErrItemNotFounded     = errors.New("item not founded")
)

type Coins interface {
	GetInfo(nickname string) (*types.InfoResponse, error)
	Send(username string, details types.SendCoinRequest) error
	ByItem(username string, item string, quantity float64) error
}

type CoinsPostgres struct {
	db  *sqlx.DB
	log *slog.Logger
}

func NewCoinsPostgres(db *sqlx.DB, log *slog.Logger) *CoinsPostgres {
	return &CoinsPostgres{db: db, log: log}
}

func (r *CoinsPostgres) GetInfo(nickname string) (*types.InfoResponse, error) {
	const op = "coins_postgres.GetInfo"

	var res types.InfoResponse
	var receivedCoins []types.Received
	var sentCoins []types.Sent

	queryReceived := "SELECT from_user, amount FROM coins_transactions WHERE to_user = $1"
	rowsReceived, err := r.db.Query(queryReceived, nickname)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rowsReceived.Close()

	for rowsReceived.Next() {
		var received types.Received
		if err := rowsReceived.Scan(&received.FromUser, &received.Amount); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		receivedCoins = append(receivedCoins, received)
	}

	querySent := "SELECT to_user, amount FROM coins_transactions WHERE from_user = $1"
	rowsSent, err := r.db.Query(querySent, nickname)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rowsSent.Close()

	for rowsSent.Next() {
		var sent types.Sent
		if err := rowsSent.Scan(&sent.ToUser, &sent.Amount); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		sentCoins = append(sentCoins, sent)
	}

	res.CoinsHistory = types.CoinsHistory{
		Received: receivedCoins,
		Sent:     sentCoins,
	}

	queryCoins := "SELECT coins FROM users WHERE username = $1"
	err = r.db.QueryRow(queryCoins, nickname).Scan(&res.Coins)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	queryInventory := "SELECT item_type, quantity FROM inventory WHERE username = $1"
	rowsInventory, err := r.db.Query(queryInventory, nickname)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rowsInventory.Close()

	var inventory []types.Inventory
	for rowsInventory.Next() {
		var item types.Inventory
		if err := rowsInventory.Scan(&item.Type, &item.Quantity); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		inventory = append(inventory, item)
	}

	res.Inventory = inventory

	return &res, nil
}

func (r *CoinsPostgres) Send(username string, details types.SendCoinRequest) error {
	const op = "coins_postgres.Send"

	var coins int
	query := "SELECT coins from users WHERE username = $1"
	err := r.db.QueryRow(query, username).Scan(&coins)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrUserNotFound
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	if coins < details.Amount {
		return ErrNotEnougthFunds
	}

	var recExists bool
	query = "SELECT 1 FROM users WHERE username=$1 LIMIT 1"
	err = r.db.QueryRow(query, details.ToUser).Scan(&recExists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrRecieverNotFounded
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("error to start transaction: %w %s", err, op)
	}
	queryCoinsUpdate := `UPDATE users SET coins = coins - $1
						WHERE username = $2`
	_, err = tx.Exec(queryCoinsUpdate, details.Amount, username)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update coins balance: %w %s", err, op)
	}

	queryCoinsUpdate = `UPDATE users SET coins = coins + $1 
						WHERE username = $2`
	_, err = tx.Exec(queryCoinsUpdate, details.Amount, details.ToUser)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update coins balance: %w %s", err, op)
	}

	return tx.Commit()
}

func (r *CoinsPostgres) ByItem(username string, item string, quantity float64) error {
	const op = "coins_postgres.ByItem"

	var coins float64
	query := "SELECT coins FROM users WHERE username = $1"

	err := r.db.QueryRow(query, username).Scan(&coins)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrUserNotFound
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	var price float64
	query = "SELECT price FROM items WHERE name = $1"
	err = r.db.QueryRow(query, item).Scan(price)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrItemNotFounded
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	if price*quantity > coins {
		return ErrNotEnougthFunds
	}

	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("error to start transaction: %w %s", err, op)
	}

	query = `UPDATE users SET coins = coins - $1
			WHERE username = $2`
	_, err = tx.Exec(query, username, price)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update coins balance: %w %s %s", err, username, op)
	}

	query = `INSERT INTO inventory (username, quantity, item_id)
			VALUES($1, $2, (SELECT id FROM items WHERE name = $3))`

	_, err = tx.Exec(query, username, quantity, item)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to insert data: %w %s %s", err, username, item)
	}

	return tx.Commit()
}
