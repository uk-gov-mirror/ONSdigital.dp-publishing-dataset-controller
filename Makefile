BINPATH ?= build
BUILD_TIME=$(shell date +%s)
GIT_COMMIT=$(shell git rev-parse HEAD)
VERSION ?= $(shell git tag --points-at HEAD | grep ^v | head -n 1)

.PHONY: audit
audit:
	dis-vulncheck

.PHONY: build
build:
	go build -tags 'production' -ldflags "-X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT) -X main.Version=$(VERSION)" -o $(BINPATH)/dp-publishing-dataset-controller

.PHONY: debug
debug:
	go build -tags 'debug' -ldflags "-X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT) -X main.Version=$(VERSION)" -o $(BINPATH)/dp-publishing-dataset-controller
	HUMAN_LOG=1 DEBUG=1 $(BINPATH)/dp-publishing-dataset-controller

.PHONY: debug-run
debug-run:
	HUMAN_LOG=1 DEBUG=1 go run -race $(LDFLAGS) -ldflags "-X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT) -X main.Version=$(VERSION)" main.go

.PHONY: test
test: 
	go test -race -cover ./...
