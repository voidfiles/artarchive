PROJECT         :=artarchive
CW              :=$(shell pwd)
GOFILES         :=$(shell find . -name '*.go' -not -path './vendor/*' -path './node_modules/*')
GOPACKAGES      :=$(shell go list ./... | grep -v /vendor/| grep -v /checkers | grep -v /node_modules)
OS              := $(shell go env GOOS)
ARCH            := $(shell go env GOARCH)
CACHE           :=download-cache

BIN             := $(CW)/bin

GITHASH         :=$(shell git rev-parse --short HEAD)
GITBRANCH       :=$(shell git rev-parse --abbrev-ref HEAD)
BUILDDATE      	:=$(shell date -u +%Y%m%d%H%M)
GO_LDFLAGS		  ?= -s -w
GO_BUILD_FLAGS  :=-ldflags "${GOLDFLAGS} -X main.BuildVersion=${GITHASH} -X main.GitHash=${GITHASH} -X main.GitBranch=${GITBRANCH} -X main.BuildDate=${BUILDDATE}"
ARTIFACT_NAME   :=$(PROJECT)-$(GITHASH).tar.gz
ARTIFACT_DIR    :=$(PROJECT_DIR)/_artifacts
WORKDIR         :=$(PROJECT_DIR)/_workdir
DATA_DIR        :=$(CW)/data
MISC_DIR        :=$(CW)/_misc
WORKDIR 	      :=$(CW)/_work
SLS             :=$(CW)/node_modules/.bin/sls

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(WORKDIR)/$(PROJECT)_linux_amd64 $(GO_BUILD_FLAGS)

build:
	CGO_ENABLED=1 go build -o $(WORKDIR)/$(PROJECT)_$(OS)_$(ARCH) $(GO_BUILD_FLAGS)


dependencies:
	go get honnef.co/go/tools/cmd/megacheck
	go get github.com/alecthomas/gometalinter
	go get github.com/golang/dep/cmd/dep
	go get github.com/stretchr/testify
	go get github.com/jstemmer/go-junit-report
	dep ensure
	gometalinter --install

lint:
	echo "metalinter..."
	gometalinter --enable=goimports --enable=unparam --enable=unused --disable=golint --disable=govet $(GOPACKAGES)
	echo "megacheck..."
	megacheck $(GOPACKAGES)
	echo "golint..."
	golint -set_exit_status $(GOPACKAGES)
	echo "go vet..."
	go vet --all $(GOPACKAGES)

yarn:
	yarn install

init: dependencies yarn

clean:
	rm -fR $(WORKDIR)

test:
	CGO_ENABLED=1 go test $(GOPACKAGES)

test-race:
	CGO_ENABLED=1 go test -race $(GOPACKAGES)

deploy: clean build-linux
	$(SLS) deploy

logs_hello:
	$(SLS) logs -f hello --startTime=30m

invoke_hello:
	$(SLS) invoke -f hello

logs_indexer:
	$(SLS) logs -f indexer --startTime=30m

invoke_indexer:
	$(SLS) invoke -f indexer


run:
	$(WORKDIR)/$(PROJECT)_$(OS)_$(ARCH)

scan:
	$(WORKDIR)/$(PROJECT)_$(OS)_$(ARCH) -scan

run_slide_editor:
	(cd slide-editor && yarn start)

deploy_slide_editor:
	(cd slide-editor && yarn build)
