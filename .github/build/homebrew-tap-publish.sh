#!/bin/bash -e

VERSION="$1"

rm -rf ./homebrew-tap
git clone https://${GITHUB_TOKEN}@github.com/subchen/homebrew-tap.git

curl -fSL https://github.com/subchen/frep/releases/download/v1.3.12/frep-1.3.12-linux-amd64 -o ./frep
chmod +x ./frep

./frep --overwrite \
    -e VERSION=${VERSION} \
    homebrew-formula/frep.rb.gotmpl:homebrew-tap/Formula/frep.rb

cd homebrew-tap \
    && git config user.name "Guoqiang Chen" \
    && git config user.email "subchen@gmail.com" \
    && git add ./Formula/frep.rb \
    && git commit -m "Automatic update frep to ${VERSION}" \
    && git push origin master
