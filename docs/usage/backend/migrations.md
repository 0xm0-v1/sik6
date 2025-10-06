# Database migrations

Migrations for the API live in `backend/internal/storage/postgres/migrations`. They are plain SQL files executed manually whenever the schema changes.

## Add a new migration

1. Create a file in the directory above, e.g. `0002_add_categories.sql`.
2. Write the SQL statements required to reach the new schema (CREATE TABLE, ALTER TABLE, etc.).
3. Run the migration against the Docker database:
   ```bash
   cd backend
   make migrate NAME=0002_add_categories.sql
   ```
   The command wraps `docker compose exec` and executes the script inside the Postgres container.
4. Commit the SQL file so that other developers can run the same migration.

## Tips

- Run the commands from a POSIX-compatible shell (Git Bash on Windows) so that `make` can disable MSYS path conversion automatically.
- Migrations are mounted inside the container at `/migrations`. The `make migrate` target verifies that the file exists before running it.
- To replay an existing migration or inspect the database manually:
  ```bash
  MSYS_NO_PATHCONV=1 docker compose exec db \
    sh -c 'psql -U "$POSTGRES_USER" -d "$POSTGRES_DB"'
  ```
- Avoid renaming migration files once they have been merged: treat them as immutable history.

