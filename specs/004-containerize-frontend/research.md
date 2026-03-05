# Research: Containerize Frontend (Host Port 3000, Container Port 5137)

**Feature**: 004-containerize-frontend  
**Date**: 2026-03-05

## Frontend Container Image (Node + Vite dev server)

- **Decision**: Use a single-stage Dockerfile in `frontend/`: Node.js LTS image (e.g. `node:20-alpine` or `node:22-alpine`), install dependencies with `npm ci`, run the dev server with `npm run dev` (Vite). Configure the dev server to listen on port **5137** inside the container (spec requirement). No production build; development mode only.
- **Rationale**: Spec requires dev server with hot reload (FR-007) and inner port 5137. Single stage is sufficient for dev; multi-stage would be for production builds (out of scope). Alpine keeps image size smaller.
- **Alternatives considered**: Multi-stage with production build (rejected—spec says dev only); different Node version (LTS aligns with typical frontend tooling).

## Port Mapping (Spec FR-002, FR-003)

- **Decision**: In docker-compose.yml, map **host port 3000** to **container port 5137** (`ports: "3000:5137"`). Inside the container, the Vite dev server MUST listen on 5137 (currently frontend uses 5173 in vite.config.ts; override via Vite config or CLI `--port 5137` so the container exposes 5137).
- **Rationale**: Spec explicitly requires outer 3000 and inner 5137. Vite supports `server.port` and `--port`; implementation will set 5137 for the containerized run.
- **Alternatives considered**: Keeping 5173 inside and mapping 3000:5173 (rejected—spec says inner 5137); dynamic port (rejected—spec is explicit).

## API Base URL (Spec FR-006)

- **Decision**: The frontend already uses `VITE_API_URL` (see `frontend/src/services/api.ts`: `import.meta.env.VITE_API_URL || 'http://localhost:8080'`). In docker-compose, set `VITE_API_URL=http://backend:8080` for the `frontend` service when running in the same composition as the backend. Document that this is the default and that users can override it (e.g. for host-run backend via `http://host.docker.internal:8080` or custom URL).
- **Rationale**: Clarification session chose default backend via internal network; existing code already supports env override. No application code change required for the default when in composition.
- **Alternatives considered**: Hardcoding backend URL (rejected); different env var name (rejected—VITE_ prefix required for Vite client-side env).

## Docker Compose Integration

- **Decision**: Add a `frontend` service to the existing root `docker-compose.yml`. Build context `./frontend`, Dockerfile `frontend/Dockerfile`; `ports: "3000:5137"`; environment `VITE_API_URL=http://backend:8080`; same default network as `db` and `backend`. Optionally `depends_on: backend` (without condition) so that when user runs `docker compose up -d`, frontend starts after backend is defined (backend may not expose healthcheck; ordering is best-effort unless backend adds healthcheck). For simplicity, document that if backend is not ready, frontend may load but API calls fail until backend is up.
- **Rationale**: Single Compose file; frontend and backend on same network so `http://backend:8080` resolves. No production build or healthcheck required for dev.
- **Alternatives considered**: Separate compose file (adds complexity); frontend depending on backend with healthcheck (backend currently has no HTTP healthcheck; could add later).

## Hot Reload / Source Mount (Optional)

- **Decision**: For a better dev experience, consider mounting the frontend source directory as a volume so that edits on the host are reflected in the container without rebuilding. This is optional in the initial implementation; document in quickstart if implemented. If not implemented, users must rebuild the image or run frontend on the host for live edits.
- **Rationale**: Spec asks for dev server with hot reload; volume mount is the standard way to get hot reload when the app runs in a container. Implementation plan can make this optional (Phase 1: no volume; Phase 2 or follow-up: add volume).
- **Alternatives considered**: No volume (simpler; rebuild to see changes); required volume (better UX; may have platform-specific behavior on some hosts).

## Startup and Readiness (Spec FR-004, SC-001)

- **Decision**: Document that the frontend is typically ready within 60 seconds after the container starts. No formal healthcheck required for dev; users can open http://localhost:3000 and refresh until the app loads. Optionally add a simple healthcheck (e.g. `curl -f http://localhost:5137` inside the container) in a follow-up if needed for CI.
- **Rationale**: Spec SC-001 allows up to 60 seconds; Vite dev server usually starts in under 30 seconds. Keeping implementation simple for dev-only.
- **Alternatives considered**: Compose healthcheck (optional; not blocking).

## Summary Table

| Topic | Decision |
|-------|----------|
| Image | Node LTS (e.g. node:20-alpine), single stage, npm ci + npm run dev |
| Port (container) | Dev server listens on 5137 (Vite config or --port 5137) |
| Port (host) | 3000:5137 in Compose |
| API URL env | VITE_API_URL=http://backend:8080 by default in Compose; overridable, documented |
| Composition | Same network as backend; optional depends_on backend |
| Hot reload | Optional volume mount for source; document if implemented |
| Readiness | Document ~60s; no mandatory healthcheck for dev |
