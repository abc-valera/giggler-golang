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
	"runtime/debug"
	"time"

	"github.com/abc-valera/template-golang/src/components/env"
	"github.com/abc-valera/template-golang/src/features/greetings"
	"github.com/abc-valera/template-golang/src/shared/app"
	"github.com/abc-valera/template-golang/src/shared/logger"
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

	mux.HandleFunc("/", greetings.Handler)
	mux.HandleFunc("/version", app.VersionHandler)

	var handler http.Handler = mux
	handler = applyCorsMiddleware(handler)
	handler = applyRecovererMiddleware(handler)

	server := http.Server{
		Addr:    ":" + env.Load("PORT"),
		Handler: handler,
	}

	go func() {
		fmt.Println("Server is running on port", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()

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

func applyRecovererMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(rw http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					rw.WriteHeader(http.StatusInternalServerError)

					// Check if the error is of type error
					if _, ok := err.(error); !ok {
						err = fmt.Errorf("%v", err)
					}

					logger.Error("PANIC_OCCURED",
						"err", err,
						"stack", debug.Stack(),
					)
				}
			}()
			next.ServeHTTP(rw, r)
		},
	)
}

func applyCorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
			} else {
				next.ServeHTTP(w, r)
			}
		},
	)
}
