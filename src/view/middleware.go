package view

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/abc-valera/giggler-golang/src/shared/logger"
)

func ApplyRecovererMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(rw http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					rw.WriteHeader(http.StatusInternalServerError)

					// Check if the error is of type error
					if _, ok := err.(error); !ok {
						err = fmt.Errorf("%v", err)
					}

					logger.Error("PANIC_OCCURED",
						"err", err,
						"stack", debug.Stack(),
					)
				}
			}()
			next.ServeHTTP(rw, r)
		},
	)
}

func ApplyCorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
			} else {
				next.ServeHTTP(w, r)
			}
		},
	)
}
