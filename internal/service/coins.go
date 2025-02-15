package service

import (
	"avito/internal/repository"
	"avito/internal/types"
	"errors"
	"fmt"
	"log/slog"
)

var (
	ErrRecieverNotFound = errors.New("reciever not found")
	ErrNotEnoughtCoins  = errors.New("not enought coins")
	ErrItemNotFound     = errors.New("item not found")
)

type Coins interface {
	GetInfo(username string) (*types.InfoResponse, error)
	Send(username string, details types.SendCoinRequest) error
	BuyItem(username string, item string, req *types.BuyRequest) error
}

type CoinsService struct {
	repo repository.Coins
	log  *slog.Logger
}

func NewCoinsService(repo repository.Coins, log *slog.Logger) *CoinsService {
	return &CoinsService{repo: repo, log: log}
}

func (s *CoinsService) GetInfo(username string) (*types.InfoResponse, error) {
	const op = "coins.get_info"

	res, err := s.repo.GetInfo(username)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			s.log.Warn("user not found", slog.String("user", username))
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return res, nil
}

func (s *CoinsService) Send(username string, details types.SendCoinRequest) error {
	const op = "coins.send"
	err := s.repo.Send(username, details)
	if err != nil {
		if errors.Is(err, repository.ErrNotEnougthFunds) {
			s.log.Warn("not enought coins", slog.String("user", username))
			return fmt.Errorf("%w: %s", ErrNotEnoughtCoins, username)
		} else if errors.Is(err, repository.ErrReceverNotFounded) {
			s.log.Warn("receiver not found", slog.String("user", details.ToUser))
			return fmt.Errorf("%w: %s", ErrRecieverNotFound, op)
		}
		return fmt.Errorf("%w: %s", err, op)
	}

	return err
}

func (s *CoinsService) BuyItem(username string, item string, req *types.BuyRequest) error {
	const op = "coins.buy_item"

	err := s.repo.BuyItem(username, item, req)
	if err != nil {
		if errors.Is(err, repository.ErrNotEnougthFunds) {
			s.log.Warn("not enought coins", slog.String("user", username))
			return fmt.Errorf("%w: %s", ErrNotEnoughtCoins, username)
		} else if errors.Is(err, repository.ErrItemNotFound) {
			s.log.Warn("item not found", slog.String("item", item))
			return fmt.Errorf("%w: %s", ErrItemNotFound, username)
		}
		return fmt.Errorf("%w: %s", err, op)
	}

	return err
}
