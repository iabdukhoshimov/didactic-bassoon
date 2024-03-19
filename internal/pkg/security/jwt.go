package security

import (
	"fmt"
	"time"

	"gitlab.com/tour/internal/config"

	"github.com/golang-jwt/jwt"
	models "gitlab.com/tour/internal/core/lib"
)

type IssuerInterface interface {
	NewAccessToken(models.UserClaims) (string, error)
	NewRefreshToken(models.UserClaims) (string, error)
	ParseAccessToken(string) (*models.UserClaims, error)
	ParseRefreshToken(string) (*models.UserClaims, error)
}

type Issuer struct {
	AccessLifeTime  int
	RefreshLifeTime int
	TokenSecret     []byte
}

func NewIssuer(cfg *config.Config) *Issuer {
	return &Issuer{
		TokenSecret:     []byte(cfg.JWT.SecretToken),
		AccessLifeTime:  cfg.JWT.AccessLifeTime,
		RefreshLifeTime: cfg.JWT.RefreshLifeTime,
	}
}

func (i *Issuer) NewAccessToken(claims models.UserClaims) (string, error) {
	claims.StandardClaims = jwt.StandardClaims{
		ExpiresAt: jwt.TimeFunc().Add(time.Duration(i.AccessLifeTime) * time.Second).Unix(),
		IssuedAt:  jwt.TimeFunc().Unix(),
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return accessToken.SignedString(i.TokenSecret)
}

func (i *Issuer) NewRefreshToken(claims models.UserClaims) (string, error) {
	claims.StandardClaims = jwt.StandardClaims{
		ExpiresAt: jwt.TimeFunc().Add(time.Duration(i.RefreshLifeTime) * time.Second).Unix(),
		IssuedAt:  jwt.TimeFunc().Unix(),
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return refreshToken.SignedString(i.TokenSecret)
}

func (i *Issuer) ParseAccessToken(accessToken string) (*models.UserClaims, error) {
	parsedAccessToken, err := jwt.ParseWithClaims(accessToken, &models.UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return i.TokenSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if !parsedAccessToken.Valid {
		return nil, fmt.Errorf("invalid access token")
	}

	return parsedAccessToken.Claims.(*models.UserClaims), nil
}

func (i *Issuer) ParseRefreshToken(refreshToken string) (*models.UserClaims, error) {
	parsedRefreshToken, err := jwt.ParseWithClaims(refreshToken, &models.UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return i.TokenSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if !parsedRefreshToken.Valid {
		return nil, fmt.Errorf("invalid refresh token")
	}

	return parsedRefreshToken.Claims.(*models.UserClaims), nil
}
