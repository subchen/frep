#!/bin/sh

ldflags="$ldflags -s -w"
ldflags="$ldflags -X main.BuildVersion=$BuildVersion"
ldflags="$ldflags -X main.BuildGitCommit=$GitCommit"
ldflags="$ldflags -X main.BuildDate=$GitCommit"

mkdir -p dist && GOOS=linux GOARCH=amd64 go build -ldflags "$ldflags" -o dist/frep-linux-amd64
mkdir -p dist && GOOS=darwin GOARCH=amd64 go build -ldflags "$ldflags" -o dist/frep-darwin-amd64
mkdir -p dist && GOOS=windows GOARCH=amd64 go build -ldflags "$ldflags" -o dist/frep-windows-amd64

