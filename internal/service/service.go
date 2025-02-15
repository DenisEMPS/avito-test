package service

import (
	"avito/internal/repository"
	"log/slog"
)

type Service struct {
	Authorization
	Coins
}

func NewService(repos *repository.Repository, log *slog.Logger) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization, log),
		Coins:         NewCoinsService(repos.Coins, log),
	}
}
