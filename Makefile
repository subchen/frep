ROOT    := $(shell pwd)
NAME    := frep
VERSION := $(shell cat VERSION)

GOPATH  := $(ROOT)/../../../../

LDFLAGS := -s -w \
           -X 'main.BuildVersion=$(VERSION)' \
           -X 'main.BuildGitRev=$(shell git rev-list HEAD --count)' \
           -X 'main.BuildGitCommit=$(shell git describe --abbrev=0 --always)' \
           -X 'main.BuildDate=$(shell date -u -R)'

PACKAGES := $(shell go list ./... | grep -v /vendor/)

clean:
	@ rm -rf $(NAME)
	@ glide vc --only-code --no-tests --no-legal-files

fmt:
	@ go fmt $(PACKAGES)

vet:
	@go vet $(PACKAGES)

test: clean fmt
	@ go test -v $(PACKAGES) $(args)

run: clean fmt
	@ go build -o $(NAME)
	@ cd test_docs; ../$(NAME) server

build: \
    build-linux \
    build-darwin \
    build-windows

build-linux: clean fmt
	@ GOOS=linux GOARCH=amd64 \
	go build -ldflags "$(LDFLAGS)" -o bin/$(VERSION)/$(NAME)-linux-amd64

build-darwin: clean fmt
	@ GOOS=darwin GOARCH=amd64 \
	go build -ldflags "$(LDFLAGS)" -o bin/$(VERSION)/$(NAME)-darwin-amd64

build-windows: clean fmt
	@ GOOS=windows GOARCH=amd64 \
	go build -ldflags "$(LDFLAGS)" -o bin/$(VERSION)/$(NAME)-windows-amd64.exe

