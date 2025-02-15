package types

import "github.com/dgrijalva/jwt-go"

type TokenClaims struct {
	jwt.StandardClaims
	User_id int64  `json:"uid"`
	Email   string `json:"email"`
}
