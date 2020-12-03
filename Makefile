CWD     := $(shell pwd)
NAME    := frep
VERSION := 1.3.11

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
	rm -rf $(NAME) ./_releases ./_build

fmt:
	go fmt ./...

lint:
	go vet ./...

build-all: \
    build-linux \
    build-darwin \
    build-windows \
    rpm \
    deb

build: build-$(shell go env GOOS)

build-linux: clean fmt
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
		go build -a -installsuffix cgo -ldflags "$(LDFLAGS)" -o _releases/$(NAME)-$(VERSION)-linux-amd64

build-darwin: clean fmt
	GOOS=darwin GOARCH=amd64 \
		go build -ldflags "$(LDFLAGS)" -o _releases/$(NAME)-$(VERSION)-darwin-amd64

build-windows: clean fmt
	GOOS=windows GOARCH=amd64 \
		go build -ldflags "$(LDFLAGS)" -o _releases/$(NAME)-$(VERSION)-windows-amd64.exe

rpm: build-linux
	mkdir -p _build/rpm/usr/local/bin/
	cp -f _releases/$(NAME)-$(VERSION)-linux-amd64 _build/rpm/usr/local/bin/$(NAME)

	docker run --rm -it \
		-v $(shell pwd):/workspace --workdir /workspace \
		subchen/centos:8-dev \
		fpm -s dir -t rpm --name $(NAME) --version $(VERSION) --iteration $(shell git rev-list HEAD --count) \
			--maintainer "subchen@gmail.com" --vendor "Guoqiang Chen" --license "Apache 2" \
			--url "https://github.com/subchen/frep" \
			--description "Generate file using template" \
			-C _build/rpm/ \
			--package ./_releases/

deb: build-linux
	mkdir -p _build/deb/usr/local/bin/
	cp -f _releases/$(NAME)-$(VERSION)-linux-amd64 _build/deb/usr/local/bin/$(NAME)

	docker run --rm -it \
		-v $(shell pwd):/workspace --workdir /workspace \
		subchen/centos:8-dev \
		fpm -s dir -t deb --name $(NAME) --version $(VERSION) --iteration $(shell git rev-list HEAD --count) \
			--maintainer "subchen@gmail.com" --vendor "Guoqiang Chen" --license "Apache 2" \
			--url "https://github.com/subchen/frep" \
			--description "Generate file using template" \
			-C _build/deb/ \
			--package ./_releases/

homebrew: build-darwin
	rm -rf homebrew-tap
	git clone https://$(GITHUB_TOKEN)@github.com/subchen/homebrew-tap.git

	go run *.go --overwrite \
	   -e VERSION=$(VERSION) \
	   homebrew-formula/frep.rb.gotmpl:homebrew-tap/Formula/frep.rb

	cd homebrew-tap \
	   && git config user.name "Guoqiang Chen" \
	   && git config user.email "subchen@gmail.com" \
	   && git add ./Formula/frep.rb \
	   && git commit -m "Automatic update frep to $(VERSION)" \
	   && git push origin master

docker:
	docker login -u subchen -p "$(DOCKER_PASSWORD)"
	docker build -t subchen/$(NAME):$(VERSION) .
	docker tag subchen/$(NAME):$(VERSION) subchen/$(NAME):latest
	docker push subchen/$(NAME):$(VERSION)
	docker push subchen/$(NAME):latest

sha256sum: build-all
	@ for f in $(shell ls ./_releases); do \
		cd $(CWD)/_releases; sha256sum "$$f" >> $$f.sha256; \
	done

release: build-all sha256sum homebrew docker
