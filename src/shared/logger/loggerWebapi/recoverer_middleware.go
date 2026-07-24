package loggerWebapi

import (
	"net/http"
	"runtime/debug"

	"giggler-golang/src/shared/logger"

	"github.com/danielgtaylor/huma/v2"
)

func InitRecovererMiddleware(logger logger.Interface) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		defer func() {
			if err := recover(); err != nil {
				ctx.SetStatus(http.StatusInternalServerError)

				logger.Error("PANIC_OCCURED",
					"err", err,
					"stack", debug.Stack(),
				)
			}
		}()

		next(ctx)
	}
}
