# Quickstart: Containerized Database (002-containerize-database)

**Feature**: Run MongoDB in a container for local development and CI  
**Stack**: Docker + Docker Compose; backend (Go) and frontend (React + TypeScript + MUI) unchanged

## Prerequisites

- Docker and Docker Compose (Compose V2: `docker compose`)
- Go and Node.js (for backend and frontend; see 001-issue-crud quickstart if needed)

## Start the database

From the **repository root**:

```bash
docker compose up -d db
```

Wait until the database is ready (see [Readiness check](#readiness-check) below). Default: MongoDB is available on port 27017 on all interfaces (reachable from any host) with no auth (dev only; do not expose on untrusted networks).

## Readiness check

Before starting the backend or running tests, verify the database accepts connections:

- **Option 1 (recommended)**: Check Compose health status (matches `docker-compose.yml` healthcheck):
  ```bash
  docker compose ps
  ```
  Wait until the `db` service shows state **healthy**. The root `docker-compose.yml` defines a healthcheck using `mongosh --eval "db.adminCommand('ping')"`.

- **Option 2**: If `mongosh` is installed on the host:
  ```bash
  mongosh --eval "db.adminCommand('ping')" mongodb://localhost:27017
  ```
  Exit code 0 means ready.

- **Option 3**: Use any documented script or command from the implementation (see [contracts/database-runtime.md](./contracts/database-runtime.md)).

## Run the application

1. Start the database (see above) and wait for readiness.
2. Start the backend:
   ```bash
   cd backend && go run ./cmd/server
   ```
   Default `MONGODB_URI` is `mongodb://localhost:27017`; override with env if needed.
3. Start the frontend:
   ```bash
   cd frontend && npm run dev
   ```

## Stop the database

From the repository root:

- **Stop but keep data (persistent)**:
  ```bash
  docker compose down
  ```
- **Stop and remove data (ephemeral)**:
  ```bash
  docker compose down -v
  ```

## Persistent vs ephemeral mode

- **Persistent (default)**: The Compose file uses a named volume `mongodb_data` for MongoDB data. Use `docker compose down` to stop; data is kept. The next `docker compose up -d db` reuses the same data.
- **Ephemeral**: Use `docker compose down -v` to stop and remove the volume. The next `docker compose up -d db` starts with a clean state (no prior data). Use this for tests or when you want a fresh DB. To verify: after `down -v`, run `up -d db`, connect with the app or mongosh, and confirm no previous data is present.

## Connection parameters

- **Default**: The backend uses `MONGODB_URI`; if unset, it defaults to `mongodb://localhost:27017` (no auth). The Compose file publishes port 27017 on all interfaces, so the DB is reachable from any host (same machine: localhost; others: use the host’s IP). Use only in trusted dev/test networks.
- **Credentials**: To use auth, set `MONGO_INITDB_ROOT_USERNAME` and `MONGO_INITDB_ROOT_PASSWORD` (e.g. in `.env`) and include them in `MONGODB_URI` (e.g. `mongodb://user:pass@localhost:27017`). Unset = no auth (dev only; do not use in production).
- **Localhost-only**: To bind only to the host, edit `docker-compose.yml` and set ports to `"127.0.0.1:27017:27017"`.

## Configuration (summary)

| Need              | How |
|-------------------|-----|
| Different port     | Override in Compose or set in `MONGODB_URI` (see contract). |
| Credentials        | Set `MONGO_INITDB_ROOT_USERNAME` / `MONGO_INITDB_ROOT_PASSWORD` and use same in `MONGODB_URI` (see contract). |
| Localhost-only     | Edit Compose: use `127.0.0.1:27017:27017` (default is all interfaces). |
| CI / parallel runs | Use unique `COMPOSE_PROJECT_NAME` per run or `docker compose down -v` before next `up` (see contract). |

## Tests

- **Backend (with containerized DB)**: From repo root, run `./scripts/test-with-db.sh`. This starts the DB, waits for healthy, runs `go test ./...` in `backend/`, then runs `docker compose down -v` so the next run gets a clean state.
- **Isolated/parallel runs**: Set `COMPOSE_PROJECT_NAME` to a unique value per run (e.g. `COMPOSE_PROJECT_NAME=issue-tracker-${CI_JOB_ID}` or a random suffix) so each run gets its own container and volume; or run `docker compose down -v` before the next `docker compose up -d db` for a clean state. See [contracts/database-runtime.md](./contracts/database-runtime.md).
- **Frontend**: Start DB and backend, then run frontend tests as in 001-issue-crud quickstart.

## Key files

- **Spec**: [spec.md](./spec.md)
- **Plan**: [plan.md](./plan.md)
- **Data model (runtime)**: [data-model.md](./data-model.md)
- **Database runtime contract**: [contracts/database-runtime.md](./contracts/database-runtime.md)
