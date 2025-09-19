# Linting Conventions – sik6

This document defines the linting rules applied in the **sik6** project.  
It covers both the **backend (Go)** and the **frontend (Angular/TypeScript)** stacks.  

Linting ensures **code quality**, **consistency**, and helps **catch bugs early**.  

## Backend (Go)

### Tool
- [golangci-lint](https://golangci-lint.run/) — industry-standard meta-linter for Go.

### Configuration
- File: `.golangci.yml` (root of the repo).
- Reference: `.golangci.reference.yml` (documented configuration for team).

### Enabled Linters
- **govet**: detects misuses (printf formatting, struct alignment, etc.).  
- **staticcheck**: advanced static analysis, best practices, bug prevention.  
- **unused**: finds unused variables, functions, imports.

### Exclusions
- Skip directories: `vendor/`, `node_modules/`, `dist/`, `tmp/`.

### Usage
```bash
make lint
```
- Runs `golangci-lint run ./...` on the whole codebase.
- Test files (`*_test.go`) are also included in linting.

## Frontend (Angular/TypeScript)

### Tool
- [ESLint](https://eslint.org/) with [@angular-eslint](https://github.com/angular-eslint/angular-eslint)

### Configuration
- File: `.eslintrc.json` (project root or `frontend/` when created).
- Reference: `.eslintrc.reference.js` (documented configuration for team).
- Ignore file: `.eslintignore`.

### Enabled Rules
- Base: `eslint:recommended` + `plugin:@angular-eslint/recommended`.
- Custom additions:
    - `no-unused-vars: error` — unused variables are not allowed.
    - `no-console: warn` — discourage leaving `console.log` in commits.

### Exclusions
- Ignored paths: `node_modules/`, `dist/`, `coverage/`, `.tmp/`.

### Usage 
```bash
+ npm run lint
```

```bash
npm run lint:fix
```
- Runs ESLint on all `.ts` files (and Angular templates).

## Pre-commit Hooks
- Husky automatically runs both linters before commits:
    - `make lint` (Go)
    - `npm run lint` (Angular/TypeScript)
> If linting fails, the commit is blocked until issues are resolved.

---

## Runtime Notes

- **Go Linter**
  - Using [`golangci-lint`](https://golangci-lint.run/) **v2**.
  - Configuration file: `.golangci.yml` (with reference `.golangci.reference.yml`).
  - Pre-commit hook blocks commits if any lint errors are found.

- **TypeScript / Angular Linter**
  - Using [ESLint](https://eslint.org/) **v8** (pinned) + `@angular-eslint`.
  - Requires `@typescript-eslint/parser` and `@typescript-eslint/eslint-plugin`.
  - Configuration file: `.eslintrc.json` (with reference `.eslintrc.reference.js`).
  - Runs with option `--no-error-on-unmatched-pattern` to avoid false failures when no `.ts` files exist yet.
  - `.eslintignore` excludes `node_modules/`, `dist/`, `coverage/`, `.tmp/`, and `backend/`.

### Note
- Formatting (Prettier, gofmt, etc.) is handled in a **separate formatter setup ticket**.
- This document will be updated as more rules are introduced over time.

_This linting system ensures that all code remain consistent, maintainable and aligned, while providing immediate feedback to developers through automated checks & fix._
