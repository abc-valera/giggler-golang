package userViewWebapi

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"

	"giggler-golang/src/features/user/userLevel"
)

func ApplyRoutes(api huma.API) {
	// TODO: add helper functions in the view package
	huma.Post(api, "/register", RegisterHandler)
	huma.Post(api, "/login", LoginHandler)
	huma.Post(api, "/refresh", RefreshHandler)
	huma.Register(api, huma.Operation{
		OperationID: "profile",
		Method:      http.MethodGet,
		Path:        "/profile",
		Middlewares: huma.Middlewares{AuthMiddleware},
		Security: []map[string][]string{
			{string(userLevel.Basic): {}},
		},
	}, ProfileHandler)
}
