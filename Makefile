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
	@ rm -rf $(NAME) ./bin

glide-vc:
	@ glide-vc --only-code --no-tests --no-legal-files

fmt:
	@ go fmt $(PACKAGES)

vet:
	@ go vet $(PACKAGES)

test: clean fmt
	@ go test -v $(PACKAGES) $(ARGS)

run: clean fmt
	@ go run main.go $(ARGS)

build: \
    build-linux \
    build-darwin \
    build-windows

build-linux: clean fmt
	@ GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o bin/$(NAME)-$(VERSION)-linux-amd64
	@ cd bin; md5sum $(NAME)-$(VERSION)-linux-amd64 > $(NAME)-$(VERSION)-linux-amd64.md5

build-darwin: clean fmt
	@ GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o bin/$(NAME)-$(VERSION)-darwin-amd64
	@ cd bin; md5sum $(NAME)-$(VERSION)-darwin-amd64 > $(NAME)-$(VERSION)-darwin-amd64.md5

build-windows: clean fmt
	@ GOOS=windows GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o bin/$(NAME)-$(VERSION)-windows-amd64.exe
	@ cd bin; md5sum $(NAME)-$(VERSION)-windows-amd64.exe > $(NAME)-$(VERSION)-windows-amd64.exe.md5

