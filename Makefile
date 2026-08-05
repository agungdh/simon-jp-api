BINARY=simon-jp-api
BUILD_DIR=bin

.PHONY: help run build build-prod clean compose-up compose-down compose-clean

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help

run: ## Run the application
	go run .

build: ## Build the application
	go build -o $(BUILD_DIR)/$(BINARY) .

build-prod: ## Build the application for production
	CGO_ENABLED=0 go build -ldflags "-s -w" -o $(BUILD_DIR)/$(BINARY) .

clean: ## Remove build artifacts
	rm -rf $(BUILD_DIR)

compose-up: ## Start docker compose services
	docker compose up -d

compose-down: ## Stop docker compose services
	docker compose down

compose-clean: ## Stop docker compose services and remove all volumes
	docker compose down -v
