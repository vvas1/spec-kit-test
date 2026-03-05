# Quickstart: Frontend in Docker (Host Port 3000, Container Port 5137)

**Feature**: 004-containerize-frontend  
**Date**: 2026-03-05

Run the frontend as a container. The frontend is reachable from the host at **http://localhost:3000**. The dev server runs inside the container on port **5137**. Development only; not for production deployment.

## Prerequisites

- Docker and Docker Compose (Compose V2: `docker compose`). The frontend container workflow requires Docker and Compose; there is no fallback for running the frontend without them.
- No need to install Node.js or the frontend toolchain (npm) on the host; the container runs the dev server.

## Start the frontend (with backend and database)

From the repository root:

**Option A – Full stack (db + backend + frontend):**

```bash
docker compose up -d
```

**Option B – Start services in order:**

```bash
docker compose up -d db
# Wait for db healthy: docker compose ps
docker compose up -d backend
docker compose up -d frontend
```

The frontend container will start and, when in the same composition as the backend, use **http://backend:8080** as the default API base URL (via `VITE_API_URL`). The app is typically ready within 60 seconds after start. Access it at **http://localhost:3000** (host port 3000 maps to container port 5137).

## Verify frontend is running

- Open **http://localhost:3000** in a browser. The frontend application should load. (Host port 3000 maps to container port 5137.)
- If the backend is running, use the app (e.g. load issues, users) to confirm API calls succeed. If the backend is not running, the UI may load but API calls will fail until the backend is up.

## Stop the frontend

```bash
docker compose stop frontend
```

Or stop all services:

```bash
docker compose down
```

To also remove the database volume:

```bash
docker compose down -v
```

## Configuring the API base URL

When the frontend runs in the same composition as the backend, the default API base URL is **http://backend:8080** (set via `VITE_API_URL` in Compose). Override by setting the `VITE_API_URL` environment variable for the frontend service. To override (e.g. backend running on the host):

- Set the environment variable `VITE_API_URL` for the frontend service (e.g. `http://host.docker.internal:8080` on Docker Desktop, or the host’s IP and port). See [contracts/frontend-runtime.md](contracts/frontend-runtime.md) for details.

## Edge cases

- **Frontend container fails to start** (e.g. build failure or port 3000 already in use): Run `docker compose logs frontend` to see errors. Free port 3000 on the host or fix the build.
- **No container runtime on host**: Docker and Docker Compose are required; see Prerequisites. There is no fallback without Docker.
- **API calls fail**: Ensure the backend is running and on the same Compose network. Default `VITE_API_URL` is http://backend:8080. If the backend runs on the host, set `VITE_API_URL` to the host-accessible URL (e.g. http://host.docker.internal:8080). See [contracts/frontend-runtime.md](contracts/frontend-runtime.md).

## Troubleshooting

| Issue | What to do |
|-------|------------|
| Frontend container exits or port 3000 in use | Run `docker compose logs frontend`. Free port 3000 on the host or change the host port in `docker-compose.yml`. |
| API calls fail (e.g. network error) | Ensure the backend is running and on the same Compose network. Default URL is http://backend:8080. If the backend runs on the host, set `VITE_API_URL` to the host-accessible URL (e.g. http://host.docker.internal:8080). |
| No container runtime on host | Docker and Docker Compose are required; see Prerequisites. |

## Summary

| Goal | Command / Info |
|------|-----------------|
| Start full stack | `docker compose up -d` |
| Access frontend | http://localhost:3000 |
| Host port → container port | 3000 → 5137 |
| Default API URL (with backend in same compose) | http://backend:8080 (override via VITE_API_URL) |
| Stop frontend | `docker compose stop frontend` or `docker compose down` |
