.PHONY: build test clean install

BINARY_NAME=basecamp
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")

build:
	go build -ldflags "-X main.version=$(VERSION)" -o $(BINARY_NAME) ./cmd/basecamp

test:
	go test -v ./...

clean:
	rm -f $(BINARY_NAME)

install: build
	mv $(BINARY_NAME) $(GOPATH)/bin/
