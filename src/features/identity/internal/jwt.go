package internal

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"

	"giggler-golang/src/shared/errutil"
	"giggler-golang/src/shared/must"
)

type JwtService struct {
	SignKey    string
	SignMethod jwt.SigningMethod
}

func InitJwtService() JwtService {
	return JwtService{
		SignKey: func() string {
			key := must.GetEnv("GIGGLER_JWT_SIGN_KEY")
			if len(key) < 32 {
				panic("sign key for JWT is too short")
			}
			return key
		}(),
		SignMethod: jwt.SigningMethodHS256,
	}
}

func (s JwtService) NewJWT(claims jwt.Claims) string {
	token := jwt.NewWithClaims(s.SignMethod, claims)

	tokenString, err := token.SignedString([]byte(s.SignKey))
	if err != nil {
		panic(err)
	}

	return tokenString
}

func (s JwtService) VerifyJWT(token string) (jwt.Claims, error) {
	var claims jwt.Claims
	if _, err := jwt.ParseWithClaims(
		token,
		claims,
		func(*jwt.Token) (any, error) {
			return []byte(s.SignKey), nil
		},
	); err != nil {
		return claims, errors.Join(errutil.ErrValidation, errors.New("invalid JWT token"), err)
	}

	return claims, nil
}
