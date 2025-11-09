package buildVersion

import (
	"os/exec"
	"strings"

	"giggler-golang/src/shared/errutil/must"
	"giggler-golang/src/shared/singleton"
)

// Get returns a git status in the defined format: <branch>:<commit hash> (<clean|dirty>)
var Get = singleton.New(func() string {
	gitStatusCmd := exec.Command("sh", "-c", `
		printf "%s:%s (%s)" \
			"$(git rev-parse --abbrev-ref HEAD)" \
			"$(git rev-parse --short=4 HEAD)" \
			"$(git status --porcelain | grep -q . && echo "dirty" || echo "clean")"
	`)

	gitStatus := must.Do(gitStatusCmd.CombinedOutput())

	return strings.TrimSpace(string(gitStatus))
})
