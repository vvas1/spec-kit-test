# Data Model: Containerize Backend (Runtime Config)

**Feature**: 003-containerize-backend  
**Date**: 2026-03-05

This feature does not introduce new application entities. The data model here describes the **backend runtime configuration** that the implementation must support (for documentation and contract consistency).

## Runtime Entities

### Backend Service (runtime)

| Attribute | Type | Required | Constraints | Notes |
|------------|------|----------|-------------|--------|
| Listen address | string | yes | Inside container only (e.g. :8080) | Not published to host |
| Host exposure | enum | yes | none | No ports mapping in Compose |
| Network | string | yes | Same Compose network as db (and other services) | Discovery by service name `backend` |
| Lifecycle | — | — | start, stop | Via Compose or container commands |

No persistent storage of "service" in a database; this is the runtime shape that the container definition and docs must satisfy.

### Configuration (env / Compose)

| Setting (logical) | Source | Required | Notes |
|-------------------|--------|----------|--------|
| MongoDB connection | MONGODB_URI | yes | In container set to `mongodb://db:27017` (service name) |
| Listen port | Application default | yes | 8080 inside container; not published |
| Service name | Compose service name | yes | `backend` for discovery on network |

## Lifecycle / State

- **Start**: One documented command (e.g. `docker compose up -d backend` or `docker compose up -d` for db + backend) brings the backend up after DB is healthy; readiness within 30 seconds (spec SC-001).
- **Stop**: One documented command (e.g. `docker compose stop backend` or `docker compose down`) stops the backend.
- **Reachability**: From host: not reachable (spec FR-002). From other containers on same network: reachable at `backend:8080`.

## Relationships

- Backend connects to **MongoDB** via MONGODB_URI (in container: `mongodb://db:27017`). No new application data entities; existing Issue/User data and API behavior unchanged.
- Other services (e.g. frontend container, future microservices) connect to the backend via **backend** hostname and port **8080** on the same Compose network.
