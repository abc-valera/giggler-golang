package jwt

import (
	"errors"

	"giggler-golang/src/core/errutil"
)

var (
	ErrInvalidToken         = errutil.NewCode(errutil.CodeUnauthenticated, errors.New("provided invalid token"))
	ErrExpiredToken         = errutil.NewCode(errutil.CodeUnauthenticated, errors.New("provided expired token"))
	ErrProvidedAccessToken  = errutil.NewCode(errutil.CodeUnauthenticated, errors.New("provided access token, expected refresh token"))
	ErrProvidedRefreshToken = errutil.NewCode(errutil.CodeUnauthenticated, errors.New("provided refresh token, expected access token"))
)
