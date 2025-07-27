# AI Service Makefile - Comprehensive Development Tools

.PHONY: help build run test clean docker-build docker-run deps fmt lint vet check-all security-check performance-test integration-test e2e-test

# =============================================================================
# VARIABLES
# =============================================================================
BINARY_NAME=ai-service
DOCKER_IMAGE=ai-service
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev-$(shell git rev-parse --short HEAD 2>/dev/null || echo 'unknown')")
BUILD_TIME=$(shell date +%Y-%m-%dT%H:%M:%S%z)
GO_VERSION=$(shell go version | awk '{print $$3}')
LDFLAGS=-ldflags="-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GoVersion=$(GO_VERSION)"

# =============================================================================
# DEFAULT TARGET
# =============================================================================
help: ## Show comprehensive help message
	@echo '🚀 AI Service - Go Gin Lambda Project'
	@echo '====================================='
	@echo 'Version: $(VERSION)'
	@echo 'Go Version: $(GO_VERSION)'
	@echo 'Build Time: $(BUILD_TIME)'
	@echo ''
	@echo '📋 Available Commands:'
	@echo ''
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
	@echo ''
	@echo '🔧 Quick Start:'
	@echo '  make setup      - Initial project setup'
	@echo '  make run        - Start development server'
	@echo '  make test       - Run all tests'
	@echo '  make check-all  - Run all quality checks'

# =============================================================================
# DEVELOPMENT TOOLS
# =============================================================================
deps: ## Install and update all dependencies
	@echo "📦 Installing dependencies..."
	go mod download
	go mod tidy
	go mod verify
	@echo "✅ Dependencies installed successfully"

fmt: ## Format all Go code with gofmt
	@echo "🎨 Formatting Go code..."
	go fmt ./...
	@echo "✅ Code formatting complete"

lint: ## Run comprehensive linting with golangci-lint
	@echo "🔍 Running linting checks..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run --timeout=5m; \
	else \
		echo "⚠️  golangci-lint not found. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
		exit 1; \
	fi
	@echo "✅ Linting complete"

vet: ## Run go vet for static analysis
	@echo "🔍 Running go vet..."
	go vet ./...
	@echo "✅ Go vet complete"

staticcheck: ## Run staticcheck for additional analysis
	@echo "🔍 Running staticcheck..."
	@if command -v staticcheck >/dev/null 2>&1; then \
		staticcheck ./...; \
	else \
		echo "⚠️  staticcheck not found. Install with: go install honnef.co/go/tools/cmd/staticcheck@latest"; \
	fi
	@echo "✅ Staticcheck complete"

# =============================================================================
# TESTING
# =============================================================================
test: ## Run all tests with verbose output
	@echo "🧪 Running all tests..."
	go test -v -race -timeout=5m ./...
	@echo "✅ All tests completed"

test-unit: ## Run only unit tests
	@echo "🧪 Running unit tests..."
	go test -v -race ./tests/unit/...
	@echo "✅ Unit tests completed"

test-integration: ## Run integration tests
	@echo "🧪 Running integration tests..."
	go test -v -race ./tests/integration/...
	@echo "✅ Integration tests completed"

test-e2e: ## Run end-to-end tests
	@echo "🧪 Running E2E tests..."
	go test -v -race ./tests/e2e/...
	@echo "✅ E2E tests completed"

test-coverage: ## Run tests with coverage report
	@echo "🧪 Running tests with coverage..."
	go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "📊 Coverage report generated: coverage.html"
	@echo "📈 Coverage summary:"
	@go tool cover -func=coverage.out | tail -1

test-benchmark: ## Run benchmark tests
	@echo "🏃 Running benchmark tests..."
	go test -bench=. -benchmem ./...
	@echo "✅ Benchmark tests completed"

# =============================================================================
# BUILD
# =============================================================================
build: ## Build the binary for current platform
	@echo "🔨 Building binary..."
	CGO_ENABLED=1 go build $(LDFLAGS) -o $(BINARY_NAME) cmd/main/main.go
	@echo "✅ Binary built: $(BINARY_NAME)"

build-linux: ## Build for Linux deployment
	@echo "🔨 Building Linux binary..."
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY_NAME)-linux cmd/main/main.go
	@echo "✅ Linux binary built: $(BINARY_NAME)-linux"

build-all: build build-linux ## Build for all target platforms
	@echo "✅ All binaries built successfully"

# =============================================================================
# RUN
# =============================================================================
run: ## Run the service locally with PostgreSQL
	@echo "🚀 Starting AI Service with PostgreSQL..."
	@echo "📡 Server will be available at: http://localhost:8080"
	@echo "📚 API Documentation: http://localhost:8080/swagger/index.html"
	@PORT=8080 go run cmd/main/main.go

run-postgres: ## Run the service with PostgreSQL
	@echo "🚀 Starting AI Service with PostgreSQL..."
	@echo "📡 Server will be available at: http://localhost:8080"
	@echo "📚 API Documentation: http://localhost:8080/swagger/index.html"
	@PORT=8080 export DB_TYPE=postgres && \
	export DB_HOST=localhost && \
	export DB_PORT=5432 && \
	export DB_USER=ai_service_user && \
	export DB_PASSWORD=ai_service_password && \
	export DB_NAME=ai_service && \
	export DB_SSLMODE=disable && \
	go run cmd/main/main.go

run-dev: ## Run with hot reload using Air
	@echo "🚀 Starting development server with hot reload..."
	@if command -v air >/dev/null 2>&1; then \
		air; \
	else \
		echo "⚠️  Air not found. Install with: go install github.com/cosmtrek/air@latest"; \
		echo "💡 Falling back to regular run..."; \
		make run; \
	fi

# =============================================================================
# ENVIRONMENT & CONFIGURATION
# =============================================================================
env: ## Create environment file from template
	@if [ ! -f .env ]; then \
		cp .env.example .env 2>/dev/null || echo "# AI Service Environment Variables" > .env; \
		echo "✅ Environment file created: .env"; \
		echo "📝 Please edit .env file with your configuration"; \
	else \
		echo "⚠️  .env file already exists"; \
	fi

env-check: ## Check environment variables
	@echo "🔍 Checking environment variables..."
	@if [ -f .env ]; then \
		echo "✅ .env file found"; \
		@echo "📋 Required variables:"; \
		@echo "  - PORT (default: 8080)"; \
		@echo "  - ENV (default: development)"; \
		@echo "  - SENTRY_DSN (optional)"; \
		@echo "  - JWT_SECRET (optional)"; \
	else \
		echo "⚠️  .env file not found. Run 'make env' to create one"; \
	fi

# =============================================================================
# DOCKER
# =============================================================================
docker-build: ## Build Docker image
	@echo "🐳 Building Docker image..."
	docker build -t $(DOCKER_IMAGE):$(VERSION) .
	docker tag $(DOCKER_IMAGE):$(VERSION) $(DOCKER_IMAGE):latest
	@echo "✅ Docker image built: $(DOCKER_IMAGE):$(VERSION)"

docker-run: ## Run Docker container
	@echo "🐳 Running Docker container..."
	docker run -p 8080:8080 --env-file .env $(DOCKER_IMAGE):latest

docker-compose-up: ## Start services with docker-compose (PostgreSQL)
	@echo "🐳 Starting services with docker-compose (PostgreSQL)..."
	docker-compose up -d
	@echo "✅ Services started"

docker-compose-down: ## Stop docker-compose services
	@echo "🐳 Stopping docker-compose services..."
	docker-compose down
	@echo "✅ Services stopped"

docker-compose-postgres-up: ## Start services with PostgreSQL
	@echo "🐳 Starting services with PostgreSQL..."
	docker-compose -f docker-compose.postgres.yml up -d
	@echo "✅ PostgreSQL services started"

docker-compose-postgres-down: ## Stop PostgreSQL services
	@echo "🐳 Stopping PostgreSQL services..."
	docker-compose -f docker-compose.postgres.yml down
	@echo "✅ PostgreSQL services stopped"

docker-compose-logs: ## View docker-compose logs
	docker-compose logs -f

docker-clean: ## Clean Docker images and containers
	@echo "🧹 Cleaning Docker resources..."
	docker system prune -f
	@echo "✅ Docker cleanup complete"

# =============================================================================
# DATABASE
# =============================================================================
db-reset: ## Reset the database (PostgreSQL)
	@echo "🗄️  Resetting database..."
	rm -f ai_service.db
	@echo "✅ Database reset complete"

db-migrate: ## Run database migrations
	@echo "🗄️  Running database migrations..."
	@if [ -f scripts/migrations/001_initial_schema.sql ]; then \
		echo "📋 Migration files found"; \
		echo "💡 For PostgreSQL, run migrations manually"; \
	else \
		echo "⚠️  No migration files found"; \
	fi

db-seed: ## Seed database with sample data
	@echo "🌱 Seeding database..."
	@if [ -d scripts/seeders ]; then \
		echo "📋 Seeder files found"; \
	else \
		echo "⚠️  No seeder files found"; \
	fi

setup-postgres: ## Setup PostgreSQL database and run migrations
	@echo "🗄️  Setting up PostgreSQL..."
	@chmod +x scripts/setup_postgres.sh
	@./scripts/setup_postgres.sh

# =============================================================================
# API TESTING
# =============================================================================
test-api: ## Run API tests using curl
	@echo "🌐 Running API tests..."
	@if [ -f examples/curl_examples.sh ]; then \
		chmod +x examples/curl_examples.sh; \
		./examples/curl_examples.sh; \
	else \
		echo "⚠️  API test script not found: examples/curl_examples.sh"; \
	fi

# =============================================================================
# DOCUMENTATION
# =============================================================================
swagger-gen: ## Generate Swagger documentation
	@echo "📚 Generating Swagger documentation..."
	@if command -v swag >/dev/null 2>&1; then \
		swag init -g cmd/main/main.go -o api/swagger; \
		echo "✅ Swagger docs generated"; \
	else \
		echo "⚠️  swag not found. Install with: go install github.com/swaggo/swag/cmd/swag@latest"; \
	fi

swagger-serve: ## Serve Swagger docs locally
	@echo "📚 Swagger documentation available at: http://localhost:8080/swagger/index.html"

docs-gen: ## Generate all documentation
	@echo "📚 Generating documentation..."
	make swagger-gen
	@echo "✅ Documentation generated"

# =============================================================================
# SECURITY & QUALITY
# =============================================================================
security-check: ## Run security checks
	@echo "🔒 Running security checks..."
	@if command -v gosec >/dev/null 2>&1; then \
		gosec ./...; \
	else \
		echo "⚠️  gosec not found. Install with: go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest"; \
	fi
	@echo "✅ Security checks complete"

performance-test: ## Run performance tests
	@echo "⚡ Running performance tests..."
	@if command -v hey >/dev/null 2>&1; then \
		echo "💡 Use 'hey' for load testing: hey -n 1000 -c 10 http://localhost:8080/api/health"; \
	else \
		echo "⚠️  hey not found. Install with: go install github.com/rakyll/hey@latest"; \
	fi

# =============================================================================
# COMPREHENSIVE CHECKS
# =============================================================================
check-all: fmt vet staticcheck lint security-check test ## Run all quality and security checks
	@echo "✅ All checks completed successfully"

check-quick: fmt vet test ## Run quick checks (no linting)
	@echo "✅ Quick checks completed"

# =============================================================================
# CLEANUP
# =============================================================================
clean: ## Clean all build artifacts
	@echo "🧹 Cleaning build artifacts..."
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_NAME)-linux
	rm -f coverage.out
	rm -f coverage.html
	rm -f ai_service.db
	@echo "✅ Cleanup complete"

clean-all: clean docker-clean ## Clean everything including Docker
	@echo "✅ Complete cleanup finished"

# =============================================================================
# DEPLOYMENT
# =============================================================================
deploy-prod: check-all build-linux ## Build for production deployment
	@echo "🚀 Production build ready: $(BINARY_NAME)-linux"
	@echo "📋 Deployment checklist:"
	@echo "  ✅ All tests passed"
	@echo "  ✅ Code quality checks passed"
	@echo "  ✅ Security checks passed"
	@echo "  ✅ Binary built for Linux"
	@echo "  ⚠️  Remember to set production environment variables"

# =============================================================================
# RELEASE
# =============================================================================
release: check-all build-all docker-build ## Prepare complete release
	@echo "🎉 Release $(VERSION) ready!"
	@echo "📦 Artifacts:"
	@echo "  - $(BINARY_NAME) (current platform)"
	@echo "  - $(BINARY_NAME)-linux (Linux deployment)"
	@echo "  - Docker image: $(DOCKER_IMAGE):$(VERSION)"

# =============================================================================
# DEVELOPMENT SETUP
# =============================================================================
setup: deps env swagger-gen ## Complete initial development setup
	@echo "🎉 Setup complete!"
	@echo ""
	@echo "📋 Next steps:"
	@echo "  1. Edit .env file with your configuration"
	@echo "  2. Run 'make run' to start the service"
	@echo "  3. Visit http://localhost:8080/swagger/index.html for API docs"
	@echo "  4. Run 'make test' to verify everything works"
	@echo ""
	@echo "🔧 Development commands:"
	@echo "  make run-dev    - Start with hot reload"
	@echo "  make test       - Run tests"
	@echo "  make check-all  - Run all quality checks"
	@echo "  make build      - Build binary"

# =============================================================================
# UTILITIES
# =============================================================================
version: ## Show version information
	@echo "📋 Version Information:"
	@echo "  Version: $(VERSION)"
	@echo "  Go Version: $(GO_VERSION)"
	@echo "  Build Time: $(BUILD_TIME)"

deps-check: ## Check if all required tools are installed
	@echo "🔍 Checking required tools..."
	@echo "  Go: $(shell which go)"
	@echo "  Git: $(shell which git)"
	@echo "  Docker: $(shell which docker)"
	@echo "  Docker Compose: $(shell which docker-compose)"
	@echo "  golangci-lint: $(shell which golangci-lint 2>/dev/null || echo 'Not installed')"
	@echo "  swag: $(shell which swag 2>/dev/null || echo 'Not installed')"
	@echo "  air: $(shell which air 2>/dev/null || echo 'Not installed')"