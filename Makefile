.PHONY: help dev-up dev-down migrate-up migrate-down migrate-create sqlc-generate run

help: ## Show this help
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'

dev-up: ## Start PostgreSQL container
	docker compose up -d

dev-down: ## Stop PostgreSQL container
	docker compose down

migrate-up: ## Run database migrations
	dbmate up

migrate-down: ## Rollback last migration
	dbmate down

migrate-create: ## Create a new migration (use name=migration_name)
	dbmate new $(name)

sqlc-generate: ## Generate sqlc code
	sqlc generate

run: ## Run the application
	go run main.go

build: ## Build the application
	go build -o bin/clockworth main.go
