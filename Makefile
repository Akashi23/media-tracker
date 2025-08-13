.PHONY: help build run test clean docker-build docker-run docker-stop docker-logs dev-up dev-down

# Default target
help:
	@echo "Available commands:"
	@echo "  build        - Build the Go backend"
	@echo "  run          - Run the Go backend locally"
	@echo "  test         - Run tests"
	@echo "  clean        - Clean build artifacts"
	@echo "  docker-build - Build Docker images"
	@echo "  docker-run   - Run the full stack with Docker Compose"
	@echo "  docker-stop  - Stop all containers"
	@echo "  docker-logs  - Show logs from all containers"
	@echo "  dev-up       - Start development environment with hot reload"
	@echo "  dev-down     - Stop development environment"

# Build the Go backend
build:
	cd back && go build -o bin/main .

# Run the Go backend locally
run:
	cd back && go run main.go

# Run tests
test:
	cd back && go test ./...

# Clean build artifacts
clean:
	cd back && rm -rf bin/ tmp/

# Build Docker images
docker-build:
	docker-compose build

# Run the full stack with Docker Compose
docker-run:
	docker-compose up -d

# Stop all containers
docker-stop:
	docker-compose down

# Show logs from all containers
docker-logs:
	docker-compose logs -f

# Start development environment with hot reload
dev-up:
	docker-compose -f docker-compose.dev.yml up -d

# Stop development environment
dev-down:
	docker-compose -f docker-compose.dev.yml down

# Database commands
db-migrate:
	psql -h localhost -U postgres -d media_tracker -f migrations/001_initial_schema.sql

db-seed:
	psql -h localhost -U postgres -d media_tracker -f scripts/seed_data.sql

db-reset:
	psql -h localhost -U postgres -d media_tracker -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"
	psql -h localhost -U postgres -d media_tracker -f migrations/001_initial_schema.sql
	psql -h localhost -U postgres -d media_tracker -f scripts/seed_data.sql

# Install dependencies
deps:
	cd back && go mod tidy

# Format code
fmt:
	cd back && go fmt ./...

# Lint code
lint:
	cd back && golangci-lint run

# Add default user programmatically
add-user:
	cd scripts && go mod tidy && go run add_default_user.go
