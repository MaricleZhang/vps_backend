.PHONY: help build run test clean migrate-up migrate-down docker-up docker-down

help: ## Show this help
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

build: ## Build the application
	@echo "Building..."
	@go build -o bin/server cmd/server/main.go

run: ## Run the application
	@echo "Running server..."
	@go run cmd/server/main.go

test: ## Run tests
	@echo "Running tests..."
	@go test -v ./...

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf bin/

deps: ## Install dependencies
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy

migrate-up: ## Run database migrations
	@echo "Running migrations..."
	@go run cmd/server/main.go migrate up

migrate-down: ## Rollback database migrations
	@echo "Rolling back migrations..."
	@go run cmd/server/main.go migrate down

docker-up: ## Start Docker services (PostgreSQL, Redis)
	@echo "Starting Docker services..."
	@docker-compose up -d

docker-down: ## Stop Docker services
	@echo "Stopping Docker services..."
	@docker-compose down

dev: docker-up run ## Start development environment

lint: ## Run linter
	@echo "Running linter..."
	@golangci-lint run

format: ## Format code
	@echo "Formatting code..."
	@go fmt ./...
