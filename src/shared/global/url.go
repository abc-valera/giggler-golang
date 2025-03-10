package global

import (
	"net/url"

	"github.com/abc-valera/giggler-golang/src/components/env"
	"github.com/abc-valera/giggler-golang/src/components/errutil"
)

var urlVar = errutil.Must(url.Parse(env.Load("URL")))

func URL() *url.URL {
	return urlVar
}
