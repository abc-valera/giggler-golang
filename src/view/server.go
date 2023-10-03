package view

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/abc-valera/giggler-golang/src/components/env"
	"github.com/abc-valera/giggler-golang/src/components/errutil"
	"github.com/abc-valera/giggler-golang/src/components/viewgen"
	"github.com/abc-valera/giggler-golang/src/features/auth/authView"
	"github.com/abc-valera/giggler-golang/src/features/joke/jokeView"
	"github.com/abc-valera/giggler-golang/src/features/user/userView"
	"github.com/abc-valera/giggler-golang/src/shared/logger"
	"github.com/abc-valera/giggler-golang/src/shared/otel"
)

var HttpServer http.Server

func init() {
	mux := http.NewServeMux()
	{
		mux.Handle("/ping", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("PONG")) }))

		type (
			signHandler = authView.View
			userHandler = userView.View
			jokeHandler = jokeView.View
		)

		var ogenHandler viewgen.Handler = struct {
			signHandler
			userHandler
			jokeHandler

			handler
		}{}

		mux.Handle("/", errutil.Must(viewgen.NewServer(ogenHandler, handler{})))
	}

	// TODO: check middlewares order
	var handler http.Handler = mux
	{
		handler = ApplyCorsMiddleware(handler)
		handler = logger.ApplyMiddleware(handler)
		handler = otel.ApplyTracerMiddleware(handler)
		handler = ApplyRecovererMiddleware(handler)
	}

	HttpServer = http.Server{
		Addr:    ":" + env.Load("PORT"),
		Handler: handler,
	}

	go func() {
		fmt.Println("Server is running on port", HttpServer.Addr)
		if err := HttpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()
}
