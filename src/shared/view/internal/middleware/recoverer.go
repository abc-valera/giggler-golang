package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"giggler-golang/src/shared/log"
)

func ApplyRecovererMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)

				if _, ok := err.(error); !ok {
					err = fmt.Errorf("%v", err)
				}

				log.Error("PANIC_OCCURED",
					"err", err,
					"stack", debug.Stack(),
				)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
