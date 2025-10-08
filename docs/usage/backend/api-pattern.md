# Backend API Pattern

Follow this checklist to introduce a new domain (e.g. `User`) to the Go backend. Each step references the detailed guides.

## 1. Database schema

1. Add a migration under `backend/internal/storage/postgres/migrations` (e.g. `0002_add_users.sql`). See [Database migrations](./migrations.md).
2. Apply it locally:
   ```bash
   cd backend
   make migrate NAME=0002_add_users.sql
   ```
3. Optional: create fixtures in `backend/testdata/fixtures` (see [Fixtures](./fixtures.md)).

## 2. Repository layer

- Define the domain interface in `internal/<domain>/repository.go` (CRUD/pagination signatures).
- Implement it in `internal/storage/postgres/repository/<domain>/<domain>_repository.go`.
- Cover it with focused tests (fixtures or stubs).

## 3. Service layer

- Add `<domain>/service.go` mirroring `internal/recipe/service.go`.
- Encapsulate validation, domain logic, and helper methods.
- Test service behaviour with stub repositories.

## 4. HTTP transport

1. Create handlers in `internal/<domain>/transport/http/handlers.go`:
   - Use the shared envelope utilities, middleware, and pagination helpers.
   - Parse/validate request payloads and call the service.
2. Add unit tests (see `handlers_internal_test.go` in the recipe module).
3. Register the routes in `internal/app/http.go` by extending `Dependencies` and the router map.

## 5. Wiring & bootstrap

- Update `cmd/sik6/main.go` to instantiate the repository + service and pass them to `app.NewHTTPHandler`.
- Update CLI utilities (`dbseed`, `cismoke`) if they need the new domain.
- Expose fixtures through the CLI seeder when appropriate.

## 6. Validation

```bash
make -C backend restart   # rebuild containers if needed
make -C backend seed-dev  # apply dev fixtures
```

Run `go test ./...` and adjust `scripts/ci/stack-smoke.sh` if the new route must be exercised in CI.

**Remember:** keep `response.Envelope` consistent, reuse services instead of inlining repository calls, and prefer idempotent fixtures for predictable tests.
