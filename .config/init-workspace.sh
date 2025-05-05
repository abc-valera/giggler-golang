#!/usr/bin/bash

echo "Running the init-workspace script 🛠️"

# Install all the tools and dependencies
echo "Downloading tools and dependencies 📦 (It can take some time...)"

go install github.com/mgechev/revive@latest
go install mvdan.cc/gofumpt@latest

npm install

go mod download

# Pull docker images
echo "Pulling docker images 🐳 (It can take even more time.....)"

# TODO: add images building here

# Create .env files
cp ./env/example.env ./env/.env

echo "Workspace initialized 🚀"
echo "You can start coding now! 👨‍💻 / 👩‍💻"
