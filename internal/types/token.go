package types

import "github.com/dgrijalva/jwt-go"

type TokenClaims struct {
	jwt.StandardClaims
	User_id  int    `json:"user_id"`
	Username string `json:"username"`
}
