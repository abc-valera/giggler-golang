package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"runtime"
	"time"

	"github.com/abc-valera/giggler-golang/src/features/auth/authView"
	"github.com/abc-valera/giggler-golang/src/features/joke/jokeView"
	"github.com/abc-valera/giggler-golang/src/features/user/userView"
	"github.com/abc-valera/giggler-golang/src/shared/app"
	"github.com/abc-valera/giggler-golang/src/shared/env"
	"github.com/abc-valera/giggler-golang/src/shared/errutil"
	"github.com/abc-valera/giggler-golang/src/shared/log"
	"github.com/abc-valera/giggler-golang/src/shared/otel"
	"github.com/abc-valera/giggler-golang/src/shared/view/viewgen"
)

func main() {
	mux := http.NewServeMux()

	if env.LoadBool("IS_PPROF_ENABLED") {
		// Enable mutex and block profiling
		runtime.SetMutexProfileFraction(1)
		runtime.SetBlockProfileRate(1)

		mux.HandleFunc("GET /debug/pprof/", pprof.Index)
		mux.HandleFunc("GET /debug/pprof/cmdline", pprof.Cmdline)
		mux.HandleFunc("GET /debug/pprof/profile", pprof.Profile)
		mux.HandleFunc("GET /debug/pprof/symbol", pprof.Symbol)
		mux.HandleFunc("GET /debug/pprof/trace", pprof.Trace)
	}

	mux.HandleFunc("/version", app.VersionHandler)

	// TODO: check middlewares order
	var ogenhandler http.Handler = errutil.Must(viewgen.NewServer(ogenHandler, securityHandler{}))
	ogenhandler = applyCorsMiddleware(ogenhandler)
	ogenhandler = ApplyLogMiddleware(ogenhandler)
	ogenhandler = otel.ApplyTracerMiddleware(ogenhandler)
	ogenhandler = applyRecovererMiddleware(ogenhandler)

	mux.Handle("/", ogenhandler)

	server := http.Server{
		Addr:    fmt.Sprintf(":%s", env.Load("API_PORT")),
		Handler: mux,
	}

	go func() {
		fmt.Println("Server is running on port", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()

	log.Info("App has started", "version", app.Version())

	// Stop program execution until receiving an interrupt signal
	gracefulShutdown := make(chan os.Signal, 1)
	signal.Notify(gracefulShutdown, os.Interrupt)
	<-gracefulShutdown

	// Gracefully shutdown the http server
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		panic(err)
	}
}

type (
	signHandler = authView.View
	userHandler = userView.View
	jokeHandler = jokeView.View
)

var ogenHandler viewgen.Handler = struct {
	signHandler
	userHandler
	jokeHandler

	errorHandler
}{}
