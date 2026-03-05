# Research: Containerize Backend (Docker-Only Access)

**Feature**: 003-containerize-backend  
**Date**: 2026-03-05

## Backend Container Image (Go)

- **Decision**: Use a multi-stage Dockerfile in `backend/`: (1) build stage with Go image to compile the binary; (2) run stage with a minimal image (e.g. `alpine` or `scratch` + ca-certificates if needed) to run the binary only. Pin Go version to match `go.mod` (e.g. 1.23).
- **Rationale**: Keeps image size small and avoids shipping source or build tools in the runtime image; aligns with common Go container best practices. Alpine is a good balance of size and compatibility (e.g. for TLS); scratch is smaller but may require static binary and extra care for certs.
- **Alternatives considered**: Single-stage Dockerfile (larger image); distroless (good alternative; Alpine chosen for simplicity and wide use).

## No Host Exposure (Spec FR-002)

- **Decision**: In docker-compose.yml, define the `backend` service with **no** `ports:` mapping. The backend listens on `:8080` inside the container; that port is not published to the host. Other containers on the same Compose network reach the backend at hostname `backend` and port `8080`.
- **Rationale**: Spec explicitly requires the backend to be available only from within Docker; omitting `ports:` is the standard way to keep a service internal in Compose.
- **Alternatives considered**: Publishing to 127.0.0.1 (rejected—still exposes to host); using `network_mode: host` (rejected—would expose to host and break isolation).

## Service Discovery (Spec FR-003)

- **Decision**: Use the default Compose network for the project. Both `db` and `backend` are on this network. Backend is reachable as `backend` (service name) on port `8080`. Backend connects to MongoDB using `MONGODB_URI=mongodb://db:27017` (service name `db`). Document in contract and quickstart that any other container on the same network should use `http://backend:8080` (or the appropriate API base URL).
- **Rationale**: Compose assigns a DNS name per service name on the default network; no extra configuration needed. Matches spec requirement for discoverability by name/alias.
- **Alternatives considered**: Custom network (optional; default is sufficient); explicit links (deprecated; Compose network DNS is preferred).

## Backend Startup and Readiness (Spec FR-004, SC-001)

- **Decision**: (1) Use `depends_on: db` with `condition: service_healthy` so the backend starts only after MongoDB is healthy. (2) Optionally add a Compose healthcheck for the backend (e.g. `curl -f http://localhost:8080/...` or a dedicated health endpoint if one exists); if the backend has no health endpoint, document that "ready" means container is running and the app listens on 8080. (3) Document that typical readiness is within 30 seconds after container start once DB is healthy.
- **Rationale**: Reduces race conditions where backend starts before DB is ready; 30-second target is achievable for a Go server connecting to a local MongoDB container.
- **Alternatives considered**: Fixed sleep (rejected—unreliable); no depends_on (rejected—would allow backend to start before DB and fail or retry unnecessarily).

## Docker Compose Integration with Existing `db` Service

- **Decision**: Add a `backend` service to the existing root `docker-compose.yml`. Backend service: build context `./backend`, Dockerfile `backend/Dockerfile`; no `ports:`; environment `MONGODB_URI=mongodb://db:27017`; `depends_on: db: condition: service_healthy`; same default network as `db`. Optionally set restart policy (e.g. `restart: "no"` for dev) and ensure backend process runs as main process (no shell wrapper unless needed).
- **Rationale**: Single Compose file keeps one place to start DB and backend; reuses 002-containerize-database setup; backend and db on same network by default.
- **Alternatives considered**: Separate compose file for backend (possible but adds complexity); running backend in same container as DB (rejected—separate services per spec).

## Failure and Edge Cases (Spec Edge Cases)

- **Decision**: (1) Container fails to start: rely on Compose and container logs; document that users should run `docker compose logs backend` and fix config (e.g. wrong MONGODB_URI or DB not reachable). (2) No container runtime on host: document in quickstart that Docker (and Compose) are required to run the backend container. (3) Discovery: document that other services use hostname `backend` and port `8080`; no hard-coded host IPs.
- **Rationale**: Spec asks for clear failure feedback and documentation for container runtime requirement and discovery; no new tooling beyond Compose and logs is required.

## Summary Table

| Topic | Decision |
|-------|----------|
| Image | Multi-stage Dockerfile in backend/; Go build stage + minimal run image (e.g. Alpine) |
| Host exposure | No `ports:` for backend service; internal only |
| Discovery | Default Compose network; backend at `backend:8080`, DB at `db:27017` |
| MONGODB_URI in container | `mongodb://db:27017` |
| Startup order | `depends_on: db` with `condition: service_healthy` |
| Readiness | Optional healthcheck; document 30s target and how to verify |
| Failure/docs | Compose logs; quickstart states Docker required and how others connect |
