ROOT_DIR:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
SHELL := /bin/bash
COMMIT:=$(shell git rev-list -1 HEAD)
VERSION:=$(COMMIT)
DATE:=$(shell date -uR)

BIN_NAME:=slotalk
GOFLAGS:=-mod=readonly
GO_BUILD:=go build $(GOFLAGS)

# include files with the `// +build mock` annotation
TEST_TAGS:=-tags mock -coverprofile cover.out

.PHONY: build generate test build-all-platforms clean install-ci-tools install-local-tools licenses

install-ci-tools:
	go install github.com/google/go-licenses@c781b427440f8ea100841eefdd308e660d26d121
	go install github.com/atombender/go-jsonschema/cmd/gojsonschema@v0.11.0

install-local-tools: install-ci-tools
	go install github.com/princjef/gomarkdoc/cmd/gomarkdoc@v0.4.1
	go install github.com/tfadeyi/aloe-cli@v0.0.2
	go install github.com/goreleaser/goreleaser@v1.18.2
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.52.2

build:
	cd $(ROOT_DIR) && $(GO_BUILD) -o builds/$(BIN_NAME) .

./builds/$(BIN_NAME)-$(GOOS)-$(GOARCH):
	cd $(ROOT_DIR) && $(GO_BUILD) -o builds/$(BIN_NAME)-$(GOOS)-$(GOARCH) .

build-all-platforms:
	$(MAKE) GOOS=linux   GOARCH=amd64 ./builds/$(BIN_NAME)-linux-amd64
	$(MAKE) GOOS=darwin  GOARCH=amd64 ./builds/$(BIN_NAME)-darwin-amd64
	$(MAKE) GOOS=windows GOARCH=amd64 ./builds/$(BIN_NAME)-windows-amd64

# used to generate struct mocks
generate:
	cd $(ROOT_DIR) && go generate ./...

test: build
	cd $(ROOT_DIR) &&  go test $(GOFLAGS) $(TEST_TAGS) ./...

coverage: test
	go tool cover -html=cover.out -o cover.html
	@echo "open ./cover.html to see coverage"

clean:
	cd $(ROOT_DIR) && \
	rm -rf ./builds

lint:
	golangci-lint run

# generates the licenses used by the tool
licenses:
	rm -rf kodata
	go-licenses save . --save_path="kodata/licenses"
