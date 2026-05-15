# ============================================================
#  wallet-service
# ============================================================

# --- Config --------------------------------------------------
APP_ENV   ?= local
DB_URL    ?= postgres://wallet:wallet@localhost:5432/wallet?sslmode=disable
MIGRATION_DIR := ./migrations

# --- Tools ---------------------------------------------------
GOOSE  := go tool goose
SQLC   := sqlc
BUF    := buf

.PHONY: help build run test test-v test-integration lint generate \
        migrate migrate-down migrate-reset migrate-status \
        docker-up docker-down

# Default target
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

# --- Build & Run ---------------------------------------------

build: ## Compile the server binary to ./bin/server
	go build -o bin/server ./cmd/server

run: ## Run the server locally (APP_ENV=local)
	APP_ENV=$(APP_ENV) go run ./cmd/server

test: ## Run all unit tests
	go test ./...

test-v: ## Run all tests with verbose output
	go test -v ./...

test-integration: ## Run integration tests with verbose output (requires Docker)
	go test -v ./internal/service/...

lint: ## Run go vet (install golangci-lint for richer linting)
	go vet ./...

# --- Code Generation -----------------------------------------

generate: ## Regenerate protobuf (buf) and SQL (sqlc)
	$(BUF) generate
	$(SQLC) generate

# --- Migrations ----------------------------------------------
# goose is tracked as a Go tool (go tool goose).
# GOOSE_DRIVER and GOOSE_DBSTRING are picked up automatically.

GOOSE_ENV := GOOSE_DRIVER=postgres GOOSE_DBSTRING="$(DB_URL)" GOOSE_MIGRATION_DIR=$(MIGRATION_DIR)

migrate: ## Apply all pending migrations (goose up)
	$(GOOSE_ENV) $(GOOSE) up

migrate-down: ## Roll back the last migration (goose down)
	$(GOOSE_ENV) $(GOOSE) down

migrate-reset: ## Roll back all migrations (goose reset)
	$(GOOSE_ENV) $(GOOSE) reset

migrate-status: ## Show migration status (goose status)
	$(GOOSE_ENV) $(GOOSE) status

# --- Docker --------------------------------------------------

docker-up: ## Start Postgres and run migrations
	docker compose up -d postgres
	@echo "Waiting for Postgres to be ready..."
	@until docker compose exec postgres pg_isready -U wallet > /dev/null 2>&1; do sleep 1; done
	$(GOOSE_ENV) $(GOOSE) up

docker-down: ## Stop and remove containers
	docker compose down
