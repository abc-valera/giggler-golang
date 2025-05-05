package app

import (
	"net/http"

	"github.com/abc-valera/giggler-golang/src/shared/env"
)

var buildVersion = env.Load("BUILD_VERSION")

func BuildVersion() string {
	return buildVersion
}

func BuildVersionHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(buildVersion))
}
