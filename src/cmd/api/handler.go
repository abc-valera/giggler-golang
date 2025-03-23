package main

import (
	"giggler-golang/src/features/auth/authView"
	"giggler-golang/src/features/joke/jokeView"
	"giggler-golang/src/features/user/userView"
	"giggler-golang/src/shared/errutil/errutilView"
	"giggler-golang/src/shared/view/viewgen"
)

type (
	authHandler = authView.Handler
	userHandler = userView.Handler
	jokeHandler = jokeView.Handler

	errorHandler = errutilView.Handler
)

var handler viewgen.Handler = struct {
	authHandler
	userHandler
	jokeHandler

	errorHandler
}{
	authHandler: authView.Handler{},
	userHandler: userView.Handler{},
	jokeHandler: jokeView.Handler{},

	errorHandler: errutilView.Handler{},
}
