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
	docker build -t sik6-api:dev .

# Start backend-dev container
init: 
	docker run --rm \
	--name sik6-api \
	-p 8080:8080 \
	-v /c/Users/frede/Desktop/sik6/backend:/app \
	sik6-api:dev
