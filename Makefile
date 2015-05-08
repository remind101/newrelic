.PHONY: build test

build:
	docker build -t nra:test .

test: build
	docker run --rm nra:test go test -v ./...