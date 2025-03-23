package authJWT

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type claims struct {
	exp time.Time
	iat time.Time

	Payload
}

func (c claims) GetExpirationTime() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(c.exp), nil
}

func (c claims) GetIssuedAt() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(c.iat), nil
}

func (c claims) GetNotBefore() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(c.iat), nil
}

func (c claims) GetIssuer() (string, error) {
	return "giggle-golang", nil
}

func (c claims) GetSubject() (string, error) {
	return c.UserID, nil
}

func (c claims) GetAudience() (jwt.ClaimStrings, error) {
	return jwt.ClaimStrings{"giggle-golang"}, nil
}
