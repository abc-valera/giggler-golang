package app

import (
	"net/http"

	"github.com/abc-valera/giggler-golang/src/shared/env"
)

var version = env.Load("VERSION")

func Version() string {
	return version
}

func VersionHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(version))
}
