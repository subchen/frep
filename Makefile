CWD     := $(shell pwd)
VERSION := $(shell cat VERSION)

LDFLAGS := -s -w \
           -X 'main.BuildVersion=$(VERSION)' \
           -X 'main.BuildGitBranch=$(shell git describe --all)' \
           -X 'main.BuildGitRev=$(shell git rev-list --count HEAD)' \
           -X 'main.BuildGitCommit=$(shell git rev-parse HEAD)' \
           -X 'main.BuildDate=$(shell date -u -R)'

export GO111MODULE=on

default:
	@ echo "no default target for Makefile"

clean:
	rm -rf frep ./_releases ./_build

fmt:
	go fmt ./...

lint:
	go vet ./...

build:
	go build -ldflags "$(LDFLAGS)"

build-all: clean fmt
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
		go build -a -installsuffix cgo -ldflags "$(LDFLAGS)" -o _releases/frep-$(VERSION)-linux-amd64

	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 \
		go build -a -installsuffix cgo -ldflags "$(LDFLAGS)" -o _releases/frep-$(VERSION)-linux-arm64

	GOOS=darwin GOARCH=amd64 \
		go build -ldflags "$(LDFLAGS)" -o _releases/frep-$(VERSION)-darwin-amd64

	GOOS=darwin GOARCH=arm64 \
		go build -ldflags "$(LDFLAGS)" -o _releases/frep-$(VERSION)-darwin-arm64

	GOOS=windows GOARCH=amd64 \
		go build -ldflags "$(LDFLAGS)" -o _releases/frep-$(VERSION)-windows-amd64.exe


	sha256sum _releases/frep-$(VERSION)-linux-amd64       > _releases/frep-$(VERSION)-linux-amd64
	sha256sum _releases/frep-$(VERSION)-linux-arm64       > _releases/frep-$(VERSION)-linux-arm64
	sha256sum _releases/frep-$(VERSION)-darwin-amd64      > _releases/frep-$(VERSION)-darwin-amd64
	sha256sum _releases/frep-$(VERSION)-darwin-arm64      > _releases/frep-$(VERSION)-darwin-arm64
	sha256sum _releases/frep-$(VERSION)-windows-amd64.exe > _releases/frep-$(VERSION)-windows-amd64.exe
