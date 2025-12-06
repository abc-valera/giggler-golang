package public

import (
	"net/url"

	"giggler-golang/src/shared/errutil/must"
)

var Url = must.Do(url.Parse(must.GetEnv("PUBLIC_URL")))
