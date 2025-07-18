// Code generated by ogen, DO NOT EDIT.

package viewgen

import (
	"context"

	ht "github.com/ogen-go/ogen/http"
)

// UnimplementedHandler is no-op Handler which returns http.ErrNotImplemented.
type UnimplementedHandler struct{}

var _ Handler = UnimplementedHandler{}

// AuthLoginPost implements AuthLoginPost operation.
//
// Performs user authentication.
//
// POST /auth/login
func (UnimplementedHandler) AuthLoginPost(ctx context.Context, req *AuthLoginPostReq) (r *AuthLoginPostOK, _ error) {
	return r, ht.ErrNotImplemented
}

// AuthRefreshPost implements AuthRefreshPost operation.
//
// Exchanges a refresh token for an access token.
//
// POST /auth/refresh
func (UnimplementedHandler) AuthRefreshPost(ctx context.Context, req *AuthRefreshPostReq) (r *AuthRefreshPostOK, _ error) {
	return r, ht.ErrNotImplemented
}

// AuthRegisterPost implements AuthRegisterPost operation.
//
// Performs user registration.
//
// POST /auth/register
func (UnimplementedHandler) AuthRegisterPost(ctx context.Context, req *AuthRegisterPostReq) error {
	return ht.ErrNotImplemented
}

// JokesDel implements JokesDel operation.
//
// Deletes joke for current user.
//
// DELETE /jokes
func (UnimplementedHandler) JokesDel(ctx context.Context, req *JokesDelReq) error {
	return ht.ErrNotImplemented
}

// JokesGet implements JokesGet operation.
//
// Returns the most relevant jokes.
//
// GET /jokes
func (UnimplementedHandler) JokesGet(ctx context.Context, params JokesGetParams) (r JokesSchema, _ error) {
	return r, ht.ErrNotImplemented
}

// JokesPost implements JokesPost operation.
//
// Creates a new joke for current user.
//
// POST /jokes
func (UnimplementedHandler) JokesPost(ctx context.Context, req *JokesPostReq) (r *JokeSchema, _ error) {
	return r, ht.ErrNotImplemented
}

// JokesPut implements JokesPut operation.
//
// Updates joke for current user.
//
// PUT /jokes
func (UnimplementedHandler) JokesPut(ctx context.Context, req *JokesPutReq) (r *JokeSchema, _ error) {
	return r, ht.ErrNotImplemented
}

// UserDel implements UserDel operation.
//
// Deletes current user profile.
//
// DELETE /user
func (UnimplementedHandler) UserDel(ctx context.Context, req *UserDelReq) error {
	return ht.ErrNotImplemented
}

// UserGet implements UserGet operation.
//
// Returns current user profile.
//
// GET /user
func (UnimplementedHandler) UserGet(ctx context.Context) (r *UserSchema, _ error) {
	return r, ht.ErrNotImplemented
}

// UserPut implements UserPut operation.
//
// Updates current user profile.
//
// PUT /user
func (UnimplementedHandler) UserPut(ctx context.Context, req *UserPutReq) (r *UserSchema, _ error) {
	return r, ht.ErrNotImplemented
}

// NewError creates *ErrorSchemaStatusCode from error returned by handler.
//
// Used for common default response.
func (UnimplementedHandler) NewError(ctx context.Context, err error) (r *ErrorSchemaStatusCode) {
	r = new(ErrorSchemaStatusCode)
	return r
}
