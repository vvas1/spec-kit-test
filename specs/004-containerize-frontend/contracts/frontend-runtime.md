# Contract: Frontend Runtime (Containerized, Host Port 3000, Container Port 5137)

**Feature**: 004-containerize-frontend  
**Date**: 2026-03-05

This contract defines how the containerized frontend is started, stopped, and how the API base URL is configured. The frontend is reachable from the host at port **3000**; the dev server inside the container listens on port **5137**. Development only; no production deployment.

## Commands (operational)

The implementation MUST provide a documented way to perform the following. Exact commands are implementation-defined (e.g. Docker Compose); this contract specifies behavior.

| Operation | Input (env/config) | Output / Behavior |
|-----------|---------------------|-------------------|
| Start | Optional: Compose project | Frontend dev server running in container; host port 3000 maps to container port 5137. |
| Stop | None (or project name) | Frontend container stopped; port 3000 no longer serves the app. |
| Readiness | None | Frontend typically ready within 60 seconds after container start (spec SC-001). |

## Environment Variables (frontend container)

| Purpose | Name | Value in container (default when in same composition as backend) | Notes |
|---------|------|-----------------------------------------------------------------|--------|
| API base URL | VITE_API_URL | http://backend:8080 | Required for frontend to call backend. Overridable (e.g. http://host.docker.internal:8080 if backend runs on host). Document in quickstart. |

## Port Mapping

| Requirement | Behavior |
|-------------|----------|
| Host port | Frontend MUST be reachable from the host at port **3000** (e.g. http://localhost:3000). |
| Container port | The dev server inside the container MUST listen on port **5137** (or Compose MUST map host 3000 to container 5137). |

## Readiness and Verification

- **From host**: Open http://localhost:3000 in a browser; the frontend application loads. If the backend is not yet running, API calls may fail until the backend is up.
- **Timing**: Frontend is typically ready within 60 seconds after the container starts. For exact start/stop commands and verification steps, see [quickstart.md](../quickstart.md).

## Defaults and Overrides

| Item | Default | Override |
|------|---------|----------|
| Host port | 3000 | Change in Compose `ports` if needed (spec requires 3000). |
| Container port | 5137 | Must match dev server config (e.g. Vite `server.port` or `--port 5137`). |
| VITE_API_URL | http://backend:8080 (when in same composition) | Set in Compose or env for different backend URL. |

## Error Behavior (spec edge cases)

- **Frontend container fails to start**: Compose will show failed state; use `docker compose logs frontend` to see errors (e.g. port 3000 already in use, build failure). Resolve by freeing the port or fixing the build.
- **No container runtime on host**: The frontend container workflow requires Docker (and Docker Compose). Document in quickstart; no fallback without Docker in scope.
- **API calls fail (e.g. network error)**: Ensure backend is running and reachable from the frontend container (same network; default http://backend:8080). If backend runs on host, set VITE_API_URL to http://host.docker.internal:8080 (or equivalent) and document.
