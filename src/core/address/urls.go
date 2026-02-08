package address

import (
	"net/url"

	"giggler-golang/src/core/must"
)

type URLs struct {
	Local *url.URL
	// Origin is the URL of the webapp's public server
	Origin *url.URL
}

func InitURLs() URLs {
	return URLs{
		Local:  must.ParseUrl("http://localhost:" + must.GetEnv("WEBAPI_PORT")),
		Origin: must.ParseUrl(must.GetEnv("ORIGIN_URL")),
	}
}
