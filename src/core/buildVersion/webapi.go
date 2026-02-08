package buildVersion

import (
	"net/http"
	"os/exec"
	"strings"

	"giggler-golang/src/core/must"
)

func InitBuildVersion() string {
	gitStatusCmd := exec.Command("sh", "-c", `
		printf "%s:%s (%s)" \
			"$(git rev-parse --abbrev-ref HEAD)" \
			"$(git rev-parse --short=4 HEAD)" \
			"$(git status --porcelain | grep -q . && echo "dirty" || echo "clean")"
	`)
	gitStatus := must.Do(gitStatusCmd.CombinedOutput())
	return strings.TrimSpace(string(gitStatus))
}

// InitWebapiHandler returns a git status in the defined format: <branch>:<commit hash> (<clean|dirty>)
func InitWebapiHandler(mux *http.ServeMux, buildVersion string) {
	mux.HandleFunc("GET /build-version", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(buildVersion))
	})
}
