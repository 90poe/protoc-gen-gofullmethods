include linting.mk

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
	go build -o bin/protoc-gen-gofullmethods main.go

.PHONY: unit-test
unit-test:
	go test -coverpkg=./... -coverprofile=unit_coverage.out ./...

.PHONY: update-example
update-example: clean-example download-example-proto build	
	buf generate

clean-example:
	rm -rf example/idl/*

download-example-proto:
	buf export --output example/idl "https://github.com/bufbuild/buf-tour.git#branch=main,subdir=start/petapis"
