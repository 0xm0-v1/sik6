# Data fixtures

This guide explains how to add a new SQL fixture so local runs and CI share the same sample data.

> **Prerequisite:** run the database migrations first (e.g. `make migrate NAME=0001_init.sql` on a fresh volume). Seeding assumes the target tables already exist.

## 1. Create the SQL file

1. Add a file under `backend/testdata/fixtures`, e.g. `users.sql`.
2. Keep the script idempotent: clean the target tables (`TRUNCATE ...`) before inserting the rows needed by your tests.
3. Stick to plain SQL so it runs seamlessly through `pgx` (no psql-specific commands).

## 2. Expose the fixture in Go

1. Create a small helper in `backend/internal/storage/postgres/fixtures` following the existing convention:
   ```go
   const usersFixture = "users.sql"

   func SeedUsers(ctx context.Context, pool *pgxpool.Pool) error {
    	return Seed(ctx, pool, usersFixture)
   }
   ```
2. If the new fixture should become the default, update `internal/cli/dbseed/run.go` (e.g. call `fixtures.SeedUsers`).

## 3. Provide a way to run it

- Makefile (preferred, matches CI):
  ```bash
  cd backend
  make seed                  # loads default fixtures (recipes.sql)
  make seed FIXTURES=users.sql,extras.sql
  ```
- Direct CLI (useful outside Docker):
  ```bash
  cd backend
  go run ./cmd/dbseed users.sql
  ```

## 4. Validate and document

1. Execute the seeder against the Docker stack (`make up` + command above) and confirm the API returns the expected data.
2. Document any extra prerequisites here if needed.

## Reminders

- The seeder reads `.env.development` (or `DB_DSN`); ensure Postgres is running before executing it.
- Keep fixtures small -- just enough rows for smoke/integration tests. Larger datasets belong in dedicated QA scripts.

## Dev vs CI fixtures

- `recipes.sql` remains the default fixture for CI. It truncates the table and is meant for clean-state test runs.
- `recipes_dev.sql` inserts the same sample data using `ON CONFLICT DO NOTHING`, so it can be applied repeatedly without losing manual changes.
- Running `make seed-dev` (from the `backend/` directory) loads the development-friendly fixture.
