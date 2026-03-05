# Quickstart: Backend in Docker (No Host Access)

**Feature**: 003-containerize-backend  
**Date**: 2026-03-05

Run the backend as a container. The backend is **only** reachable from within Docker (e.g. by other containers on the same network). It is **not** exposed to the host.

## Prerequisites

- Docker and Docker Compose (Compose V2: `docker compose`). The backend container workflow requires Docker and Compose; there is no fallback for running the backend without them.
- No need to install Go (or any backend runtime) on the host; the container image includes the built binary.

## Start the backend (with database)

From the repository root:

```bash
docker compose up -d db
```

Wait for the database to be healthy:

```bash
docker compose ps
```

When `db` shows state **healthy**, start the backend:

```bash
docker compose up -d backend
```

Or start both in one go (backend will wait for db to be healthy):

```bash
docker compose up -d
```

The backend container will start and connect to MongoDB at `mongodb://db:27017`. It is typically ready to accept requests within 30 seconds after start, once the database is healthy. It listens on port 8080 **inside** the container only; no port is published to the host.

## Verify backend is running

- **From the host**: You **cannot** reach the backend on localhost. For example, `curl http://localhost:8080` (or connecting to 127.0.0.1:8080) will fail or be refused. This is intentional (spec: backend available only from within Docker).
- **From another container on the same network**: Run a container on the same Compose network and call `http://backend:8080` (e.g. a quick `curl http://backend:8080/...` from a temporary container, or run the frontend in a container that uses `backend:8080` as the API URL).

Example one-off check from another container (run from repo root after `docker compose up -d`):

```bash
docker compose run --rm curlimages/curl curl -s -o /dev/null -w "%{http_code}" http://backend:8080/
```

A non-zero exit or non-2xx status may mean the backend is not ready yet (wait a few seconds and retry). Success (e.g. 200 or 404) means the backend is listening.

## Stop the backend

```bash
docker compose stop backend
```

Or stop all services (db and backend):

```bash
docker compose down
```

To also remove the database volume (ephemeral):

```bash
docker compose down -v
```

## Connecting other services to the backend

Any service that runs in a container on the same Compose default network should use:

- **Hostname**: `backend`
- **Port**: `8080`
- **Example base URL**: `http://backend:8080`

Do not use `localhost` or the host machine IP from inside another container; use the service name `backend`.

## Edge cases

- **Backend container fails to start**: Run `docker compose logs backend`; check `MONGODB_URI` and that `db` is healthy. See Troubleshooting below.
- **No container runtime on host**: Docker and Docker Compose are required; document this in Prerequisites above. There is no fallback without Docker.
- **Other container cannot reach backend**: Use hostname `backend` and port `8080` on the same Compose network; do not use localhost or host IP from inside containers. See [contracts/backend-runtime.md](contracts/backend-runtime.md) for full error behavior.

## Troubleshooting

| Issue | What to do |
|-------|-------------|
| Backend container exits immediately | Run `docker compose logs backend`. Check MongoDB connection (is `db` healthy?) and `MONGODB_URI` (must be `mongodb://db:27017` in container). |
| Need to reach backend from host | Out of scope for this feature; backend is intentionally not exposed to the host. Run the backend locally (e.g. `go run ./cmd/server`) if you need host access for development. |
| Other container cannot reach backend | Ensure both are on the same Compose network and use `http://backend:8080`, not localhost. |

## Summary

| Goal | Command / Info |
|------|-----------------|
| Start DB + backend | `docker compose up -d` (or `up -d db` then `up -d backend`) |
| Backend reachable from host? | No (by design) |
| Backend reachable from other containers? | Yes, at `http://backend:8080` |
| Stop backend | `docker compose stop backend` or `docker compose down` |
