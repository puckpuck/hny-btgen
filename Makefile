VERSION ?= $(shell cat version.txt | tr -d '\n')

.PHONY: build
build:
	@echo "Building hny-btgen $(VERSION)"
	go build -ldflags="-X 'main.VERSION=$(VERSION)'" -o hny-btgen
