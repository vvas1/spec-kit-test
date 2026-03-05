# Implementation Plan: Containerize Frontend (Host Port 3000, Container Port 5137)

**Branch**: `004-containerize-frontend` | **Date**: 2026-03-05 | **Spec**: [spec.md](./spec.md)  
**Input**: Feature specification from `/specs/004-containerize-frontend/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/plan-template.md` for the execution workflow.

## Summary

Containerize the React/Vite frontend for local development so it runs as a service inside Docker, reachable from the host at port **3000** with the dev server listening on port **5137** inside the container. Use a Dockerfile that runs the frontend dev server (e.g. Vite) and extend docker-compose with a `frontend` service: port mapping 3000:5137, default API base URL `http://backend:8080` when in the same composition (overridable via `VITE_API_URL`). Development only; no production build or deployment in scope.

## Technical Context

**Language/Version**: Node.js (LTS) and existing frontend stack (TypeScript, React, Vite per constitution); no new application code; container definition and tooling only.  
**Primary Dependencies**: Docker, Docker Compose; existing frontend (Vite dev server, React, MUI).  
**Storage**: N/A (frontend is stateless; API calls go to backend).  
**Testing**: Existing frontend tests; add path to run frontend in container and verify host port 3000 serves the app and API URL is configurable.  
**Target Platform**: Developer machines (Linux, macOS, Windows with Docker).  
**Project Type**: Infrastructure / developer experience (add frontend container to existing web application).  
**Performance Goals**: Frontend ready within 60 seconds of container start (spec SC-001).  
**Constraints**: Host port 3000, container port 5137; dev server only (no production build); default `VITE_API_URL=http://backend:8080` when in same composition.  
**Scale/Scope**: Single frontend instance per Compose project; local dev only.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Gate | Status |
|-----------|------|--------|
| I. Library-First | No new application modules; Dockerfile and Compose are tooling with clear boundaries | Pass |
| II. API-First & Testable Services | No API change; frontend still consumes backend API; contract documents runtime and env only | Pass |
| III. Test-First | New container workflow validated by existing tests; verification that frontend in container loads and can call backend | Pass |
| IV. Integration Testing | Integration tests can run with frontend in container; contract and quickstart document how to start and set API URL | Pass |
| V. Observability & Simplicity | Frontend error/loading states unchanged; startup failure surfaces via container exit/logs | Pass |
| Technology Stack | No new stack; React + TypeScript + MUI, Go, MongoDB unchanged. Docker is tooling only | Pass |

## Project Structure

### Documentation (this feature)

```text
specs/004-containerize-frontend/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output (runtime config)
├── quickstart.md        # Phase 1 output (run frontend in Docker)
├── contracts/           # Phase 1 output (frontend runtime contract)
└── tasks.md             # Phase 2 output (/speckit.tasks – not created by this command)
```

### Source Code (repository root)

Existing layout unchanged. Add only:

```text
frontend/
├── Dockerfile           # Dev server: Node image, install deps, run dev server on port 5137 (e.g. npm run dev with --port 5137 or vite config)
├── .dockerignore        # Optional: exclude node_modules, .env*, etc.
├── ...                  # (existing frontend code unchanged; may add env or config for port 5137 when in container)
docker-compose.yml       # Add frontend service: build from frontend/, ports 3000:5137, VITE_API_URL=http://backend:8080, same network as backend
```

**Structure Decision**: Add `frontend/Dockerfile` and extend root `docker-compose.yml` with a `frontend` service. Port mapping **3000:5137** (host:container). Frontend runs the dev server (e.g. `npm run dev`) configured to listen on 5137 inside the container. When run with backend in the same composition, set `VITE_API_URL=http://backend:8080` by default; document override. Existing frontend application code is unchanged except possibly port configuration for 5137 (e.g. Vite `server.port` or CLI flag).

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

N/A — no violations.
