package address

import (
	"net/url"

	"giggler-golang/src/shared/must"
)

type URLs struct {
	Port      uint
	Localhost *url.URL
	Origin    *url.URL // Origin is the URL where the webapi is accessible from the internet.
}

func InitURLs() URLs {
	return URLs{
		Port:      must.GetEnvUint("GIGGLER_WEBAPI_PORT"),
		Localhost: must.ParseUrl("http://localhost:" + must.GetEnv("GIGGLER_WEBAPI_PORT")),
		Origin:    must.ParseUrl(must.GetEnv("GIGGLER_ORIGIN_URL")),
	}
}
