#!/usr/bin/bash

echo "Running the init-workspace script 🛠️"

# Install all the tools and dependencies
echo "Downloading tools and dependencies 📦 (It can take some time...)"

go install github.com/go-task/task/v3/cmd/task@latest
go install github.com/mgechev/revive@latest
go install mvdan.cc/gofumpt@latest

npm install --save-dev @redocly/cli@latest
npm install --save-dev @grafana/openapi-to-k6
npm install --save-dev @types/k6

go mod download

# Pull docker images
echo "Pulling docker images 🐳 (It can take even more time.....)"

docker pull golang:1.24-alpine
docker pull alpine
docker pull grafana/k6
docker pull openapitools/openapi-generator-cli

# Create docker network
docker network create giggler-network

# Create .env files
cp env/.example.dev.env env/dev.env
cp env/.example.test.env env/test.env

echo "Workspace initialized 🚀"
echo "You can start coding now! 👨‍💻 / 👩‍💻"
