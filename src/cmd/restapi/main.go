package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"time"

	"giggler-golang/src/features/user/userView"
	"giggler-golang/src/shared/buildVersion"
	"giggler-golang/src/shared/env"
	"giggler-golang/src/shared/initPhase"
	"giggler-golang/src/shared/log"
	"giggler-golang/src/shared/view"
)

func init() {
	if env.LoadBool("IS_MUTEX_BLOCK_PPROF_ENABLED") {
		runtime.SetMutexProfileFraction(1)
		runtime.SetBlockProfileRate(1)
	}

	userView.Init()
}

func main() {
	initPhase.Finish()

	server := http.Server{
		Addr:    fmt.Sprintf(":%s", env.Load("RESTAPI_PORT")),
		Handler: view.API().Adapter(),
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()

	log.Info(
		"API is running",
		"port", server.Addr,
		"build-version", buildVersion.Get(),
	)

	// Stop program execution until receiving an interrupt signal
	gracefulShutdown := make(chan os.Signal, 1)
	signal.Notify(gracefulShutdown, os.Interrupt)
	<-gracefulShutdown

	log.Info("received interrupt signal, shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		panic(err)
	}

	log.Info("server gracefully stopped")
}
