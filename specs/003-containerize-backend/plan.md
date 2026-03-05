# Implementation Plan: Containerize Backend (Docker-Only Access)

**Branch**: `003-containerize-backend` | **Date**: 2026-03-05 | **Spec**: [spec.md](./spec.md)  
**Input**: Feature specification from `/specs/003-containerize-backend/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/plan-template.md` for the execution workflow.

## Summary

Containerize the Go backend so it runs as a service inside Docker and is reachable only from within the Docker environment (e.g., by other containers on the same network). No ports are published to the host. Use a Dockerfile for the backend and extend the existing docker-compose setup with a `backend` service that joins the same network as `db`; backend connects to MongoDB via service name `db` and is discoverable as `backend` on the internal network. Documentation and a backend runtime contract describe how to start/stop and how other services connect.

## Technical Context

**Language/Version**: Go 1.23+ (existing backend); no new application code; container definition and tooling only.  
**Primary Dependencies**: Docker, Docker Compose; existing Go backend (standard library, MongoDB driver).  
**Storage**: MongoDB (unchanged; backend uses MONGODB_URI; in Compose set to `mongodb://db:27017`).  
**Testing**: Existing backend unit/integration tests; add path to run backend in container and verify no host exposure and same-network reachability.  
**Target Platform**: Developer machines and CI (Linux, macOS, Windows with Docker).  
**Project Type**: Infrastructure / developer experience (add backend container to existing web application).  
**Performance Goals**: Backend ready within 30 seconds of container start (spec SC-001).  
**Constraints**: Backend MUST NOT be exposed to the host (no `ports:` for backend service); MUST be reachable only from other containers on the same Compose network.  
**Scale/Scope**: Single backend instance per Compose project; dev and CI only.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Gate | Status |
|-----------|------|--------|
| I. Library-First | No new application modules; Dockerfile and Compose are tooling with clear boundaries | Pass |
| II. API-First & Testable Services | No API change; backend still exposes same HTTP API; contract documents runtime only | Pass |
| III. Test-First | New container workflow validated by existing tests run against containerized backend where applicable; verification tests for no host exposure | Pass |
| IV. Integration Testing | Integration tests can run with backend in container; contract and quickstart document how to start and connect | Pass |
| V. Observability & Simplicity | Backend logs and errors unchanged; startup failure surfaces via container exit/logs (spec edge case) | Pass |
| Technology Stack | No new stack; React + TypeScript + MUI, Go, MongoDB unchanged. Docker is tooling only | Pass |

## Project Structure

### Documentation (this feature)

```text
specs/003-containerize-backend/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output (runtime config)
├── quickstart.md        # Phase 1 output (run backend in Docker)
├── contracts/           # Phase 1 output (backend runtime contract)
└── tasks.md             # Phase 2 output (/speckit.tasks – not created by this command)
```

### Source Code (repository root)

Existing layout unchanged. Add only:

```text
backend/
├── Dockerfile           # Multi-stage build: build Go binary, run in minimal image (e.g. scratch or alpine)
├── ...                  # (existing backend code unchanged)
docker-compose.yml       # Add backend service: build from backend/, no ports, same network as db, MONGODB_URI=mongodb://db:27017, depends_on db
# Optional: scripts/ for run/verify backend container (e.g. smoke from another container)
```

**Structure Decision**: Add `backend/Dockerfile` and extend root `docker-compose.yml` with a `backend` service. Backend service has no `ports:` mapping so it is not exposed to the host; it is reachable as `backend:8080` from other containers on the same Compose default network. Existing `backend/` application code is unchanged; only container definition and Compose wiring are added.

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

N/A — no violations.
