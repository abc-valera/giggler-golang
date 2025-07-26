package userView

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"

	"giggler-golang/src/features/user/userLevel"
	"giggler-golang/src/shared/view"
)

func Init() {
	// TODO: add helper functions in the view package
	huma.Post(view.API(), "/register", RegisterHandler)
	huma.Post(view.API(), "/login", LoginHandler)
	huma.Post(view.API(), "/refresh", RefreshHandler)
	huma.Register(view.API(), huma.Operation{
		OperationID: "profile",
		Method:      http.MethodGet,
		Path:        "/profile",
		Middlewares: huma.Middlewares{AuthMiddleware},
		Security: []map[string][]string{
			{string(userLevel.Basic): {}},
		},
	}, ProfileHandler)
}
