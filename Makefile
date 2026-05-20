.PHONY: help run build test test-integration test-all lint fmt tidy clean sqlc-gen

BINARY := footy-forecast
PKG := ./...

help:
	@echo "Targets:"
	@echo "  run              - run server locally"
	@echo "  build            - build binary into ./bin"
	@echo "  test             - run unit tests (fast)"
	@echo "  test-integration - run integration tests (needs Docker)"
	@echo "  test-all         - run both"
	@echo "  lint             - run golangci-lint"
	@echo "  fmt              - format code"
	@echo "  tidy             - tidy go modules"
	@echo "  clean            - remove build artifacts"
	@echo ""
	@echo "Database:"
	@echo "  migrate-new name=<name>  - create new migration"
	@echo "  migrate-up               - apply pending migrations"
	@echo "  migrate-down             - roll back last migration"
	@echo "  migrate-status           - show migration status"
	@echo "  migrate-reset            - roll back ALL migrations (destructive)"
	@echo "  sqlc-gen        - regenerate sqlc code from queries"
	@echo ""
	@echo "AWS lifecycle:"
	@echo "  aws-status   - show current state of the AWS infrastructure"
	@echo "  aws-pause    - take final backup, stop EC2, release EIP (~\$$0.80/mo)"
	@echo "  aws-resume   - start EC2, allocate fresh EIP, wait for healthy"

run:
	go run ./cmd/server

build:
	mkdir -p bin
	go build -o bin/$(BINARY) ./cmd/server

build-linux-arm64:
	mkdir -p bin
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o bin/$(BINARY)-linux-arm64 ./cmd/server

test:
	go test -race -count=1 ./...

test-integration:
	go test -race -count=1 -tags=integration ./...

test-all: test test-integration

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

sqlc-gen:
	go tool sqlc generate

# ============================================================================
# AWS lifecycle
# ============================================================================
.PHONY: aws-pause aws-resume aws-status

AWS_PROFILE ?= hexa

aws-status:
	@AWS_PROFILE=$(AWS_PROFILE) ./scripts/aws/status.sh

aws-pause:
	@AWS_PROFILE=$(AWS_PROFILE) ./scripts/aws/pause.sh

aws-resume:
	@AWS_PROFILE=$(AWS_PROFILE) ./scripts/aws/resume.sh
