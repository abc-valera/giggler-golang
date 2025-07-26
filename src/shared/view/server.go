package view

import (
	"net/http"
	"net/http/pprof"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"

	"giggler-golang/src/features/user/userLevel"
	"giggler-golang/src/shared/buildVersion"
	"giggler-golang/src/shared/env"
	"giggler-golang/src/shared/log/logView"
)

var API = func() func() huma.API {
	mux := http.NewServeMux()

	if env.LoadBool("IS_HTTP_PPROF_INTERFACE_ENABLED") {
		mux.HandleFunc("GET /debug/pprof/", pprof.Index)
		mux.HandleFunc("GET /debug/pprof/cmdline", pprof.Cmdline)
		mux.HandleFunc("GET /debug/pprof/profile", pprof.Profile)
		mux.HandleFunc("GET /debug/pprof/symbol", pprof.Symbol)
		mux.HandleFunc("GET /debug/pprof/trace", pprof.Trace)
	}

	mux.HandleFunc("/build-version", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(buildVersion.Get()))
	})

	config := huma.DefaultConfig("giggler-golang Docs", "1.0.0")
	config.OpenAPIPath = "/openapi"
	config.DocsPath = "/docs/"
	config.SchemasPath = "/schemas"
	// Note, that the security schemas are defined only for documentation purposes.
	config.Components.SecuritySchemes = map[string]*huma.SecurityScheme{
		string(userLevel.Basic): {
			Type:         "apiKey",
			Description:  "JWT token for authentication",
			In:           "header",
			Name:         "Authorization",
			Scheme:       "bearer",
			BearerFormat: "JWT",
		},
	}

	api := humago.New(mux, config)

	// TODO: check middlewares order
	api.UseMiddleware(
		recovererMiddleware,
		corsMiddleware,
		logView.ApplyLogMiddleware,
	)

	return func() huma.API {
		return api
	}
}()
