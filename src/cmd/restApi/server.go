package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/pprof"

	"giggler-golang/src/features/auth/authJWT"
	"giggler-golang/src/features/auth/authView"
	"giggler-golang/src/features/joke/jokeView"
	"giggler-golang/src/features/user/userView"
	"giggler-golang/src/shared/app"
	"giggler-golang/src/shared/contexts"
	"giggler-golang/src/shared/env"
	"giggler-golang/src/shared/errutil"
	"giggler-golang/src/shared/errutil/errutilView"
	"giggler-golang/src/shared/otel/otelView"
	"giggler-golang/src/shared/view/viewgen"
)

func initServer() http.Server {
	mux := http.NewServeMux()

	if env.LoadBool("IS_HTTP_PPROF_INTERFACE_ENABLED") {
		mux.HandleFunc("GET /debug/pprof/", pprof.Index)
		mux.HandleFunc("GET /debug/pprof/cmdline", pprof.Cmdline)
		mux.HandleFunc("GET /debug/pprof/profile", pprof.Profile)
		mux.HandleFunc("GET /debug/pprof/symbol", pprof.Symbol)
		mux.HandleFunc("GET /debug/pprof/trace", pprof.Trace)
	}

	mux.HandleFunc("/build-version", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(app.BuildVersion()))
	})

	// TODO: check middlewares order
	var handler http.Handler = errutil.Must(viewgen.NewServer(handler, securityHandler{}))
	handler = applyCorsMiddleware(handler)
	handler = ApplyLogMiddleware(handler)
	handler = otelView.ApplyTracerMiddleware(handler)
	handler = applyRecovererMiddleware(handler)

	mux.Handle("/", handler)

	return http.Server{
		Addr:    fmt.Sprintf(":%s", env.Load("API_PORT")),
		Handler: mux,
	}
}

type (
	authHandler = authView.Handler
	userHandler = userView.Handler
	jokeHandler = jokeView.Handler

	errorHandler = errutilView.Handler
)

var handler viewgen.Handler = struct {
	authHandler
	userHandler
	jokeHandler

	errorHandler
}{
	authHandler: authView.Handler{},
	userHandler: userView.Handler{},
	jokeHandler: jokeView.Handler{},

	errorHandler: errutilView.Handler{},
}

type securityHandler struct{}

func (securityHandler) HandleBearerAuth(ctx context.Context, _ string, t viewgen.BearerAuth) (context.Context, error) {
	p, err := authJWT.VerifyToken(t.Token)
	if err != nil {
		return ctx, err
	}

	return contexts.SetUserID(ctx, p.UserID), nil
}
