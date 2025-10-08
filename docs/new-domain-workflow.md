# New Domain Workflow

Use this page as the entry point when you have to build a full feature (database + API + UI).
It only links to the canonical guides; refer to them for the detailed steps.

## Backend

- **Migrations & fixtures:** [Database migrations](./usage/backend/migrations.md) / [Fixtures](./usage/backend/fixtures.md)
- **Service/handler structure:** [Backend API pattern](./usage/backend/api-pattern.md)

## Frontend

- **Environment & API client:** [Frontend environment configuration](./usage/frontend/config/environment.md)
- **Feature implementation (models, store, component, route):** [Frontend feature workflow](./usage/frontend/feature-workflow.md)

## Validation

- Backend tests: `go test ./...`
- Frontend lint/tests: `npm run lint:check`, `npm run test:ci`
- Update smoke/CI scripts if your new feature must be covered.
