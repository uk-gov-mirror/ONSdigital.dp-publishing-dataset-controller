BINPATH ?= build

.PHONY: build
build:
	go build -tags 'production' -o $(BINPATH)/dp-publishing-dataset-controller

.PHONY: debug
debug:
	go build -tags 'debug' -o $(BINPATH)/dp-publishing-dataset-controller
	HUMAN_LOG=1 DEBUG=1 $(BINPATH)/dp-publishing-dataset-controller