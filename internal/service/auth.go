package service

import (
	token "avito/internal/lib/jwt"
	"avito/internal/repository"
	"avito/internal/types"
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/dgrijalva/jwt-go"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
	ErrUserExists         = errors.New("user already exists")
	ErrInvalidAppID       = errors.New("invalid app id")
	ErrInvalidToken       = errors.New("invalid token")
)

type Authorization interface {
	RegisterNewUser(user types.UserCreate) (int64, error)
	LoginUser(user types.UserLoginDTO) (string, error)
	ParseToken(token string) (string, error)
}

type AuthService struct {
	repo repository.Authorization
	log  *slog.Logger
}

func NewAuthService(repo repository.Authorization, log *slog.Logger) *AuthService {
	return &AuthService{repo: repo, log: log}
}

func (a *AuthService) RegisterNewUser(user types.UserCreate) (int64, error) {
	const op = "auth.RegisterNewUser"
	log := a.log.With(
		slog.String("op", op),
		slog.String("username", user.Username),
	)

	log.Info("registering user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash", slog.String("error", err.Error()))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := a.repo.RegisterNewUser(user, passHash)
	if err != nil {
		if errors.Is(err, repository.ErrUserExists) {
			a.log.Info("user allready exists")
			return 0, fmt.Errorf("%s: %w", op, ErrUserExists)
		}
		log.Error("failed to register user", slog.String("error", err.Error()))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user registered")

	return id, nil
}

func (a *AuthService) LoginUser(input types.UserLoginDTO) (string, error) {
	const op = "auth.LoginUser"

	log := a.log.With(
		slog.String("op", op),
		slog.String("username", input.Username),
	)

	log.Info("attempting to login user")

	user, err := a.repo.LoginUser(input.Username)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			a.log.Info("user not found", slog.String("error", err.Error()))
			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}
		a.log.Error("failed to get user", slog.String("error", err.Error()))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		a.log.Info("invalid credentials", slog.String("error", err.Error()))
		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}
	token, err := token.NewToken(user)
	if err != nil {
		a.log.Error("failed to generate token", slog.String("error", err.Error()))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}

func (a *AuthService) ParseToken(token string) (string, error) {
	const op = "auth.ParseToken"

	tokenParsed, err := jwt.ParseWithClaims(token, &types.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			a.log.Error("invalid signing method", slog.String("method", token.Method.Alg()), slog.String("token", token.Raw))
			return nil, fmt.Errorf("%s: invalid signing method", op)
		}

		return []byte(os.Getenv("SECRET")), nil
	})

	if err != nil {
		a.log.Error("failed to parse token", slog.String("error", err.Error()), slog.String("token", token))
		return "", fmt.Errorf("%s: failed to parse token: %w", op, err)
	}

	if !tokenParsed.Valid {
		a.log.Error("invalid token", slog.String("token", token))
		return "", ErrInvalidToken
	}

	claims, ok := tokenParsed.Claims.(*types.TokenClaims)
	if !ok {
		a.log.Warn("token claims type assertion failed", slog.String("token", token))
		return "", ErrInvalidToken
	}

	a.log.Info("token successfully parsed", slog.String("user_id", fmt.Sprintf("%d", claims.User_id)))

	return claims.Username, nil
}
