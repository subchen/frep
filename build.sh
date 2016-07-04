#!/bin/sh

set -e

DIR=$(cd $(dirname $0); pwd)
DIST=$DIR/dist

rm -rf $DIST

ldflags="-s -w"
ldflags="$ldflags -X 'main.BuildVersion=$(git rev-list HEAD --count)'"
ldflags="$ldflags -X 'main.BuildGitCommit=$(git describe --abbrev=0 --always)'"
ldflags="$ldflags -X 'main.BuildDate=$(date -u -R)'"

mkdir -p $DIST && GOOS=linux GOARCH=amd64 go build -ldflags "$ldflags" -o $DIST/frep-linux-amd64
mkdir -p $DIST && GOOS=darwin GOARCH=amd64 go build -ldflags "$ldflags" -o $DIST/frep-darwin-amd64
mkdir -p $DIST && GOOS=windows GOARCH=amd64 go build -ldflags "$ldflags" -o $DIST/frep-windows-amd64.exe

cd $DIST && zip -r frep-linux-amd64.zip frep-linux-amd64
cd $DIST && zip -r frep-darwin-amd64.zip frep-darwin-amd64
cd $DIST && zip -r frep-windows-amd64.zip frep-windows-amd64.exe

cd $DIST && md5sum * > md5sum.txt

