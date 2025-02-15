package service

import (
	"avito/internal/repository"
	"avito/internal/types"
	"log/slog"
)

type Coins interface {
	GetInfo(nickname string) (*types.InfoResponse, error)
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

func (s *CoinsService) GetInfo(nickname string) (*types.InfoResponse, error) {
	return s.repo.GetInfo(nickname)
}

func (s *CoinsService) Send(nickname string, details types.SendCoinRequest) error {
	return s.repo.Send(nickname, details)
}

func (s *CoinsService) BuyItem(username string, item string, req *types.BuyRequest) error {
	return s.repo.BuyItem(username, item, req)
}
