#!/usr/bin/env bash

# TODO: adapt this to both classic/nixos/devcontainers setups
init-dev-tooling() {
	cp .githooks/* .git/hooks/
	cp ./template.env ./secrets/local.env

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

	echo "Project initialized 🚀"
}
