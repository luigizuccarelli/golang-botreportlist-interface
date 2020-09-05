.PHONY: all test build clean

all: clean test build

build: 
	mkdir -p build
	go build -o build -tags real ./...

build-dev:
	mkdir -p build
	GOOS=linux go build -ldflags="-s -w" -o build ./...
	chmod 755 build/microservice
	chmod 755 build/uid_entrypoint.sh

test:
	go test -v -coverprofile=tests/results/cover.out -tags fake ./...

cover:
	go tool cover -html=tests/results/cover.out -o tests/results/cover.html

clean:
	rm -rf build/*
	go clean ./...

container:
	podman build -t  tfld-docker-prd-local.repo.14west.io/servisbot-reportlist-interface:1.14.2 .

push:
	podman push tfld-docker-prd-local.repo.14west.io/servisbot-reportlist-interface:1.14.2 
