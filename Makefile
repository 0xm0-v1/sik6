# Go project Makefile

.PHONY: lint format

# Run golangci-lint on the whole project
lint:
	golangci-lint run ./...

# Format Go
format:
	gofumpt -w .