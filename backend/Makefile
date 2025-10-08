# Go project Makefile

COMPOSE = docker compose --project-directory .. -f ../docker-compose.yml --env-file ../.env.development
MIGRATIONS_DIR = internal/storage/postgres/migrations

.PHONY: lint format build up up-detach down logs ps restart migrate seed seed-dev seed-ci

# Run golangci-lint on the whole project
lint:
	golangci-lint run ./...

# Format Go sources
format:
	gofumpt -l -w .

# Build API and web images defined in docker-compose.yml
build:
	$(COMPOSE) build api web

# Start the full stack in the foreground (API, web, db)
up:
	$(COMPOSE) up --build

# Start the stack in detached mode
up-detach:
	$(COMPOSE) up --build -d

# Stop and remove all containers in the stack
down:
	$(COMPOSE) down

# Tail logs from the running stack
logs:
	$(COMPOSE) logs -f

# Show container status
ps:
	$(COMPOSE) ps

# Recreate the stack from scratch
restart:
	$(COMPOSE) down
	$(COMPOSE) up --build

# Apply a SQL migration located in $(MIGRATIONS_DIR)
migrate:
	$(if $(NAME),,$(error NAME is required. usage: make migrate NAME=0002_add_table.sql))
	$(if $(wildcard $(MIGRATIONS_DIR)/$(NAME)),,$(error Missing $(MIGRATIONS_DIR)/$(NAME)))
	MSYS_NO_PATHCONV=1 $(COMPOSE) exec -T db sh -c 'psql -U "$$POSTGRES_USER" -d "$$POSTGRES_DB" -f /migrations/$(NAME)'

# Seed database fixtures
# Usage:
#   make seed                  # loads default fixtures (recipes.sql)
#   make seed FIXTURES=foo.sql # loads specific fixture files
seed:
	$(COMPOSE) run --rm api go run ./cmd/dbseed $(FIXTURES)

# Seed development-friendly datasets (idempotent fixtures).
seed-dev:
	$(COMPOSE) run --rm api go run ./cmd/dbseed recipes_dev.sql

# Explicit alias for CI/test fixtures.
seed-ci: seed

