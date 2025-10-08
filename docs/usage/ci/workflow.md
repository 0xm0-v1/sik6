# CI workflow overview

The GitHub Actions workflow (`.github/workflows/ci-bootstrap.yml`) validates pull requests in two stages.

## 1. Quality checks

- Ensure Go formatting matches `gofumpt` (`gofumpt -l` must report no files).
- Run `golangci-lint` across the backend module (`./backend`).
- Build the API binary via `go build ./cmd/sik6`.
- Install frontend dependencies (`npm ci`), run ESLint (`npm run lint:check`), and build the Angular application (`npm run build -- --configuration production`).

## 2. Docker smoke tests

- Pull the PostgreSQL base image and build the `api` and `web` images using `docker-compose.yml` (the CI job sets `APP_ENV_FILE=.env.ci` so the API container consumes the CI-specific env file).
- Execute `scripts/ci/stack-smoke.sh`, which:
  - launches the Compose stack using `.env.ci`,
  - waits for all health checks,
  - applies every SQL migration found in `backend/internal/storage/postgres/migrations`,
  - seeds the default fixtures inside the API container,
  - probes API (`/readyz`, `/livez`), web (`/`), and database (`pg_isready`) services,
  - runs the Go smoke CLI (`go run ./cmd/cismoke`) against `/api/recipes` through the web proxy to confirm data flows from web -> API -> database.
- On failure, the workflow prints the Compose logs to aid debugging.

## Reproducing locally

```bash
# From the repository root
export APP_ENV_FILE=.env.ci
docker compose --project-name sik6-ci --env-file .env.ci down -v
docker compose --project-name sik6-ci --env-file .env.ci pull db
docker compose --project-name sik6-ci --env-file .env.ci build api web
bash scripts/ci/stack-smoke.sh
```

```bash
# Remove variable
unset APP_ENV_FILE
```
