# Go project Makefile

.PHONY: lint format build init

# Run golangci-lint on the whole project
lint:
	golangci-lint run ./...

# Format Go
format:
	gofumpt -l -w .

# Build backend-dev image
build:
	docker compose build backend

# Start backend-dev container
init: 
	docker compose up backend