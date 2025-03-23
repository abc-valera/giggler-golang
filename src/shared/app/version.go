package app

import (
	"net/http"
	"os/exec"
	"strings"

	"github.com/abc-valera/giggler-golang/src/components/errutil"
)

var lastGitCommitID = strings.TrimSpace(string(errutil.Must(exec.Command("git", "rev-parse", "HEAD").Output())))

func VersionHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(lastGitCommitID))
}
