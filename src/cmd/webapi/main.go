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

	user "giggler-golang/src/features/identity"
	"giggler-golang/src/features/joke"
	"giggler-golang/src/shared/address"
	"giggler-golang/src/shared/db"
	"giggler-golang/src/shared/emailer"
	"giggler-golang/src/shared/logger"
	"giggler-golang/src/shared/logger/loggerWebapi"
	"giggler-golang/src/shared/must"
	"giggler-golang/src/shared/version"
	"giggler-golang/src/shared/webapi"
	"giggler-golang/src/shared/webapi/authToken"
)

func main() {
	// Set global configurations here
	if must.GetEnvBool("GIGGLER_IS_MUTEX_BLOCK_PPROF_ENABLED") {
		runtime.SetMutexProfileFraction(1)
		runtime.SetBlockProfileRate(1)
	}

	urls := address.InitURLs()

	// Init infrastructure
	dbConn := db.InitPostgresConnection()
	dbConn.AutoMigrate(&user.User{}, &joke.Joke{})

	var emailerInstance emailer.Interface
	switch must.GetEnv("GIGGLER_EMAILER") {
	case "dummy":
		emailerInstance = emailer.InitDummy()
	default:
		panic(must.ErrInvalidEnvValue)
	}

	var loggerInstance logger.Interface
	switch must.GetEnv("GIGGLER_LOGGER") {
	case "stdout":
		loggerInstance = logger.InitSlog()
	case "nop":
		loggerInstance = logger.InitNoop()
	default:
		panic(must.ErrInvalidEnvValue)
	}

	// Init services and usecases
	authTokenService := authToken.InitService()
	userUsecase := user.New(dbConn, emailerInstance, authTokenService)
	jokeUsecase := joke.New(dbConn, authTokenService)

	// Init undocumented APIs
	mux := http.NewServeMux()
	if must.GetEnvBool("GIGGLER_IS_HTTP_PPROF_INTERFACE_ENABLED") {
		mux.HandleFunc("GET /debug/pprof/", pprof.Index)
		mux.HandleFunc("GET /debug/pprof/cmdline", pprof.Cmdline)
		mux.HandleFunc("GET /debug/pprof/profile", pprof.Profile)
		mux.HandleFunc("GET /debug/pprof/symbol", pprof.Symbol)
		mux.HandleFunc("GET /debug/pprof/trace", pprof.Trace)
	}
	version.InitWebapiHandler(mux)

	// Init documented APIs
	humaApi := initHuma(mux, urls)
	humaApi.UseMiddleware( // TODO: check middlewares order
		loggerWebapi.InitRecovererMiddleware(loggerInstance),
		corsMiddleware,
		loggerWebapi.InitLoggerMiddleware(loggerInstance),
	)

	// Init handlers
	user.InitRoutes(humaApi, userUsecase)
	joke.InitRoutes(humaApi, jokeUsecase, webapi.InitAuthMiddleware(authTokenService))

	// Create a default HTTP server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", must.GetEnv("GIGGLER_WEBAPI_PORT")),
		Handler: humaApi.Adapter(),
	}

	go func() {
		loggerInstance.Info(
			"HTTP server is running",
			"build-version", version.Version,
			"local url:", urls.Localhost.String(),
			"public url:", urls.Origin.String(),
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
