# Contract: Backend Runtime (Containerized, Docker-Only Access)

**Feature**: 003-containerize-backend  
**Date**: 2026-03-05

This contract defines how the containerized backend is started, stopped, and how other services connect to it. The backend is **not** exposed to the host; it is reachable only from within the Docker (Compose) environment.

## Commands (operational)

The implementation MUST provide a documented way to perform the following. Exact commands are implementation-defined (e.g. Docker Compose); this contract specifies behavior.

| Operation | Input (env/config) | Output / Behavior |
|-----------|---------------------|-------------------|
| Start | Optional: Compose project | Backend process running in a container; **no** port published to host. Reachable from other containers at `backend:8080`. |
| Stop | None (or project name) | Backend container stopped; no longer reachable. |
| Readiness | None | Backend ready to accept requests within 30 seconds after container start (once DB is healthy). |

## Environment Variables (backend container)

| Purpose | Name | Value in container | Notes |
|---------|------|--------------------|--------|
| MongoDB connection | `MONGODB_URI` | `mongodb://db:27017` | Required; use Compose service name `db`. Overridable for different DB (e.g. auth) if documented. |

## Network and Reachability

| Requirement | Behavior |
|-------------|----------|
| Host access | Backend MUST NOT be reachable from the host when using the intended configuration (no `ports:` for backend service). |
| Same-network access | Backend MUST be reachable from other containers on the same Compose network at hostname **backend**, port **8080** (e.g. `http://backend:8080` for HTTP API). |
| Startup order | Backend SHOULD start only after MongoDB is healthy (e.g. `depends_on: db` with `condition: service_healthy`). |

## Readiness and Verification

- **From another container**: A container on the same network can verify backend is up by issuing an HTTP request to `http://backend:8080/` (or a documented health/API path). Success (e.g. 200 or 404 for unknown path) indicates the backend is listening.
- **From host**: Connection to the backend MUST fail (e.g. no port bound to host). Do not publish backend port to the host.
- **Timing**: Backend is typically ready within 30 seconds after the container starts, assuming DB is already healthy. For exact start/stop commands and verification steps, see [quickstart.md](../quickstart.md).

## Defaults and Overrides

| Item | Default | Override |
|------|---------|----------|
| Backend port (inside container) | 8080 | Change only if application is configured to listen on another port. |
| MONGODB_URI | `mongodb://db:27017` | Set in Compose or env for auth (e.g. `mongodb://user:pass@db:27017`) or different host. |
| Host exposure | None | Do not add `ports:` for backend to comply with spec. |

## Error Behavior (spec edge cases)

- **Backend container fails to start**: Compose will show failed state; use `docker compose logs backend` to see application errors (e.g. MongoDB connection refused). Resolve by ensuring `db` is healthy and MONGODB_URI is correct.
- **No container runtime on host**: The backend container workflow requires Docker (and Docker Compose). Document this in quickstart; no fallback for running without Docker in scope.
- **Other service cannot reach backend**: Ensure the client container is on the same Compose network and uses hostname `backend` and port `8080` (not localhost or host IP).
