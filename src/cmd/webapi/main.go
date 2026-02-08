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

	"giggler-golang/src/core/address"
	"giggler-golang/src/core/buildVersion"
	"giggler-golang/src/core/must"
	"giggler-golang/src/features/example"
	"giggler-golang/src/features/joke/jokeData"
	"giggler-golang/src/features/user/userData"
	"giggler-golang/src/features/user/userUsecase"
	"giggler-golang/src/features/user/userViewWebapi"
	"giggler-golang/src/infra/db"
	"giggler-golang/src/infra/emailer"
	"giggler-golang/src/infra/logger"
	"giggler-golang/src/infra/logger/loggerWebapi"
)

func main() {
	// Set global configurations here
	if must.GetEnvBool("IS_MUTEX_BLOCK_PPROF_ENABLED") {
		runtime.SetMutexProfileFraction(1)
		runtime.SetBlockProfileRate(1)
	}

	buildVersionString := buildVersion.InitBuildVersion()

	urls := address.InitURLs()

	// Init infrastructure
	dbConn := db.InitPostgresConnection()
	dbConn.AutoMigrate(&userData.User{}, &jokeData.Joke{})

	var emailerInstance emailer.Interface
	switch must.GetEnv("EMAILER") {
	case "dummy":
		emailerInstance = emailer.InitDummy()
	default:
		panic(must.ErrInvalidEnvValue)
	}

	var loggerInstance logger.Interface
	switch must.GetEnv("LOGGER") {
	case "stdout":
		loggerInstance = logger.InitSlog()
	case "nop":
		loggerInstance = logger.InitNoop()
	default:
		panic(must.ErrInvalidEnvValue)
	}

	// Init usecases
	userUsecase := userUsecase.New(dbConn, emailerInstance)

	// Init undocumented APIs
	mux := http.NewServeMux()
	if must.GetEnvBool("IS_HTTP_PPROF_INTERFACE_ENABLED") {
		mux.HandleFunc("GET /debug/pprof/", pprof.Index)
		mux.HandleFunc("GET /debug/pprof/cmdline", pprof.Cmdline)
		mux.HandleFunc("GET /debug/pprof/profile", pprof.Profile)
		mux.HandleFunc("GET /debug/pprof/symbol", pprof.Symbol)
		mux.HandleFunc("GET /debug/pprof/trace", pprof.Trace)
	}
	buildVersion.InitWebapiHandler(mux, buildVersionString)

	// Init documented APIs
	humaApi := initHuma(mux, urls)
	humaApi.UseMiddleware( // TODO: check middlewares order
		loggerWebapi.InitRecovererMiddleware(loggerInstance),
		corsMiddleware,
		loggerWebapi.InitLoggerMiddleware(loggerInstance),
	)

	// Init handlers
	example.ApplyRoutes(humaApi)
	userViewWebapi.InitRoutes(humaApi, userUsecase)

	// Create a default HTTP server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", must.GetEnv("WEBAPI_PORT")),
		Handler: humaApi.Adapter(),
	}

	go func() {
		loggerInstance.Info("HTTP server is running",
			"build-version", buildVersionString,
			"local", urls.Local.String(),
			"origin", urls.Origin.String(),
		)

		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()

	// Stop program execution until receiving an interrupt signal
	gracefulShutdown := make(chan os.Signal, 1)
	signal.Notify(gracefulShutdown, os.Interrupt)
	<-gracefulShutdown

	loggerInstance.Info("received interrupt signal, shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		panic(err)
	}

	loggerInstance.Info("server gracefully stopped")
}
