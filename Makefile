include linting.mk

export GO111MODULE=on

.PHONY: ci
ci: build test lint

.PHONY: all
all: deps build test lint

.PHONY: test
test: unit-test

.PHONY: deps
deps: tidy

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: verify
verify:
	go mod verify

.PHONY: build
build: test lint
	go build ./...

.PHONY: unit-test
unit-test:
	go test -coverpkg=./... -coverprofile=unit_coverage.out ./...
