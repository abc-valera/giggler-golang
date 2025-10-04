#!/usr/bin/env bash

# run.sh is written using an eponymous pattern for organizing project’s CLI commands.
# Read more: https://run.jotaen.net/

# Load env from .env file if exists
[[ -f ./env/.env ]] && set -a && source ./env/.env && set +a

# BUILD_VERSION is a git status in a format: <branch>:<commit hash> (<clean|dirty>)
BUILD_VERSION=$(
	printf "%s:%s (%s)" \
		"$(git rev-parse --abbrev-ref HEAD)" \
		"$(git rev-parse --short=4 HEAD)" \
		"$(git status --porcelain | grep -q . && echo "dirty" || echo "clean")"
)
export BUILD_VERSION

GOPATH=$(go env GOPATH)
export GOPATH
GOMODCACHE=$(go env GOMODCACHE)
export GOMODCACHE
GOCACHE=$(go env GOCACHE)
export GOCACHE

run::restapi:dev() {
	export POSTGRES_HOST=localhost
	docker compose --profile dev up -d
	sleep 2
	air
}

run::restapi:dev:down() {
	echo_warning
	docker compose --profile dev down -v
}

run::restapi:release() {
	docker compose --profile release up --build -d
}

run::restapi:release:logs() {
	docker compose --profile release logs -f
}

run::restapi:release:down() {
	echo_warning
	docker compose --profile release down -v
}

run::build() {
	docker compose --profile release build
}

run::pprof:cpu() {
	go tool pprof -http=:3010 "$URL/debug/pprof/profile"
}

run::pprof:heap() {
	go tool pprof -http=:3010 "$URL/debug/pprof/heap"
}

run::pprof:heap:collect() {
	curl "$URL/debug/pprof/heap?gc=1" >"local/pprof/heap.$(date "+%y-%m-%d--%H-%M-%S")"
}

run::pprof:heap:diff() {
	go tool pprof -http=:3010 -diff_base "$1" "$2"
}

run::pprof:allocs() {
	go tool pprof -http=:3010 "$URL/debug/pprof/allocs"
}

run::pprof:goroutine() {
	go tool pprof -http=:3010 "$URL/debug/pprof/goroutine"
}

run::gorm:generate() {
	go run ./src/cmd/gormgen
}

run::init-dev-env() {
	echo "Downloading tools and dependencies 📦 (It can take some time...)"

	go install mvdan.cc/gofumpt@latest
	go install github.com/air-verse/air@latest
	go install github.com/go-delve/delve/cmd/dlv@latest

	go mod download

	npm install

	echo "Pulling and building docker images 🐳 (It can take even more time.....)"

	run::restapi:release

	# Create a default .env file
	cp ./env/example.env ./env/.env

	echo "Project initialized 🚀"
}

echo_warning() {
	echo "This is a dangerous command... Do you want to continue? (y/N)"
	read -r response
	if [[ "$response" =~ ^[Yy]$ ]]; then
		echo "Proceeding with the command..."
	else
		echo "Cancelled"
		exit 0
	fi
}

# "$@" represents all the arguments passed to the script
"$@"
