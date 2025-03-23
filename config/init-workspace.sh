#!/usr/bin/bash

echo "Downloading tools and dependencies 📦 (It can take some time...)"

go install mvdan.cc/gofumpt@latest

go mod download

npm install

echo "Pulling and building docker images 🐳 (It can take even more time.....)"

task release:build
task dev:build

# Create a default .env file
cp ./env/example.env ./env/.env

echo "Workspace initialized 🚀"
