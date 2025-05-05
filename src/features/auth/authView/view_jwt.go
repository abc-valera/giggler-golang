package authView

import (
	"time"

	"giggler-golang/src/shared/env"
	"giggler-golang/src/shared/errutil"

	"github.com/golang-jwt/jwt/v5"
)

var (
	errInvalidToken        = errutil.NewCodeMessage(errutil.CodeUnauthenticated, "Provided invalid token")
	errExpiredToken        = errutil.NewCodeMessage(errutil.CodeUnauthenticated, "Provided expired token")
	errProvidedAccessToken = errutil.NewCodeMessage(errutil.CodeUnauthenticated, "Provided access token")
)

var (
	viewJwtSignKey = func() string {
		key := env.Load("JWT_SIGN_KEY")
		if len(key) < 32 {
			panic("sign key for JWT is too short")
		}
		return key
	}()
	viewJwtSignMethod      = jwt.SigningMethodHS256
	viewJwtAccessDuration  = env.LoadDuration("JWT_ACCESS_DURATION")
	viewJwtRefreshDuration = env.LoadDuration("JWT_REFRESH_DURATION")
)

type viewPayload struct {
	UserID    string
	IsRefresh bool
}

func viewCreateToken(payload viewPayload, duration time.Duration) (string, error) {
	token := jwt.New(viewJwtSignMethod)

	claims := jwt.MapClaims{
		"user_id":    payload.UserID,
		"is_refresh": payload.IsRefresh,
		"issued_at":  time.Now().Format(time.RFC3339),
		"expired_at": time.Now().Add(duration).Format(time.RFC3339),
	}
	token.Claims = claims

	tokenString, err := token.SignedString([]byte(viewJwtSignKey))
	if err != nil {
		return "", errutil.NewInternalErr(err)
	}

	return tokenString, nil
}

func viewCreateAccessToken(userID string) (string, error) {
	return viewCreateToken(
		viewPayload{
			UserID:    userID,
			IsRefresh: false,
		},
		viewJwtAccessDuration,
	)
}

func viewCreateRefreshToken(userID string) (string, error) {
	return viewCreateToken(
		viewPayload{
			UserID:    userID,
			IsRefresh: true,
		},
		viewJwtRefreshDuration,
	)
}

func ViewVerifyToken(tokenString string) (viewPayload, error) {
	var claims jwt.MapClaims
	_, err := jwt.ParseWithClaims(
		tokenString,
		&claims,
		func(_ *jwt.Token) (any, error) {
			return []byte(viewJwtSignKey), nil
		},
	)
	if err != nil {
		return viewPayload{}, errInvalidToken
	}

	expiredAt, err := time.Parse(time.RFC3339, claims["expired_at"].(string))
	if err != nil {
		return viewPayload{}, errutil.NewInternalErr(err)
	}

	if time.Now().After(expiredAt) {
		return viewPayload{}, errExpiredToken
	}

	return viewPayload{
		UserID:    claims["user_id"].(string),
		IsRefresh: claims["is_refresh"].(bool),
	}, nil
}
