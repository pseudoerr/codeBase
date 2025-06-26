PHONY: build run test clean docker-up docker-down migrate-up migrate-down

# Variables
BINARY_NAME=auth-service
DOCKER_COMPOSE_FILE=docker-compose.yml
MIGRATIONS_PATH=./migrations
DATABASE_URL=postgres://postgres:auth@localhost:5437/postgres?sslmode=disable

# Build the application
build:
	go build -o bin/$(BINARY_NAME) ./cmd

# Run the application locally
run:
	go run ./cmd

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out

# Docker commands
docker-up:
	docker-compose -f $(DOCKER_COMPOSE_FILE) up -d

docker-down:
	docker-compose -f $(DOCKER_COMPOSE_FILE) down

docker-logs:
	docker-compose -f $(DOCKER_COMPOSE_FILE) logs -f

# Database migrations
migrate-up:
	migrate -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL)" down

migrate-force:
	migrate -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL)" force $(VERSION)

# Development helpers
dev-setup:
	docker network create codebase-network || true
	docker-compose up -d auth-db redis
	sleep 5
	make migrate-up

# Linting
lint:
	golangci-lint run

# Format code
fmt:
	go fmt ./...

# Generate mocks (if using mockery)
mocks:
	mockery --all --output=./mocks

# Security scan
security:
	gosec ./...