package JWT

import (
	"time"

	"giggler-golang/src/shared/env"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
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
	UserID uuid.UUID
}

type claims struct {
	Payload   Payload
	isRefresh bool
	exp       time.Time
	iat       time.Time
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
	return c.Payload.UserID.String(), nil
}

func (c claims) GetAudience() (jwt.ClaimStrings, error) {
	return jwt.ClaimStrings{"giggle-golang"}, nil
}

// NewAccessToken returns a JWT access token for the payload.
func NewAccessToken(p Payload) string {
	return newToken(p, false)
}

// NewRefreshToken returns a JWT refresh token for the payload.
func NewRefreshToken(p Payload) string {
	return newToken(p, true)
}

func newToken(p Payload, isRefresh bool) string {
	claims := claims{
		Payload:   p,
		isRefresh: isRefresh,
		iat:       time.Now(),
		exp: func() time.Time {
			if isRefresh {
				return time.Now().Add(jwtRefreshDuration)
			}
			return time.Now().Add(jwtAccessDuration)
		}(),
	}

	token := jwt.NewWithClaims(jwtSignMethod, claims)

	tokenString, err := token.SignedString([]byte(jwtSignKey))
	if err != nil {
		panic(err) // should never happen
	}

	return tokenString
}

func VerifyAccessToken(token string) (p Payload, err error) {
	c, err := newPayload(token)
	if err != nil {
		return Payload{}, err
	}
	if c.isRefresh {
		return Payload{}, ErrProvidedRefreshToken
	}
	return c.Payload, nil
}

func VerifyRefreshToken(token string) (p Payload, err error) {
	c, err := newPayload(token)
	if err != nil {
		return Payload{}, err
	}
	if !c.isRefresh {
		return Payload{}, ErrProvidedAccessToken
	}
	return c.Payload, nil
}

func newPayload(token string) (claims, error) {
	var c claims
	if _, err := jwt.ParseWithClaims(
		token,
		&c,
		func(*jwt.Token) (any, error) {
			return []byte(jwtSignKey), nil
		},
	); err != nil {
		return claims{}, ErrInvalidToken
	}

	if time.Now().After(c.exp) {
		return claims{}, ErrExpiredToken
	}

	return c, nil
}
