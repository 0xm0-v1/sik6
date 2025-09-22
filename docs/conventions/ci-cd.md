# CI/CD

This document describes the CI/CD bootstrap for the **sik6** project: objectives, triggers, job stages, workflow structure, and policies.

## Workflow file

Path: `.github/workflows/ci-bootstrap.yml`

Workflow name: **CI Bootstrap**

## What it does today (bootstrap skeleton)

- Triggers on:
  - Pull requests targeting `main`
  - Pushes to `main`
  - Manual dispatch (`workflow_dispatch`)

## Local fixes / CI verifies policy

- No auto-fix in CI. Formatters / linters will run only in verify mode.
- Developers must fix issues locally; CI verifies.

## Planned stages (future stories)

| Stage         | Purpose                                          |
| ------------- | ------------------------------------------------ |
| Static Checks | format / lint / static analysis                  |
| Build         | compile / build artifacts                        |
| Test          | unit, integration, coverage, reports             |
| Deploy        | staging / sandbox / registry / environment aware |

_This CI / CD system ensures consistent quality and visible verification of every change before it reaches `main`._
