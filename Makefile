.PHONY: run build test lint migrate-up migrate-down

run:
	go run ./cmd/server

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o bin/server ./cmd/server

test:
	go test -race ./...

lint:
	golangci-lint run

migrate-up:
	goose -dir migrations postgres "$$DATABASE_URL" up

migrate-down:
	goose -dir migrations postgres "$$DATABASE_URL" down
