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
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"

	"giggler-golang/src/features/example"
	"giggler-golang/src/features/user/userLevel"
	"giggler-golang/src/features/user/userViewWebapi"
	"giggler-golang/src/shared/buildVersion"
	"giggler-golang/src/shared/errutil/must"
	"giggler-golang/src/shared/log"
	"giggler-golang/src/shared/log/logView"
)

func main() {
	// Set global configurations here
	if must.GetEnvBool("IS_MUTEX_BLOCK_PPROF_ENABLED") {
		runtime.SetMutexProfileFraction(1)
		runtime.SetBlockProfileRate(1)
	}

	// Create a default ServeMux first, these routes won't be shown in the generated API docs
	mux := http.NewServeMux()

	if must.GetEnvBool("IS_HTTP_PPROF_INTERFACE_ENABLED") {
		mux.HandleFunc("GET /debug/pprof/", pprof.Index)
		mux.HandleFunc("GET /debug/pprof/cmdline", pprof.Cmdline)
		mux.HandleFunc("GET /debug/pprof/profile", pprof.Profile)
		mux.HandleFunc("GET /debug/pprof/symbol", pprof.Symbol)
		mux.HandleFunc("GET /debug/pprof/trace", pprof.Trace)
	}

	mux.HandleFunc("/build-version", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(buildVersion.Get()))
	})

	// Create Huma http server configuration
	humaConfig := huma.Config{
		OpenAPI: &huma.OpenAPI{
			OpenAPI: "3.1.0",
			Info: &huma.Info{
				Title:   "giggler-golang Docs",
				Version: "0.1.0",
			},
			Components: &huma.Components{
				Schemas: huma.NewMapRegistry("#/components/schemas/", huma.DefaultSchemaNamer),
				// Note, that the security schemas are defined only for documentation purposes.
				SecuritySchemes: map[string]*huma.SecurityScheme{
					string(userLevel.Basic): {
						Type:         "apiKey",
						Description:  "JWT token for authentication",
						In:           "header",
						Name:         "Authorization",
						Scheme:       "bearer",
						BearerFormat: "JWT",
					},
				},
			},
		},
		OpenAPIPath:   "/openapi",
		DocsPath:      "/docs/",
		SchemasPath:   "/schemas",
		Formats:       huma.DefaultFormats,
		DefaultFormat: "application/json",
		CreateHooks: []func(huma.Config) huma.Config{
			func(c huma.Config) huma.Config {
				// Add a link transformer to the API. This adds `Link` headers and
				// puts `$schema` fields in the response body which point to the JSON
				// Schema that describes the response structure.
				// This is a create hook so we get the latest schema path setting.
				linkTransformer := huma.NewSchemaLinkTransformer("#/components/schemas/", c.SchemasPath)
				c.OnAddOperation = append(c.OnAddOperation, linkTransformer.OnAddOperation)
				c.Transformers = append(c.Transformers, linkTransformer.Transform)
				return c
			},
		},
	}

	humaApi := humago.New(mux, humaConfig)

	// Set middlewares
	// TODO: check middlewares order
	humaApi.UseMiddleware(
		recovererMiddleware,
		corsMiddleware,
		logView.ApplyLogMiddleware,
	)

	// Register features
	example.ApplyRoutes(humaApi)
	userViewWebapi.ApplyRoutes(humaApi)

	// Create a default HTTP server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", must.GetEnv("WEBAPI_PORT")),
		Handler: humaApi.Adapter(),
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}

		log.Info("HTTP server is running",
			"HTTP server is running",
			"port", server.Addr,
			"build-version", buildVersion.Get(),
		)
	}()

	// Stop program execution until receiving an interrupt signal
	gracefulShutdown := make(chan os.Signal, 1)
	signal.Notify(gracefulShutdown, os.Interrupt)
	<-gracefulShutdown

	log.Info("received interrupt signal, shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		panic(err)
	}

	log.Info("server gracefully stopped")
}
