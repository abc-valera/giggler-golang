package identity

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"giggler-golang/src/features/identity/internal"
	"giggler-golang/src/shared/errutil"
	"giggler-golang/src/shared/must"
)

type AuthTokenService struct {
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
	jwtService           internal.JwtService
}

func initAuthTokenService() AuthTokenService {
	return AuthTokenService{
		accessTokenDuration:  must.GetEnvDuration("GIGGLER_ACCESS_TOKEN_DURATION"),
		refreshTokenDuration: must.GetEnvDuration("GIGGLER_REFRESH_TOKEN_DURATION"),
		jwtService:           internal.InitJwtService(),
	}
}

func (s AuthTokenService) NewAccessTokenJWT(userID uuid.UUID) string {
	return s.jwtService.NewJWT(Token{
		UserID:    userID,
		IsRefresh: false,
		IssuedAt:  *jwt.NewNumericDate(time.Now()),
		ExpiresAt: *jwt.NewNumericDate(time.Now().Add(s.accessTokenDuration)),
	})
}

func (s AuthTokenService) VerifyAccessTokenJWT(token string) (Token, error) {
	claims, err := s.jwtService.VerifyJWT(token)
	if err != nil {
		return Token{}, err
	}

	t, ok := claims.(Token)
	if !ok {
		return Token{}, errors.Join(errutil.ErrValidation, errors.New("invalid auth token"))
	}

	if t.IsRefresh {
		return Token{}, errors.Join(errutil.ErrValidation, errors.New("refresh token used as access token"))
	}

	return t, nil
}

func (s AuthTokenService) NewRefreshTokenJWT(userID uuid.UUID) string {
	return s.jwtService.NewJWT(Token{
		UserID:    userID,
		IsRefresh: true,
		IssuedAt:  *jwt.NewNumericDate(time.Now()),
		ExpiresAt: *jwt.NewNumericDate(time.Now().Add(s.refreshTokenDuration)),
	})
}

func (s AuthTokenService) VerifyRefreshTokenJWT(token string) (Token, error) {
	claims, err := s.jwtService.VerifyJWT(token)
	if err != nil {
		return Token{}, err
	}

	t, ok := claims.(Token)
	if !ok {
		return Token{}, errors.Join(errutil.ErrValidation, errors.New("invalid auth token"))
	}

	if !t.IsRefresh {
		return Token{}, errors.Join(errutil.ErrValidation, errors.New("access token used as refresh token"))
	}

	return t, nil
}

type Token struct {
	UserID    uuid.UUID
	IsRefresh bool
	IssuedAt  jwt.NumericDate
	ExpiresAt jwt.NumericDate
}

const tokenContextKey = "auth_token"

func SetAuthToken(ctx context.Context, token Token) context.Context {
	return context.WithValue(ctx, tokenContextKey, token)
}

func GetAuthToken(ctx context.Context) Token {
	token, ok := ctx.Value(tokenContextKey).(Token)
	if !ok {
		panic("auth token not found in context")
	}
	return token
}

func (p Token) GetExpirationTime() (*jwt.NumericDate, error) {
	return &p.ExpiresAt, nil
}

func (p Token) GetIssuedAt() (*jwt.NumericDate, error) {
	return &p.IssuedAt, nil
}

func (p Token) GetNotBefore() (*jwt.NumericDate, error) {
	return &p.IssuedAt, nil
}

func (p Token) GetIssuer() (string, error) {
	return "giggler-golang", nil
}

func (p Token) GetSubject() (string, error) {
	return p.UserID.String(), nil
}

func (p Token) GetAudience() (jwt.ClaimStrings, error) {
	return jwt.ClaimStrings{"giggler-golang"}, nil
}
