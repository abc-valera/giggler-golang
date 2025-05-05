package app

import (
	"net/url"

	"github.com/abc-valera/giggler-golang/src/shared/env"
	"github.com/abc-valera/giggler-golang/src/shared/errutil"
)

var urlVar = initUrl()

// initUrl initializes the URL link from the environment variable.
//
// The function returns the string, because we don't want to store
// the globally acessible mutable pointer variable.
func initUrl() string {
	urlEnvVar := env.Load("URL")
	errutil.Must(url.Parse(urlEnvVar))
	return urlEnvVar
}

// URL returns the URL where the app can be accessed.
//
// With each call the new *url.URL is created, so it can be modified without consequences.
func URL() *url.URL {
	return errutil.Must(url.Parse(urlVar))
}
