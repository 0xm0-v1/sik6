# Go project Makefile

.PHONY: lint

# Run golangci-lint on the whole project
lint:
	golangci-lint run ./...
