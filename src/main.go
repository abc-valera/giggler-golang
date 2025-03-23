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
	"runtime/debug"
	"time"

	"github.com/abc-valera/giggler-golang/src/components/contexts"
	"github.com/abc-valera/giggler-golang/src/components/env"
	"github.com/abc-valera/giggler-golang/src/components/errutil"
	"github.com/abc-valera/giggler-golang/src/components/viewgen"
	"github.com/abc-valera/giggler-golang/src/features/auth/authView"
	"github.com/abc-valera/giggler-golang/src/features/joke/jokeView"
	"github.com/abc-valera/giggler-golang/src/features/user/userView"
	"github.com/abc-valera/giggler-golang/src/shared/app"
	"github.com/abc-valera/giggler-golang/src/shared/log"
	"github.com/abc-valera/giggler-golang/src/shared/otel"
)

func main() {
	mux := http.NewServeMux()

	if env.LoadBool("IS_PPROF_ENABLED") {
		// Enable mutex and block profiling
		runtime.SetMutexProfileFraction(1)
		runtime.SetBlockProfileRate(1)

		mux.HandleFunc("GET /debug/pprof/", pprof.Index)
		mux.HandleFunc("GET /debug/pprof/cmdline", pprof.Cmdline)
		mux.HandleFunc("GET /debug/pprof/profile", pprof.Profile)
		mux.HandleFunc("GET /debug/pprof/symbol", pprof.Symbol)
		mux.HandleFunc("GET /debug/pprof/trace", pprof.Trace)
	}

	mux.HandleFunc("/version", app.VersionHandler)

	mux.Handle("/", errutil.Must(viewgen.NewServer(ogenHandler, securityHandler{})))

	// TODO: check middlewares order
	var handler http.Handler = mux
	handler = applyCorsMiddleware(handler)
	handler = ApplyLogMiddleware(handler)
	handler = otel.ApplyTracerMiddleware(handler)
	handler = applyRecovererMiddleware(handler)

	server := http.Server{
		Addr:    ":" + env.Load("PORT"),
		Handler: handler,
	}

	go func() {
		fmt.Println("Server is running on port", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()

	// Stop program execution until receiving an interrupt signal
	gracefulShutdown := make(chan os.Signal, 1)
	signal.Notify(gracefulShutdown, os.Interrupt)
	<-gracefulShutdown

	// Gracefully shutdown the http server
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		panic(err)
	}
}

type (
	signHandler = authView.View
	userHandler = userView.View
	jokeHandler = jokeView.View
)

var ogenHandler viewgen.Handler = struct {
	signHandler
	userHandler
	jokeHandler

	errorHandler
}{}

type errorHandler struct{}

func (errorHandler) NewError(ctx context.Context, err error) *viewgen.ErrorSchemaStatusCode {
	var codeError viewgen.ErrorSchema
	if errutil.ErrorCode(err) == errutil.CodeInternal {
		codeError = viewgen.ErrorSchema{ErrorMessage: "Internal error"}
	} else {
		codeError = viewgen.ErrorSchema{ErrorMessage: err.Error()}
	}

	switch errutil.ErrorCode(err) {
	case errutil.CodeInvalidArgument, errutil.CodeNotFound, errutil.CodeAlreadyExists:
		return &viewgen.ErrorSchemaStatusCode{
			StatusCode: 400,
			Response:   codeError,
		}
	case errutil.CodeUnauthenticated:
		return &viewgen.ErrorSchemaStatusCode{
			StatusCode: 401,
			Response:   codeError,
		}
	case errutil.CodePermissionDenied:
		return &viewgen.ErrorSchemaStatusCode{
			StatusCode: 403,
			Response:   codeError,
		}
	default:
		log.Error("REQUEST_ERROR", "err", err.Error())
		return &viewgen.ErrorSchemaStatusCode{
			StatusCode: 500,
			Response:   codeError,
		}
	}
}

type securityHandler struct{}

func (securityHandler) HandleBearerAuth(ctx context.Context, _ string, t viewgen.BearerAuth) (context.Context, error) {
	payload, err := authView.ViewVerifyToken(t.Token)
	if err != nil {
		return ctx, err
	}

	return contexts.SetUserID(ctx, payload.UserID), nil
}

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
