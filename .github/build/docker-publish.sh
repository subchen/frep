#!/bin/bash -e

CWD=$(cd $(dirname $0); cd ../../; pwd)
VERSION=$(cat $CWD/VERSION)

docker login -u subchen -p "${DOCKER_PASSWORD}"

docker build -t subchen/frep:${VERSION} .
docker tag subchen/frep:${VERSION} subchen/frep:latest
docker push subchen/frep:${VERSION}
docker push subchen/frep:latest
