# Contract: Database Runtime (Containerized MongoDB)

**Feature**: 002-containerize-database  
**Date**: 2026-03-04

This contract defines how the containerized database is started, stopped, configured, and verified ready. It is the interface between developers/CI and the database runtime.

## Commands (operational)

The implementation MUST provide a single, documented way to perform the following. Exact commands are implementation-defined (e.g. Docker Compose); this contract specifies behavior and inputs/outputs.

| Operation   | Input (env/config)     | Output / Behavior |
|------------|------------------------|-------------------|
| Start      | Optional: port, volume, credentials env | Database process running; bound to localhost by default (spec FR-008). |
| Stop       | None (or project name) | Database process stopped; resources released. |
| Remove     | Optional (e.g. `-v` flag) | Containers and optionally volumes removed (ephemeral). |
| Readiness  | None                   | Verifiable "ready to accept connections" (spec FR-007). |

## Environment Variables (configurable)

The following are documentable and overridable (spec FR-005, FR-009). Implemented names:

| Purpose           | Name(s)                     | Default / Notes |
|-------------------|-----------------------------|------------------|
| Connection string | `MONGODB_URI`               | `mongodb://localhost:27017` (no auth for local dev; do not use in production). |
| DB credentials    | `MONGO_INITDB_ROOT_USERNAME`, `MONGO_INITDB_ROOT_PASSWORD` | Optional; if set, include in `MONGODB_URI`. |
| Port / bind       | Compose `ports` (e.g. `27017:27017`) | 27017 on all interfaces (any host). For localhost-only use `127.0.0.1:27017:27017` in `docker-compose.yml`. |
| Project isolation | `COMPOSE_PROJECT_NAME`      | For CI/parallel: set to a unique value per run (e.g. `issue-tracker-${CI_JOB_ID}`) so each run gets its own container and volume. |

## Readiness Check (FR-007)

The implementation MUST document how to verify the database is ready, e.g.:

- **Option A**: Docker Compose healthcheck (e.g. `mongosh --eval "db.adminCommand('ping')"`); consumers wait until container is healthy (e.g. `docker compose ps` shows healthy, or poll).
- **Option B**: Script or command that attempts connection (e.g. `mongosh` to connection string, or TCP connect to port) and exits 0 when ready, non-zero otherwise.

Documentation MUST state the exact command or steps so scripts and CI can wait reliably (no fixed sleep as the only option).

## Defaults and Overrides

| Requirement   | Default behavior           | Override (documented) |
|---------------|----------------------------|------------------------|
| Bind address  | 0.0.0.0 (all interfaces; any host) | For localhost-only use `127.0.0.1:27017:27017` in Compose. |
| Port          | 27017                       | Document how to set different port (Compose or env). |
| Persistence   | Persistent (volume)         | Document how to run ephemeral (e.g. no volume, or `down -v`). |
| Credentials   | No auth (local dev)         | Document env vars for username/password and URI. |

## Isolation (FR-010)

For multiple test runs (parallel or back-to-back), the implementation supports both strategies:

- **Default (single run)**: One Compose project; port 27017. For a clean state between runs, run `docker compose down -v` before the next `docker compose up -d db`.
- **Isolated runs**: Set `COMPOSE_PROJECT_NAME` to a unique value per run (e.g. `COMPOSE_PROJECT_NAME=issue-tracker-${CI_JOB_ID}` or `issue-tracker-$$`). Each run gets its own container and volume; no shared state.

The script `scripts/test-with-db.sh` uses the single project and `docker compose down -v` after tests so the next invocation gets a fresh DB.

## Error Behavior (spec edge cases)

- **Port in use / start failure**: Docker Compose reports bind errors when port 27017 is already in use (e.g. "port is already allocated" or "address already in use"). Resolve by stopping the other process using 27017 or changing the host port in `docker-compose.yml`.
- **Resource exhaustion**: The MongoDB process may exit or fail with an identifiable error if the host runs out of disk or memory; check container logs with `docker compose logs db`.
- **Second start (same port)**: Starting a second Compose project (or a second container on the same host port) while the first is running will fail with a port conflict message. Stop the existing container with `docker compose down` before starting again, or use a different `COMPOSE_PROJECT_NAME` and port for an isolated second instance.
