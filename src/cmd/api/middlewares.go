package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"giggler-golang/src/shared/log"
)

func applyRecovererMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(rw http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					rw.WriteHeader(http.StatusInternalServerError)

					// Check if the error is of type error
					if _, ok := err.(error); !ok {
						err = fmt.Errorf("%v", err)
					}

					log.Error("PANIC_OCCURED",
						"err", err,
						"stack", debug.Stack(),
					)
				}
			}()
			next.ServeHTTP(rw, r)
		},
	)
}

func applyCorsMiddleware(next http.Handler) http.Handler {
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

func ApplyLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Wrap the response writer so we can capture the status code and body
			wrapped := newResponseWriter(w)
			// Call the next middleware/handler in the chain
			next.ServeHTTP(wrapped, r)

			// If the status code is not explicitly set, assume 200 OK
			if wrapped.status == 0 {
				wrapped.status = 200
			}

			logMsg := []any{
				"status", wrapped.status,
				"method", r.Method,
				"path", r.URL.EscapedPath(),
				"duration(ms)", time.Since(start).Milliseconds(),
			}
			if wrapped.status < 500 {
				log.Info("REQUEST", logMsg...)
			} else {
				log.Error("REQUEST", logMsg...)
			}
		},
	)
}

// responseWriter is a minimal wrapper for http.ResponseWriter that allows the
// written HTTP status code and body to be captured for logging
type responseWriter struct {
	http.ResponseWriter

	status int
	body   []byte
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{
		ResponseWriter: w,
	}
}

// WriteHeader captures the status code before it is written to the response
func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

// Write captures the response body before it is written to the response
func (rw *responseWriter) Write(data []byte) (int, error) {
	rw.body = data
	return rw.ResponseWriter.Write(data)
}
