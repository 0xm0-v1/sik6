#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
COMPOSE_FILE="${COMPOSE_FILE:-docker-compose.yml}"
ENV_FILE="${ENV_FILE:-.env.ci}"
PROJECT_NAME="${PROJECT_NAME:-sik6-ci}"

cd "${ROOT_DIR}"

if [[ ! -f "${ENV_FILE}" ]]; then
  echo "Missing ${ENV_FILE} at repository root" >&2
  exit 1
fi

# shellcheck disable=SC1090
set -a
source "${ENV_FILE}"
set +a

compose() {
  docker compose --project-name "${PROJECT_NAME}" --env-file "${ENV_FILE}" -f "${COMPOSE_FILE}" "$@"
}

apply_migrations() {
  echo "Applying database migrations..."
  compose exec -T db sh -c '
set -e
if ! ls /migrations/*.sql >/dev/null 2>&1; then
  echo "No migrations detected."
  exit 0
fi
for file in $(ls /migrations/*.sql | sort); do
  echo "Running migration: ${file}"
  psql -v ON_ERROR_STOP=1 -U "${POSTGRES_USER}" -d "${POSTGRES_DB}" -f "${file}"
done
'
}

cleanup() {
  compose down -v --remove-orphans >/dev/null 2>&1 || true
}

trap cleanup EXIT

cleanup
compose up --wait --build
apply_migrations

curl_retry() {
  local url=$1
  curl --fail --silent --show-error --retry 5 --retry-delay 2 --retry-all-errors "${url}"
}

echo "Seeding default fixtures..."
compose run --rm api go run ./cmd/dbseed

echo "Checking API readiness..."
curl_retry "http://localhost:8080/readyz" >/tmp/sik6-readyz.json

echo "Checking API liveness..."
curl_retry "http://localhost:8080/livez" >/tmp/sik6-livez.json

echo "Checking web frontend..."
curl_retry "http://localhost:4200/" >/tmp/sik6-web.html

echo "Checking database readiness via pg_isready..."
compose exec -T db pg_isready -U "${DB_USER:-postgres}" -d "${DB_NAME:-postgres}"

echo "Running end-to-end recipes smoke test via web proxy..."
(
  cd backend
  go run ./cmd/cismoke --url http://localhost:4200/api/recipes
)

echo "All smoke checks completed."
