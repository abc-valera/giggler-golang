package view

import (
	"context"
	"fmt"
	"net/http"
	"net/http/pprof"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"

	"giggler-golang/src/shared/app"
	"giggler-golang/src/shared/env"
	"giggler-golang/src/shared/serviceLocator"
)

var api = func() huma.API {
	mux := http.NewServeMux()

	if env.LoadBool("IS_HTTP_PPROF_INTERFACE_ENABLED") {
		mux.HandleFunc("GET /debug/pprof/", pprof.Index)
		mux.HandleFunc("GET /debug/pprof/cmdline", pprof.Cmdline)
		mux.HandleFunc("GET /debug/pprof/profile", pprof.Profile)
		mux.HandleFunc("GET /debug/pprof/symbol", pprof.Symbol)
		mux.HandleFunc("GET /debug/pprof/trace", pprof.Trace)
	}

	mux.HandleFunc("/build-version", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(app.BuildVersion()))
	})

	config := huma.DefaultConfig("giggler-golang Docs", "1.0.0")
	config.OpenAPIPath = "/openapi"
	config.DocsPath = "/docs/"
	config.SchemasPath = "/schemas"
	config.Components.SecuritySchemes = map[string]*huma.SecurityScheme{
		"myAuth": {
			Type: "oauth2",
			Flows: &huma.OAuthFlows{
				AuthorizationCode: &huma.OAuthFlow{
					AuthorizationURL: "https://example.com/oauth/authorize",
					TokenURL:         "https://example.com/oauth/token",
					Scopes: map[string]string{
						"scope1": "Scope 1 description...",
						"scope2": "Scope 2 description...",
					},
				},
			},
		},
	}

	api := humago.New(mux, config)

	huma.Get(api, "/hello/{name}", helloHandler)

	// TODO: finish middleware setup
	// var handler http.Handler = mux
	// handler = applyCorsMiddleware(handler)
	// handler = ApplyLogMiddleware(handler)
	// handler = otelView.ApplyTracerMiddleware(handler)
	// handler = applyRecovererMiddleware(handler)

	// TODO: check if it registers correctly
	serviceLocator.Set(http.Server{
		Addr:    fmt.Sprintf(":%s", env.Load("RESTAPI_PORT")),
		Handler: api.Adapter(),
	})

	return api
}()

func Get[I, O any](path string, isAuthRequired bool, handler func(context.Context, *I) (*O, error)) {
	huma.Get(api, path, handler, func(o *huma.Operation) {
		o.Security = []map[string][]string{
			{"myAuth": {"scope1"}},
		}
	})
}

func Post[I, O any](path string, handler func(context.Context, *I) (*O, error)) {
	huma.Post(api, path, handler)
}

func Put[I, O any](path string, handler func(context.Context, *I) (*O, error)) {
	huma.Put(api, path, handler)
}

func Delete[I, O any](path string, handler func(context.Context, *I) (*O, error)) {
	huma.Delete(api, path, handler)
}
