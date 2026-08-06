API_BINARY=simon-jp-api
WORKER_BINARY=simon-jp-worker
SCHEDULER_BINARY=simon-jp-scheduler
BUILD_DIR=bin
MIGRATIONS_DIR=internal/db/migrations

.PHONY: help run run-api run-worker run-scheduler clean-start build build-api build-worker build-scheduler build-prod build-api-prod build-worker-prod build-scheduler-prod clean compose-up compose-down compose-clean wait-db migrate-up migrate-down migrate-create

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help

run: run-api ## Run the API server

clean-start: compose-clean compose-up wait-db run-api ## Clean, start services, wait for postgres, then run the API

run-api: ## Run the API server
	go run ./cmd/api

run-worker: ## Run the background worker
	go run ./cmd/worker

run-scheduler: ## Run the cron scheduler
	go run ./cmd/scheduler

build: build-api build-worker build-scheduler ## Build all binaries

build-api: ## Build the API server
	go build -o $(BUILD_DIR)/$(API_BINARY) ./cmd/api

build-worker: ## Build the background worker
	go build -o $(BUILD_DIR)/$(WORKER_BINARY) ./cmd/worker

build-scheduler: ## Build the cron scheduler
	go build -o $(BUILD_DIR)/$(SCHEDULER_BINARY) ./cmd/scheduler

build-prod: build-api-prod build-worker-prod build-scheduler-prod ## Build all binaries for production

build-api-prod: ## Build the API server for production
	CGO_ENABLED=0 go build -ldflags "-s -w" -o $(BUILD_DIR)/$(API_BINARY) ./cmd/api

build-worker-prod: ## Build the background worker for production
	CGO_ENABLED=0 go build -ldflags "-s -w" -o $(BUILD_DIR)/$(WORKER_BINARY) ./cmd/worker

build-scheduler-prod: ## Build the cron scheduler for production
	CGO_ENABLED=0 go build -ldflags "-s -w" -o $(BUILD_DIR)/$(SCHEDULER_BINARY) ./cmd/scheduler

clean: ## Remove build artifacts
	rm -rf $(BUILD_DIR)

compose-up: ## Start docker compose services
	docker compose up -d

compose-down: ## Stop docker compose services
	docker compose down

compose-clean: ## Stop docker compose services and remove all volumes
	docker compose down -v

wait-db: ## Wait until postgres is ready
	@until docker compose exec -T postgres pg_isready -U admin -d simonjp >/dev/null 2>&1; do sleep 1; done

migrate-up: ## Apply all migrations
	go run ./cmd/migrate up

migrate-down: ## Rollback the last migration
	go run ./cmd/migrate down

migrate-create: ## Create a new migration: make migrate-create NAME=add_users
	goose -dir $(MIGRATIONS_DIR) create $(NAME) sql
