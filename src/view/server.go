package view

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/abc-valera/giggler-golang/gen/view"
	"github.com/abc-valera/giggler-golang/src/components/env"
	"github.com/abc-valera/giggler-golang/src/components/errutil"
	"github.com/abc-valera/giggler-golang/src/features/auth"
	"github.com/abc-valera/giggler-golang/src/features/joke"
	"github.com/abc-valera/giggler-golang/src/features/user"
)

var HttpServer http.Server

func init() {
	mux := http.NewServeMux()

	mux.Handle("/ping", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("PONG")) }))

	type (
		signHandler = auth.View
		userHandler = user.View
		jokeHandler = joke.View
	)

	var ogenHandler view.Handler = struct {
		signHandler
		userHandler
		jokeHandler

		handler
	}{}

	mux.Handle("/", errutil.Must(view.NewServer(ogenHandler, handler{})))

	HttpServer = http.Server{
		Addr:    ":" + env.Load("PORT"),
		Handler: mux,
	}

	go func() {
		fmt.Println("Server is running on port", HttpServer.Addr)
		if err := HttpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()
}
