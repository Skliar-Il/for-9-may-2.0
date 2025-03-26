package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type mode string

const (
	RefreshToken           mode = "refresh"
	AccessToken            mode = "access"
	RefreshTokenCookieName      = "refresh-token"
	AccessTokenCookieName       = "access-token"
)

func (j *ServiceJWT) GetClaims(id string, tokenMode mode) *jwt.RegisteredClaims {
	var expiration *jwt.NumericDate

	if tokenMode == RefreshToken {
		expiration = jwt.NewNumericDate(time.Now().Add(j.RefreshTimeExp))
	} else if tokenMode == AccessToken {
		expiration = jwt.NewNumericDate(time.Now().Add(j.AccessTimeExp))
	} else {
		panic("invalid type")
	}

	return &jwt.RegisteredClaims{
		Subject:   id,
		ExpiresAt: expiration,
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}
}
