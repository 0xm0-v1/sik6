# Commit Conventions

This document defines the commit conventions to follow for the **sik6** project.  
The goal is to ensure **clear traceability**, **consistency** in Git history, and to **enable automation** (changelog, release, CI/CD).

## Format

Each commit must follow this structure:\
`[SIK6-<num>] <type>(<scope>): <short, imperative message> #<ref>`

### Details:

- **[SIK6-<num>]**: project ticket ID (`SIK6` is constant, `<num>` is the ticket number).
- **type**: the nature of the change.
- **scope** _(optional)_: the module, file, or component affected.
- **message**: a short description written in the imperative mood.
- **ref**: short ID reference of the issue.

## Rules

1. **Ticket ID required**
   - Always use the **Ticket number** (`SIK6-123`), **never a User Story ID**.
   - Commits trace back to a **Ticket**, and the Ticket itself is linked to its User Story.

2. **Conventional commit types**
   - `feat`: a new feature
   - `fix`: a bug fix
   - `docs`: documentation only changes
   - `style`: formatting, missing semi colons, etc.; no code change
   - `refactor`: code change that neither fixes a bug nor adds a feature
   - `perf`: performance improvements
   - `test`: adding or updating tests
   - `chore`: build process or auxiliary tool changes
   - `design`: design related

3. **Example of Scope**
   #### Backend

- `api`: REST/GraphQL endpoints, controllers
- `auth`: authentication, authorization, session handling
- `db`: database models, schema, migrations
- `core`: core domain logic, services, business rules
- `cache`: caching layer (Redis, in-memory, etc.)
- `config`: configuration files, environment setup

#### Frontend

- `ui`: user interface, layout, components
- `ux`: user experience, interactions, accessibility
- `forms`: form handling, validation
- `routing`: client-side routing, navigation

#### Infrastructure & Ops

- `ci`: CI/CD pipelines (GitHub Actions, GitLab CI, etc.)
- `build`: build system, bundlers, compilers
- `deps`: dependencies, libraries, package updates
- `docker`: Dockerfiles, containerization
- `infra`: infrastructure as code, deployment scripts

#### Quality & Testing

- `test`: unit/integration tests
- `e2e`: end-to-end tests
- `lint`: linting, formatting, code quality

#### Documentation

- `docs`: documentation, README, ADRs
- `changelog`: versioning, release notes

#### Security

- `sec`: security fixes, vulnerabilities, patches

4. **Message**
   - Keep it short and descriptive.
   - Use the imperative mood: _"add"_, _"fix"_, _"update"_ instead of _"added"_, _"fixed"_, _"updated"_.

## Hooks Enforcement

To ensure commit conventions are respected consistently, Git hooks are configured with **Husky**.

### Prepare-commit-msg

- Automatically injects the correct **ticket prefix** (`[SIK6-<num>]`) into the commit message template.
- This ensures contributors don’t forget to include the ticket ID.
- Developers can then focus on writing the commit `type(scope): message #ref`.

### Commit-msg

- Validates the final commit message format against the defined conventions.
- If the message does not match the required pattern, the commit is **blocked** until corrected.
- This guarantees that all commits are tied to a valid ticket and follow the structure:

### Usage

```bash
git commit
```

## Why this matters

- Ensures every commit is tied to a **Ticket**.
- Tickets are then connected to **User Stories**, **PRs**, and **Issues**, providing full traceability.
- Keeps commit history clean and consistent.

_This commits system ensures every commit is tied to a **Ticket** and are then connected to **User Stories**, **PRs**, and **Issues**, providing full traceability, making commit history clean and consistent._
