#!/bin/sh

set -e

ldflags="-s -w"
ldflags="$ldflags -X 'main.BuildVersion=$(git rev-list HEAD --count)'"
ldflags="$ldflags -X 'main.BuildGitCommit=$(git describe --abbrev=0 --always)'"
ldflags="$ldflags -X 'main.BuildDate=$(date +%c)'"

mkdir -p dist && GOOS=linux GOARCH=amd64 go build -ldflags "$ldflags" -o dist/frep-linux-amd64
mkdir -p dist && GOOS=darwin GOARCH=amd64 go build -ldflags "$ldflags" -o dist/frep-darwin-amd64
mkdir -p dist && GOOS=windows GOARCH=amd64 go build -ldflags "$ldflags" -o dist/frep-windows-amd64

