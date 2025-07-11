package authJWT

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"giggler-golang/src/shared/env"
	"giggler-golang/src/shared/errutil"
)

var (
	ErrInvalidToken        = errutil.NewCode(errutil.CodeUnauthenticated, errors.New("provided invalid token"))
	ErrExpiredToken        = errutil.NewCode(errutil.CodeUnauthenticated, errors.New("provided expired token"))
	ErrProvidedAccessToken = errutil.NewCode(errutil.CodeUnauthenticated, errors.New("provided access token"))

	jwtSignKey = func() string {
		key := env.Load("JWT_SIGN_KEY")
		if len(key) < 32 {
			panic("sign key for JWT is too short")
		}
		return key
	}()
	jwtSignMethod      = jwt.SigningMethodHS256
	jwtAccessDuration  = env.LoadDuration("JWT_ACCESS_DURATION")
	jwtRefreshDuration = env.LoadDuration("JWT_REFRESH_DURATION")
)

type Payload struct {
	UserID string

	IsRefresh bool
}

func CreateToken(p Payload) (string, error) {
	c := claims{
		iat:     time.Now(),
		Payload: p,
	}

	if p.IsRefresh {
		c.exp = time.Now().Add(jwtRefreshDuration)
	} else {
		c.exp = time.Now().Add(jwtAccessDuration)
	}

	token := jwt.NewWithClaims(jwtSignMethod, c)

	tokenString, err := token.SignedString([]byte(jwtSignKey))
	if err != nil {
		return "", errutil.NewInternal(err)
	}

	return tokenString, nil
}

func VerifyToken(token string) (p Payload, err error) {
	var claims claims
	if _, err := jwt.ParseWithClaims(
		token,
		&claims,
		func(*jwt.Token) (any, error) {
			return []byte(jwtSignKey), nil
		},
	); err != nil {
		return Payload{}, ErrInvalidToken
	}

	if time.Now().After(claims.exp) {
		return Payload{}, ErrExpiredToken
	}

	return claims.Payload, nil
}
