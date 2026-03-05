# Data Model: Containerize Frontend (Runtime Config)

**Feature**: 004-containerize-frontend  
**Date**: 2026-03-05

This feature does not introduce new application entities. The data model here describes the **frontend runtime configuration** that the implementation must support (for documentation and contract consistency).

## Runtime Entities

### Frontend Service (runtime)

| Attribute | Type | Required | Constraints | Notes |
|-----------|------|----------|-------------|--------|
| Host port | number | yes | 3000 | Outer port; users access http://localhost:3000 |
| Container port | number | yes | 5137 | Inner port; dev server listens on 5137 |
| Mode | enum | yes | development | Dev server only; no production build |
| Lifecycle | — | — | start, stop | Via Compose or container commands |
| API base URL | string | configurable | Default http://backend:8080 when in same composition | Overridable via VITE_API_URL |

No persistent storage of "service" in a database; this is the runtime shape that the container definition and docs must satisfy.

### Configuration (env / Compose)

| Setting (logical) | Source | Required | Notes |
|-------------------|--------|----------|--------|
| API base URL | VITE_API_URL | no (default when in composition) | Default http://backend:8080; override for host-run backend or other envs |
| Host port | Compose ports | yes | 3000 |
| Container port | App/config | yes | 5137 (Vite or dev server config) |

## Lifecycle / State

- **Start**: One documented command (e.g. `docker compose up -d frontend` or `docker compose up -d`) brings the frontend up; typically ready within 60 seconds (spec SC-001).
- **Stop**: One documented command (e.g. `docker compose stop frontend` or `docker compose down`) stops the frontend.
- **Reachability**: From host at http://localhost:3000 (port 3000). From other containers, frontend can be reached at `frontend:5137` on the same network if needed.

## Relationships

- Frontend calls the **backend** API using the configured API base URL (default http://backend:8080 when both are in the same Compose network). No new application data entities; existing UI and API consumption unchanged.
