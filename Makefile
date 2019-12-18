BINPATH ?= build
VERSION ?= $(shell git rev-parse --short HEAD)

.PHONY: build
build:
	go build -tags 'production' -ldflags "-X main.version=$(VERSION)" -o $(BINPATH)/dp-publishing-dataset-controller 

.PHONY: debug
debug:
	go build -tags 'debug' -o $(BINPATH)/dp-publishing-dataset-controller
	HUMAN_LOG=1 DEBUG=1 $(BINPATH)/dp-publishing-dataset-controller

.PHONY: test
test: 
	go test -race -cover ./...