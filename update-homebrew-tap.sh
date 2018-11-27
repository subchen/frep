#!/bin/bash -e

GITHUB_API_BASE="https://api.github.com/repos/subchen/homebrew-tap"
GITHUB_BRANCH="master"
GITHUB_FILE="Formula/frep.rb"
#GITHUB_TOKEN="..."

APP_VERSION="1.3.3"
APP_SHA256=$(curl -sL https://github.com/subchen/frep/releases/download/v$APP_VERSION/frep-$APP_VERSION-darwin-amd64.sha256 | cut -b 1-64)
#APP_SHA256=$(sha256sum releases/frep-$APP_VERSION-darwin-amd64 | cut -b 1-64)

cat > /tmp/formula.data << EOF
class Frep < Formula
  desc "Generate file using template from environment, arguments, json/yaml/toml config files"
  homepage "https://github.com/subchen/frep"
  url "https://github.com/subchen/frep/releases/download/v$APP_VERSION/frep-$APP_VERSION-darwin-amd64"
  version "$APP_VERSION"
  sha256 "$APP_SHA256"

  bottle :unneeded

  def install
    bin.install "frep-$APP_VERSION-darwin-amd64" => "frep"
  end

  def test
    system "frep --version"
  end
end
EOF

GITHUB_FILE_SHA=$(curl -sL "$GITHUB_API_BASE/contents/$GITHUB_FILE?ref=$GITHUB_BRANCH" -H "Authorization: token $GITHUB_TOKEN" | jq -r ".sha")

cat > /tmp/formula.post << EOF
{
    "message": "Update $GITHUB_FILE",
    "committer": {
        "name": "Guoqiang Chen",
        "email": "subchen@gmail.com"
    },
    "content": "$(cat /tmp/formula.data | base64)",
    "branch": "$GITHUB_BRANCH",
    "sha": "$GITHUB_FILE_SHA"
}
EOF

curl -s -X PUT $GITHUB_API_BASE/contents/$GITHUB_FILE \
     -H "Authorization: token $GITHUB_TOKEN" \
     -H "Content-Type: application/json" \
     --data @/tmp/formula.post

rm -f /tmp/formula.*
