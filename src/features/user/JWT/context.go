package JWT

import (
	"context"
)

type contextKeyType string

var contextKey = contextKeyType("jwt")

func Get(ctx context.Context) Payload {
	if v := ctx.Value(contextKey); v != nil {
		if p, ok := v.(Payload); ok {
			return p
		} else {
			panic("JWT: payload in context is not of type Payload")
		}
	}

	panic("JWT: payload not set in context")
}

func Set(ctx context.Context, p Payload) context.Context {
	return context.WithValue(ctx, contextKey, p)
}
