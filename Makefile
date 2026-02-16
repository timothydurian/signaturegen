# Makefile for SNAP Signature Generator

.PHONY: help build run test clean docker-build docker-up docker-down docker-logs

# Variables
APP_NAME=signaturegen
DOCKER_IMAGE=snap-signature-generator
DOCKER_CONTAINER=snap-signature-generator

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build: ## Build the Go application
	@echo "Building $(APP_NAME)..."
	go build -o $(APP_NAME) .
	@echo "Build complete!"

run: ## Run the application locally
	@echo "Running $(APP_NAME)..."
	go run .

test: ## Run tests
	@echo "Running tests..."
	go test -v ./...

clean: ## Clean build artifacts
	@echo "Cleaning..."
	rm -f $(APP_NAME)
	go clean
	@echo "Clean complete!"

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE):latest .
	@echo "Docker image built!"

docker-up: ## Start services with docker-compose
	@echo "Starting services..."
	docker-compose up -d
	@echo "Services started!"

docker-down: ## Stop services with docker-compose
	@echo "Stopping services..."
	docker-compose down
	@echo "Services stopped!"

docker-logs: ## View docker-compose logs
	docker-compose logs -f

docker-restart: docker-down docker-up ## Restart docker services

docker-rebuild: docker-down docker-build docker-up ## Rebuild and restart docker services

docker-shell: ## Open shell in running container
	docker exec -it $(DOCKER_CONTAINER) sh

docker-clean: ## Remove Docker image and containers
	@echo "Cleaning Docker resources..."
	docker-compose down -v
	docker rmi $(DOCKER_IMAGE):latest 2>/dev/null || true
	@echo "Docker resources cleaned!"

install-deps: ## Download Go dependencies
	@echo "Downloading dependencies..."
	go mod download
	@echo "Dependencies downloaded!"

tidy: ## Tidy Go modules
	@echo "Tidying modules..."
	go mod tidy
	@echo "Modules tidied!"

lint: ## Run linter (requires golangci-lint)
	@echo "Running linter..."
	golangci-lint run
	@echo "Lint complete!"

dev: ## Run in development mode with hot reload (requires air)
	@echo "Starting development server..."
	air

all: clean build test ## Clean, build, and test

.DEFAULT_GOAL := help
