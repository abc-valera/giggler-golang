package version

import (
	"net/http"
)

// Version represents the version of the application build.
// It's injected at compile time via -ldflags "-X ...".
var Version string

func InitWebapiHandler(mux *http.ServeMux) {
	mux.HandleFunc("GET /version", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(Version))
	})
}
