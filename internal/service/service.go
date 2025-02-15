package service

import (
	"avito/internal/repository"
	"log/slog"
)

type Product interface {
}

type Coins interface {
}

type Order interface {
}

type Service struct {
	Authorization
	Product
	Coins
	Order
}

func NewService(repos *repository.Repository, log *slog.Logger) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization, log),
	}
}
