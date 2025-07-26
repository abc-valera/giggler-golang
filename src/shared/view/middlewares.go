package view

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"giggler-golang/src/shared/log"

	"github.com/danielgtaylor/huma/v2"
)

func corsMiddleware(ctx huma.Context, next func(huma.Context)) {
	// DEBUG:
	fmt.Println("CORS")

	ctx.SetHeader("Access-Control-Allow-Origin", "*")
	ctx.SetHeader("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
	ctx.SetHeader("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")

	if ctx.Method() == "OPTIONS" {
		ctx.SetStatus(http.StatusOK)
		return
	}

	next(ctx)
}

func recovererMiddleware(ctx huma.Context, next func(huma.Context)) {
	// DEBUG:
	fmt.Println("RECOVERER")

	defer func() {
		if err := recover(); err != nil {
			ctx.SetStatus(http.StatusInternalServerError)

			log.Error("PANIC_OCCURED",
				"err", err,
				"stack", debug.Stack(),
			)
		}
	}()

	next(ctx)
}
