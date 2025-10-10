# New Domain Workflow

Use this guide as the entry point when you need to ship a full feature (database + API + UI). It links to the canonical guides and explains the order of operations so a newcomer can follow the same path taken for the `Recipe` table.

> 📌 **Prerequisite:** read the [project architecture overview](./architecture-overview.md) to understand where backend, database, and frontend code live.

## 1. Plan the contract

1. Identify the payload that the frontend must display (list view, detail view, etc.).
2. Decide on the API endpoints you need (collection, resource, mutations) and name them following the existing REST conventions.
3. Sketch the table schema (columns, indexes, relations) and note any fixtures required for local testing.

## 2. Database & repository layer

- Create a SQL migration under `backend/internal/storage/postgres/migrations` (see [Database migrations](./usage/backend/migrations.md)).
- Implement the repository in `backend/internal/storage/postgres/repository/<domain>` and expose an interface in `backend/internal/<domain>/repository.go`.
- Add fixtures when appropriate (see [Fixtures](./usage/backend/fixtures.md)).

## 3. Backend service & HTTP transport

1. Add a service in `backend/internal/<domain>/service.go` that orchestrates validation and repository calls.
2. Implement HTTP handlers under `backend/internal/<domain>/transport/http` using the shared helpers (see [Backend API pattern](./usage/backend/api-pattern.md) and [HTTP routing guide](./usage/backend/http/routing.md)).
3. Wire the dependencies in `internal/app/http.go` and `cmd/sik6/main.go` so the new handlers are reachable.

## 4. Frontend integration

1. Extend `frontend/src/app/core/models/api.models.ts` with the new API response types.
2. Add adapters if the UI needs transformed data (`frontend/src/app/core/adapters`).
3. Expose API calls from `frontend/src/app/core/services/api.service.ts`.
4. Create the feature module directory (`frontend/src/app/<domain>`) containing:
   - A store service modelled after `root.store.ts`.
   - A standalone component and template to render the data.
   - Route registration in `frontend/src/app/app.routes.ts`.
   Detailed instructions live in the [Frontend feature workflow](./usage/frontend/feature-workflow.md) and [HTTP routing guide](./usage/frontend/http/routing.md).

_Following these steps ensures every new feature reaches parity with the existing `Recipe` workflow from database to UI._