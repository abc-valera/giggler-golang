package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/abc-valera/giggler-golang/src/view"
)

func main() {
	// Stop program execution until receiving an interrupt signal
	gracefulShutdown := make(chan os.Signal, 1)
	signal.Notify(gracefulShutdown, os.Interrupt)
	<-gracefulShutdown

	// Gracefully shutdown the http server
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := view.HttpServer.Shutdown(ctx); err != nil {
		panic(err)
	}
}
