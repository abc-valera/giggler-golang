package main

import (
	"net/http"

	"giggler-golang/src/core/address"
	"giggler-golang/src/features/user/userData"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
)

func initHuma(mux *http.ServeMux, urls address.URLs) huma.API {
	// Create Huma http server configuration
	humaConfig := huma.Config{
		OpenAPI: &huma.OpenAPI{
			OpenAPI: "3.1.0",
			Info: &huma.Info{
				Title:   "giggler-golang Docs",
				Version: "0.1.0",
			},
			// TODO: maybe apply conditionaly based on environment (dev/release)
			Servers: []*huma.Server{
				{
					URL:         urls.Local.String(),
					Description: "Local dev server",
				},
				{
					URL:         urls.Origin.String(),
					Description: "Public production server",
				},
			},
			Components: &huma.Components{
				Schemas: huma.NewMapRegistry("#/components/schemas/", huma.DefaultSchemaNamer),
				// Note, that the security schemas are defined only for documentation purposes.
				SecuritySchemes: map[string]*huma.SecurityScheme{
					string(userData.Basic): {
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
		DocsPath:      "/",
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

	return humaApi
}

func corsMiddleware(ctx huma.Context, next func(huma.Context)) {
	ctx.SetHeader("Access-Control-Allow-Origin", "*")
	ctx.SetHeader("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
	ctx.SetHeader("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")

	if ctx.Method() == "OPTIONS" {
		ctx.SetStatus(http.StatusOK)
		return
	}

	next(ctx)
}
