# AI Service Makefile

.PHONY: help build run test clean docker-build docker-run deps fmt lint vet

# Variables
BINARY_NAME=ai-service
DOCKER_IMAGE=ai-service
VERSION=$(shell git describe --tags --always --dirty)
BUILD_TIME=$(shell date +%Y-%m-%dT%H:%M:%S%z)

# Default target
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Development
deps: ## Install dependencies
	go mod download
	go mod tidy

fmt: ## Format Go code
	go fmt ./...

lint: ## Run golangci-lint
	golangci-lint run

vet: ## Run go vet
	go vet ./...

test: ## Run tests
	go test -v ./...

test-coverage: ## Run tests with coverage
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Build
build: ## Build the binary
	CGO_ENABLED=1 go build -ldflags="-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)" -o $(BINARY_NAME) main.go

build-linux: ## Build for Linux
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags="-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)" -o $(BINARY_NAME)-linux main.go

# Run
run: ## Run the service locally
	go run main.go

run-dev: ## Run with air for hot reload (requires air: go install github.com/cosmtrek/air@latest)
	air

# Environment
env: ## Copy environment template
	cp .env.example .env
	@echo "Please edit .env file with your API keys"

# Docker
docker-build: ## Build Docker image
	docker build -t $(DOCKER_IMAGE):$(VERSION) .
	docker tag $(DOCKER_IMAGE):$(VERSION) $(DOCKER_IMAGE):latest

docker-run: ## Run Docker container
	docker run -p 8080:8080 --env-file .env $(DOCKER_IMAGE):latest

docker-compose-up: ## Start with docker-compose
	docker-compose up -d

docker-compose-down: ## Stop docker-compose
	docker-compose down

docker-compose-logs: ## View docker-compose logs
	docker-compose logs -f

# Database
db-reset: ## Reset the database
	rm -f ai_service.db

# API Testing
test-api: ## Run API tests using curl
	chmod +x examples/curl_examples.sh
	./examples/curl_examples.sh

# Swagger
swagger-gen: ## Generate Swagger documentation (requires swag: go install github.com/swaggo/swag/cmd/swag@latest)
	swag init

swagger-serve: ## Serve Swagger docs locally
	@echo "Swagger docs available at: http://localhost:8080/swagger/index.html"

# Cleanup
clean: ## Clean build artifacts
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_NAME)-linux
	rm -f coverage.out
	rm -f coverage.html
	rm -f ai_service.db

# Deployment
deploy-prod: build-linux ## Build for production deployment
	@echo "Binary ready for deployment: $(BINARY_NAME)-linux"
	@echo "Don't forget to set production environment variables"

# Quality checks
check: fmt vet lint test ## Run all quality checks

# Release
release: check build docker-build ## Prepare release
	@echo "Release $(VERSION) ready"

# Development setup
setup: deps env swagger-gen ## Initial development setup
	@echo "Setup complete!"
	@echo "1. Edit .env file with your API keys"
	@echo "2. Run 'make run' to start the service"
	@echo "3. Visit http://localhost:8080/swagger/index.html for API docs"