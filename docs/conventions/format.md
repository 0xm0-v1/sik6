# Code Formatting Conventions

The **sik6** project enforces automated formatting to keep the codebase consistent and easy to read.

## Backend (Go)

- **Tool**: [gofumpt](https://github.com/mvdan/gofumpt), a stricter version of `gofmt`.
- **Manual command**:

```bash
  make format
```

- **Git hook**: staged `.go` files are automatically formatted on commit.

## Frontend (Angular / TypeScript / HTML)

- **Tool**: [Prettier](https://prettier.io/) with Angular-friendly settings:
  - `singleQuote: true` (TypeScript → single quotes),
  - `parser: angular` (for HTML templates).
- **Manual commands:**

```bash
npm run format       # format all frontend files
npm run format:check # verify formatting without writing
```

- **Git hook**: staged frontend files (`.ts`, `.html`, `.scss`, `.css`, `.json`, `.md`) are automatically formatted using [lint-staged](https://github.com/lint-staged/lint-staged).

## Lint vs Format

- **Prettier** → handles **style** (indentation, quotes, spacing, etc.).
- **ESLint** → handles **quality** (unused variables, `console.log`, Angular/TS rules).
  - All style-related rules are disabled via `eslint-config-prettier`.
