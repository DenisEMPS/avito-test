package jwt

import (
	"avito/internal/types"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func NewToken(user types.UserDAO) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &types.TokenClaims{
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(12 * time.Hour).Unix(),
		},
		Username: user.Username,
		User_id:  user.ID,
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
