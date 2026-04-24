package jwt

import (
	"giggler-golang/src/core/errutil"
)

var (
	ErrInvalidToken         = errutil.WrapStr("provided invalid token", errutil.ErrorAuthentication)
	ErrExpiredToken         = errutil.WrapStr("provided expired token", errutil.ErrorAuthentication)
	ErrProvidedAccessToken  = errutil.WrapStr("provided access token, expected refresh token", errutil.ErrorAuthentication)
	ErrProvidedRefreshToken = errutil.WrapStr("provided refresh token, expected access token", errutil.ErrorAuthentication)
)
