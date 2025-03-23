package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"time"

	"giggler-golang/src/shared/app"
	"giggler-golang/src/shared/env"
	"giggler-golang/src/shared/log"
	"giggler-golang/src/shared/serviceLocator"
)

func init() {
	if env.LoadBool("IS_MUTEX_BLOCK_PPROF_ENABLED") {
		runtime.SetMutexProfileFraction(1)
		runtime.SetBlockProfileRate(1)
	}

	serviceLocator.Disable()
}

func main() {
	server := initServer()
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()

	log.Info(
		"API is running",
		"port", server.Addr,
		"build-version", app.BuildVersion(),
	)

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
