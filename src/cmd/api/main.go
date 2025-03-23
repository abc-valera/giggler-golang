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

	"giggler-golang/src/shared/app"
	"giggler-golang/src/shared/env"
	"giggler-golang/src/shared/errutil"
	"giggler-golang/src/shared/log"
	"giggler-golang/src/shared/otel"
	"giggler-golang/src/shared/view/viewgen"
)

func main() {
	mux := http.NewServeMux()

	if env.LoadBool("IS_GENERAL_PPROF_ENABLED") {
		mux.HandleFunc("GET /debug/pprof/", pprof.Index)
		mux.HandleFunc("GET /debug/pprof/cmdline", pprof.Cmdline)
		mux.HandleFunc("GET /debug/pprof/profile", pprof.Profile)
		mux.HandleFunc("GET /debug/pprof/symbol", pprof.Symbol)
		mux.HandleFunc("GET /debug/pprof/trace", pprof.Trace)
	}

	if env.LoadBool("IS_MUTEX_BLOCK_PPROF_ENABLED") {
		runtime.SetMutexProfileFraction(1)
		runtime.SetBlockProfileRate(1)
	}

	mux.HandleFunc("/build-version", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(app.BuildVersion()))
	})

	// TODO: check middlewares order
	var handler http.Handler = errutil.Must(viewgen.NewServer(handler, securityHandler{}))
	handler = applyCorsMiddleware(handler)
	handler = ApplyLogMiddleware(handler)
	handler = otel.ApplyTracerMiddleware(handler)
	handler = applyRecovererMiddleware(handler)

	mux.Handle("/", handler)

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

	log.Info("App has started", "build-version", app.BuildVersion())

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
