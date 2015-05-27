.PHONY: build test docker docker_test

build:
	go build .

test: build
	go test -tags nra_enabled -v ./...

docker:
	docker build -t remind101/nra .

docker_test: docker
	docker run remind101/nra bash -c "cd /go/src/github.com/remind101/nra && make test"
