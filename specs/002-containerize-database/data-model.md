# Data Model: Containerize Database (Runtime Config)

**Feature**: 002-containerize-database  
**Date**: 2026-03-04

This feature does not introduce new application entities. The data model here describes the **database runtime configuration** that the implementation must support (for documentation and contract consistency).

## Runtime Entities

### Database Instance (runtime)

| Attribute       | Type   | Required | Constraints                    | Notes                                  |
|----------------|--------|----------|--------------------------------|----------------------------------------|
| Bind address   | string | yes      | Default 127.0.0.1              | Overridable for remote access          |
| Port           | number | yes      | Default 27017                   | Must be documentable and overridable   |
| Persistence    | enum   | yes      | persistent \| ephemeral        | Volume vs no volume / down -v           |
| Lifecycle      | —      | —        | start, stop, (optional) remove | Single documented path (e.g. Compose) |

No persistent storage of "instance" in a database; this is the runtime shape that the container definition and docs must satisfy.

### Configuration (env / Compose)

| Setting (logical)   | Source              | Required | Notes                                      |
|---------------------|---------------------|----------|--------------------------------------------|
| Connection host     | MONGODB_URI or host | yes      | Default localhost for app connection       |
| Connection port     | MONGODB_URI or port | yes      | Default 27017                              |
| Credentials         | Env vars / URI      | configurable | Documented; optional for local no-auth |
| Data directory      | Volume or path      | for persistent | Named volume or host path              |
| Compose project name| COMPOSE_PROJECT_NAME| for isolation | Per-run for CI / parallel runs          |

## Lifecycle / State

- **Start**: One documented command (e.g. `docker compose up -d db`) brings the database up; readiness is verified by documented check (healthcheck or probe).
- **Stop**: One documented command (e.g. `docker compose down`) stops the database; with optional volume removal for ephemeral.
- **Restart**: Same as start after stop; persistent mode retains data; ephemeral mode gives clean state.

## Relationships

- Backend (and any service that uses the database) connects via **Configuration** (host, port, credentials). No new application entities; existing Issue/User data lives in MongoDB as before.
