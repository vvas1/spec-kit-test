# Research: Containerize Database

**Feature**: 002-containerize-database  
**Date**: 2026-03-04

## Container Runtime

- **Decision**: Docker with Docker Compose for defining and running the MongoDB service.
- **Rationale**: Spec requires a single, documented way that works on "standard development and CI environments." Docker is the de facto standard for dev/CI containerization; Compose gives one-file definition, default project name, and simple `up`/`down` lifecycle. No need for Kubernetes or other orchestrators (out of scope).
- **Alternatives considered**: Podman (good alternative but less universal in CI); raw `docker run` (Compose is easier to document and version); no container (rejected—spec explicitly asks to containerize).

## MongoDB Image and Version

- **Decision**: Use the official MongoDB image (e.g. `mongo:7`) with version pinned in Compose to match or align with backend driver compatibility.
- **Rationale**: Constitution specifies MongoDB; official image is maintained and well-documented. Pinning avoids surprise upgrades; 7.x is widely supported by the Go driver.
- **Alternatives considered**: Custom image (unnecessary for dev/CI); latest tag (avoid for reproducibility).

## Readiness / Health Check

- **Decision**: Document a readiness check that scripts and CI can use: (1) Docker Compose `healthcheck` using `mongosh` or `mongosh --eval "db.adminCommand('ping')"` inside the container, and (2) a simple TCP connection to the exposed port (e.g. `localhost:27017`) or a small script that runs `mongosh` from the host if available. Prefer Compose healthcheck so `docker compose up -d` can be followed by polling `docker compose ps` or waiting for healthy status.
- **Rationale**: Spec FR-007 requires a documented way to verify the database is ready; avoids fixed sleeps and flaky CI.
- **Alternatives considered**: Sleep-only (rejected per spec); custom HTTP health endpoint (MongoDB doesn’t expose HTTP by default; TCP or mongosh is standard).

## Binding and Port

- **Decision**: Bind MongoDB to `127.0.0.1:27017` by default in the Compose service. Document how to override (e.g. change port or bind to `0.0.0.0`) for remote access if needed.
- **Rationale**: Spec FR-008 requires localhost by default; 27017 is the default MongoDB port and matches existing backend default `MONGODB_URI`.
- **Alternatives considered**: Bind to all interfaces by default (rejected per spec); different default port (would force every app to set URI).

## Credentials

- **Decision**: Support configurable credentials via environment variables (e.g. `MONGO_INITDB_ROOT_USERNAME`, `MONGO_INITDB_ROOT_PASSWORD` in Compose; application uses same in `MONGODB_URI`). Document in quickstart and contract. If no env is set, run without auth for local dev (document security caveat).
- **Rationale**: Spec FR-009 allows "environment variables or a documented alternative"; env vars are portable and CI-friendly.
- **Alternatives considered**: File-based secrets (acceptable as documented alternative); no auth only (simplest for dev; document that production must not use this setup).

## Persistence vs Ephemeral

- **Decision**: (1) **Persistent**: Compose service uses a named volume (or host path) for MongoDB data dir; `docker compose down` keeps data unless `-v` is used. (2) **Ephemeral**: Document a profile or override (e.g. no volume, or `docker compose down -v`) so tests get a clean state.
- **Rationale**: Spec FR-004 requires both modes; FR-010 requires distinct instance or clean state per test run.
- **Alternatives considered**: Always ephemeral (would break "persist across restarts" for dev); always persistent (would break "clean state" for tests).

## Isolation for Parallel / Back-to-Back Test Runs

- **Decision**: For CI or parallel runs, use one of: (1) distinct Compose project name per run (e.g. `COMPOSE_PROJECT_NAME=issue-tracker-${CI_JOB_ID}` or random suffix), so each run gets its own container and optional volume; or (2) single project but run `docker compose down -v` before `up` so the next run gets a fresh container and no volume. Document both: "single run" (deterministic port 27017) vs "parallel/isolated run" (project name or down/up with volume remove).
- **Rationale**: Spec FR-010 requires each run to get a distinct instance or clean state with no shared state; Compose project name isolates containers and optionally volumes.
- **Alternatives considered**: Single shared instance with serialized tests (rejected per spec); dynamic port per run (adds complexity; project name is simpler).

## Summary Table

| Topic           | Decision                                              |
|----------------|--------------------------------------------------------|
| Runtime        | Docker + Docker Compose                               |
| Image          | Official MongoDB (e.g. mongo:7), version pinned        |
| Readiness      | Compose healthcheck (mongosh ping) + document for CI  |
| Bind           | 127.0.0.1:27017 by default; document override         |
| Credentials    | Env vars (MONGO_*, MONGODB_URI); document no-auth dev  |
| Persistence    | Named volume for dev; profile or down -v for ephemeral |
| Test isolation | Compose project name per run or down -v before up      |
