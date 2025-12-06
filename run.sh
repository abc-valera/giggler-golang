#!/usr/bin/env bash

# run.sh is written using an eponymous pattern for organizing project’s CLI commands.
# Read more: https://run.jotaen.net/

# Load env from dotenv files
[[ -f ./env/example.env ]] && set -a && source ./env/example.env && set +a
[[ -f ./env/.env ]] && set -a && source ./env/.env && set +a

run::webapi:dev() {
	docker compose -f compose.yaml -f compose.dev.yaml up --build
}

run::webapi:dev:stop() {
	docker compose -f compose.yaml -f compose.dev.yaml stop
}

run::webapi:dev:down() {
	echo_warning
	docker compose -f compose.yaml -f compose.dev.yaml down -v
}

run::webapi:release() {
	docker compose -f compose.yaml -f compose.release.yaml up --build
}

run::webapi:release:stop() {
	echo_warning
	docker compose -f compose.yaml -f compose.release.yaml stop
}

run::webapi:release:down() {
	echo_warning
	docker compose -f compose.yaml -f compose.release.yaml down -v
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

# TODO: adapt this to both classic/nixos/devcontainers setups
run::init-dev-tooling() {
	echo "Downloading tools and dependencies 📦 (It can take some time...)"

	go install mvdan.cc/gofumpt@latest
	go install github.com/air-verse/air@latest
	go install github.com/go-delve/delve/cmd/dlv@latest
	# gopls
	# staticcheck
	# goimports

	go mod download

	echo "Pulling and building docker images 🐳 (It can take even more time.....)"

	# TODO: pull the images instead
	run::webapi:release

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
