#!/bin/bash -e

VERSION="$1"

GIT_REV=$(git rev-list HEAD --count)
CWD=$(cd $(dirname $0); cd ../../; pwd)

mkdir -p _build/deb/usr/local/bin/
cp -f _releases/frep-${VERSION}-linux-amd64 _build/deb/usr/local/bin/frep

docker run --rm -it \
    -v ${CWD}:/workspace \
    --workdir /workspace \
    subchen/centos:8-dev \
    fpm -s dir -t deb --name frep --version ${VERSION} --iteration ${GIT_REV} \
        --maintainer "subchen@gmail.com" --vendor "Guoqiang Chen" --license "Apache 2" \
        --url "https://github.com/subchen/frep" \
        --description "Generate file using template" \
        -C _build/deb/ \
        --package ./_releases/


FILE=./_releases/frep_${VERSION}-${GIT_REV}_amd64.deb
sha256sum $FILE > $FILE.sha256
