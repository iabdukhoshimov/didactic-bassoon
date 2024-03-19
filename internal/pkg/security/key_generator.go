package security

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

type OrganizationClaims struct {
	Id string `json:"org_id"`
	jwt.StandardClaims
}

type OrganizationKeyGenerator struct {
	TokenLifeTime time.Duration
	SecretKey     []byte
}

func NewOrganizationKeyGenerator(secretKey string, tokenLifeTime string) *OrganizationKeyGenerator {
	tokenLifeTimeDuration, err := time.ParseDuration(tokenLifeTime)
	if err != nil {
		panic(err)
	}
	return &OrganizationKeyGenerator{
		TokenLifeTime: tokenLifeTimeDuration,
		SecretKey:     []byte(secretKey),
	}
}

func (g *OrganizationKeyGenerator) GenerateOrganizationKey(claims OrganizationClaims) (string, error) {
	claims.StandardClaims = jwt.StandardClaims{
		ExpiresAt: jwt.TimeFunc().Add(g.TokenLifeTime).Unix(),
		IssuedAt:  jwt.TimeFunc().Unix(),
	}
	orgKey := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return orgKey.SignedString(g.SecretKey)
}

func (g *OrganizationKeyGenerator) ParseOrganizationKey(tokenString string) (*OrganizationClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &OrganizationClaims{}, func(token *jwt.Token) (interface{}, error) {
		return g.SecretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(*OrganizationClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
