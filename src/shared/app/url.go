package app

import (
	"net/url"

	"github.com/abc-valera/giggler-golang/src/components/env"
	"github.com/abc-valera/giggler-golang/src/components/errutil"
)

var urlVar = initURL()

func initURL() string {
	link := env.Load("URL")
	errutil.Must(url.Parse(link))
	return link
}

// URL returns the URL where the app can be accessed.
//
// Can be modified.
func URL() *url.URL {
	return errutil.Must(url.Parse(urlVar))
}
