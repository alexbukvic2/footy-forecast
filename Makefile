.PHONY: help run build test lint fmt tidy clean

BINARY := footy-forecast
PKG := ./...

help:
	@echo "Targets:"
	@echo "  run    - run server locally"
	@echo "  build  - build binary into ./bin"
	@echo "  test   - run all tests with race detector"
	@echo "  lint   - run golangci-lint"
	@echo "  fmt    - format code"
	@echo "  tidy   - tidy go modules"
	@echo "  clean  - remove build artifacts"

run:
	go run ./cmd/server

build:
	mkdir -p bin
	go build -o bin/$(BINARY) ./cmd/server

build-linux-arm64:
	mkdir -p bin
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o bin/$(BINARY)-linux-arm64 ./cmd/server

test:
	go test -race -count=1 $(PKG)

test-cover:
	go test -race -count=1 -coverprofile=coverage.out $(PKG)
	go tool cover -html=coverage.out -o coverage.html

lint:
	golangci-lint run

fmt:
	gofmt -w .
	go mod tidy

tidy:
	go mod tidy

clean:
	rm -rf bin tmp coverage.out coverage.html
