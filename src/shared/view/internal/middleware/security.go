package middleware

import (
	"net/http"
)

func ApplyAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: Implement JWT authentication
		// ref: https://huma.rocks/how-to/oauth2-jwt/#huma-auth-middleware
		next.ServeHTTP(w, r)
	})
}

// // NewJWKSet creates an auto-refreshing key set to validate JWT signatures.
// func NewJWKSet(jwkUrl string) jwk.Set {
// 	jwkCache := jwk.NewCache(context.Background())

// 	// register a minimum refresh interval for this URL.
// 	// when not specified, defaults to Cache-Control and similar resp headers
// 	err := jwkCache.Register(jwkUrl, jwk.WithMinRefreshInterval(10*time.Minute))
// 	if err != nil {
// 		panic("failed to register jwk location")
// 	}

// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	// fetch once on application startup
// 	_, err = jwkCache.Refresh(ctx, jwkUrl)
// 	if err != nil {
// 		panic("failed to fetch on startup")
// 	}
// 	// create the cached key set
// 	return jwk.NewCachedSet(jwkCache, jwkUrl)
// }

// // NewAuthMiddleware creates a middleware that will authorize requests based on
// // the required scopes for the operation.
// func NewAuthMiddleware(api huma.API, jwksURL string) func(ctx huma.Context, next func(huma.Context)) {
// 	keySet := NewJWKSet(jwksURL)

// 	return func(ctx huma.Context, next func(huma.Context)) {
// 		var anyOfNeededScopes []string
// 		isAuthorizationRequired := false
// 		for _, opScheme := range ctx.Operation().Security {
// 			var ok bool
// 			if anyOfNeededScopes, ok = opScheme["myAuth"]; ok {
// 				isAuthorizationRequired = true
// 				break
// 			}
// 		}

// 		if !isAuthorizationRequired {
// 			next(ctx)
// 			return
// 		}

// 		token := strings.TrimPrefix(ctx.Header("Authorization"), "Bearer ")
// 		if len(token) == 0 {
// 			huma.WriteErr(api, ctx, http.StatusUnauthorized, "Unauthorized")
// 			return
// 		}

// 		// Parse and validate the JWT.
// 		parsed, err := jwt.ParseString(token,
// 			jwt.WithKeySet(keySet),
// 			jwt.WithValidate(true),
// 			jwt.WithIssuer("my-issuer"),
// 			jwt.WithAudience("my-audience"),
// 		)
// 		if err != nil {
// 			huma.WriteErr(api, ctx, http.StatusUnauthorized, "Unauthorized")
// 			return
// 		}

// 		// Ensure the claims required for this operation are present.
// 		scopes, _ := parsed.Get("scopes")
// 		if scopes, ok := scopes.([]string); ok {
// 			for _, scope := range scopes {
// 				if slices.Contains(anyOfNeededScopes, scope) {
// 					next(ctx)
// 					return
// 				}
// 			}
// 		}

// 		huma.WriteErr(api, ctx, http.StatusForbidden, "Forbidden")
// 	}
// }
