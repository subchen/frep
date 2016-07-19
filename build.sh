#!/bin/sh

set -e

ROOT=$(cd $(dirname $0); pwd)
DIST=$ROOT/dist

rm -rf $DIST && mkdir -p $DIST

ldflags="-s -w"
ldflags="$ldflags -X 'main.BuildVersion=$(git rev-list HEAD --count)'"
ldflags="$ldflags -X 'main.BuildGitCommit=$(git describe --abbrev=0 --always)'"
ldflags="$ldflags -X 'main.BuildDate=$(date -u -R)'"

# build and zip
echo "building for linux"
cd $ROOT && GOOS=linux GOARCH=amd64 go build -ldflags "$ldflags" -o $DIST/frep
cd $DIST && zip -r frep-linux-amd64.zip frep

echo "building for darwin"
cd $ROOT && GOOS=darwin GOARCH=amd64 go build -ldflags "$ldflags" -o $DIST/frep
cd $DIST && zip -r frep-darwin-amd64.zip frep

echo "building for windows"
cd $ROOT && GOOS=windows GOARCH=amd64 go build -ldflags "$ldflags" -o $DIST/frep.exe
cd $DIST && zip -r frep-windows-amd64.zip frep.exe

# clean
cd $DIST && rm -rf frep frep.exe

# md5sum
cd $DIST && md5sum *.zip > md5sum.txt

