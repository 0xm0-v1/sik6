# Project Architecture Overview

This document summarises how the sik6 stack is organised so that new contributors can quickly locate the backend, frontend, and infrastructure pieces they need when shipping a feature end to end.

## High-level topology

- **Communication** – The Angular client talks to the Go API through REST endpoints that return a standard JSON envelope.
- **State management** – Each feature exposes its own store service that loads data through `ApiService` and provides signals to the view.
- **Database** – Schema changes are expressed as SQL migrations under `backend/internal/storage/postgres/migrations`. Fixtures in `backend/testdata` keep integration tests deterministic.

## Local development workflow

1. Run the stack through Docker: `cd backend && make up`. Air reloads the Go server, Angular runs in watch mode.
2. Apply migrations with `make migrate NAME=<file.sql>` and seed fixtures using `make seed-dev` when you need demo data.
3. The Angular dev server proxies `/api` requests to the Go backend (see `frontend/src/proxy.conf.json`), so frontend pages automatically call the local API.

## Feature anatomy

When introducing a new domain (e.g. `Category`):

- **Backend** – Create a domain folder under `backend/internal/<domain>` with `model.go`, `repository.go`, and `service.go`. Add HTTP handlers under `backend/internal/<domain>/transport/http` and register them in `internal/app/http.go`.
- **Database** – Place SQL migrations in `backend/internal/storage/postgres/migrations` and the corresponding repository implementation in `backend/internal/storage/postgres/repository/<domain>`.
- **Frontend** – Extend core models/adapters in `frontend/src/app/core`, expose API methods in `ApiService`, and create a feature directory (e.g. `frontend/src/app/<domain>`) containing the page component and store.

_This layered structure keeps data contracts explicit and mirrors the existing `recipe` example from database to UI._