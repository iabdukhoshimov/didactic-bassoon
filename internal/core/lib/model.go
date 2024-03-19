package lib

import "github.com/golang-jwt/jwt"

type UserClaims struct {
	Id   string `json:"id"`
	Role string `json:"role"`
	jwt.StandardClaims
}
