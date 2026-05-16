.PHONY: help run build test lint fmt tidy clean

BINARY := footy-forecast
PKG := ./...

help:
	@echo "Targets:"
	@echo "  run             - run server locally"
	@echo "  build           - build binary into ./bin"
	@echo "  test            - run all tests with race detector"
	@echo "  lint            - run golangci-lint"
	@echo "  fmt             - format code"
	@echo "  tidy            - tidy go modules"
	@echo "  clean           - remove build artifacts"
	@echo ""
	@echo "Database:"
	@echo "  migrate-new name=<name>  - create new migration"
	@echo "  migrate-up               - apply pending migrations"
	@echo "  migrate-down             - roll back last migration"
	@echo "  migrate-status           - show migration status"
	@echo "  migrate-reset            - roll back ALL migrations (destructive)"

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

# Migrations
DB_URL ?= postgres://footy:footy_dev_password@localhost:5432/footy_forecast?sslmode=disable
MIGRATIONS_DIR := ./migrations

.PHONY: migrate-new migrate-up migrate-down migrate-status migrate-reset

migrate-new:
	@if [ -z "$(name)" ]; then echo "usage: make migrate-new name=<snake_case_name>"; exit 1; fi
	go tool goose -dir $(MIGRATIONS_DIR) create $(name) sql

migrate-up:
	go tool goose -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" up

migrate-down:
	go tool goose -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" down

migrate-status:
	go tool goose -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" status

migrate-reset:
	go tool goose -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" reset
