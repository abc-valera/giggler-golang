package logView

import (
	"fmt"
	"time"

	"github.com/danielgtaylor/huma/v2"

	"giggler-golang/src/shared/log"
)

func ApplyLogMiddleware(ctx huma.Context, next func(huma.Context)) {
	// DEBUG:
	fmt.Println("LOG")

	start := time.Now()

	// Wrap the response writer so we can capture the status code and body
	wrapped := newResponseWriter(ctx)
	// Call the next middleware/handler in the chain
	next(ctx)

	// If the status code is not explicitly set, assume 200 OK
	if wrapped.status == 0 {
		wrapped.status = 200
	}

	logMsg := []any{
		"status", wrapped.status,
		"method", ctx.Method(),
		"path", ctx.URL().Path,
		"duration(ms)", time.Since(start).Milliseconds(),
	}
	if wrapped.status < 500 {
		log.Info("REQUEST", logMsg...)
	} else {
		log.Error("REQUEST", logMsg...)
	}
}

// responseWriter is a minimal wrapper for http.ResponseWriter that allows the
// written HTTP status code and body to be captured for logging
type responseWriter struct {
	huma.Context

	status int
	body   []byte
}

func newResponseWriter(ctx huma.Context) *responseWriter {
	return &responseWriter{
		Context: ctx,
	}
}

// WriteHeader captures the status code before it is written to the response
func (w *responseWriter) WriteHeader(code int) {
	w.status = code
	w.SetStatus(code)
}

// Write captures the response body before it is written to the response
func (w *responseWriter) Write(data []byte) (int, error) {
	w.body = data
	return w.BodyWriter().Write(data)
}
