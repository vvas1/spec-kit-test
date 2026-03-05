# Tasks: Containerize Backend (Docker-Only Access)

**Input**: Design documents from `/specs/003-containerize-backend/`  
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/

**Tests**: No new test tasks requested by spec; validation is via existing backend tests and quickstart/verification steps (start container, verify no host exposure, verify same-network reachability).

**Organization**: Tasks are grouped by user story so each story can be implemented and validated independently.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- Repository root: `docker-compose.yml`, optional `scripts/`
- Backend: `backend/` (existing), `backend/Dockerfile` (new)
- Feature docs: `specs/003-containerize-backend/`

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Add the backend container definition so the backend can be started with Compose and is only reachable from within Docker.

- [x] T001 Create `backend/Dockerfile` with multi-stage build: build stage using Go image (pin version to match backend/go.mod, e.g. 1.23), compile binary; run stage using minimal image (e.g. alpine) to run the binary only; ensure backend listens on :8080 and uses MONGODB_URI from env per specs/003-containerize-backend/research.md

- [x] T002 Add `backend` service to `docker-compose.yml` at repository root: build context `./backend`, Dockerfile `backend/Dockerfile`; no `ports:` (backend not exposed to host); environment `MONGODB_URI=mongodb://db:27017`; `depends_on: db: condition: service_healthy`; same default network as `db` per specs/003-containerize-backend/research.md and specs/003-containerize-backend/contracts/backend-runtime.md

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Document start/stop, readiness, and how other services connect so all user stories can rely on the same contract.

**⚠️ CRITICAL**: No user story work can begin until this phase is complete.

- [x] T003 Update specs/003-containerize-backend/quickstart.md with exact backend start/stop commands: `docker compose up -d db` then `docker compose up -d backend` (or `docker compose up -d`), `docker compose stop backend`, `docker compose down`; add verification step for "from another container" (e.g. one-off container curling http://backend:8080) and note that host connection must fail

- [x] T004 [P] Ensure specs/003-containerize-backend/contracts/backend-runtime.md documents MONGODB_URI in container (`mongodb://db:27017`), no host exposure, discovery as `backend:8080`, and readiness within 30 seconds; align with quickstart commands

**Checkpoint**: Foundation ready — user story implementation and verification can begin

---

## Phase 3: User Story 1 - Run Backend in a Container (Priority: P1) 🎯 MVP

**Goal**: Developer can run the backend as a container without installing Go on the host; backend starts and responds to requests from within the container environment.

**Independent Test**: Run `docker compose up -d db`, wait for db healthy, run `docker compose up -d backend`; from another container on the same network (e.g. `docker compose run --rm <image> curl -s http://backend:8080/` or equivalent), confirm backend responds; stop with `docker compose down`. No Go installed on host.

### Implementation for User Story 1

- [x] T005 [US1] Verify backend starts and responds from within Docker: after T001–T002, run `docker compose up -d db`, wait for healthy, run `docker compose up -d backend`; confirm backend container is running and that a request to http://backend:8080 from another container on the same network succeeds (e.g. 200 or 404); document or add a one-line verification to specs/003-containerize-backend/quickstart.md

- [x] T006 [P] [US1] Document in specs/003-containerize-backend/quickstart.md that the backend container workflow requires Docker/Compose and that no backend runtime (Go) need be installed on the host; add to Prerequisites section

**Checkpoint**: User Story 1 complete — backend runs in container and is reachable from same network; host does not need Go

---

## Phase 4: User Story 2 - Backend Reachable Only From Within Docker (Priority: P2)

**Goal**: Backend is not exposed to the host; connections from the host fail. Connections from other containers on the same network succeed.

**Independent Test**: With backend running via Compose, confirm `curl http://localhost:8080` (or equivalent from host) fails; confirm from another container `curl http://backend:8080` succeeds.

### Implementation for User Story 2

- [x] T007 [US2] Confirm backend service in docker-compose.yml has no `ports:` mapping; add a short verification step in specs/003-containerize-backend/quickstart.md (or contract) that connection from host to backend port must fail (e.g. "From the host: You cannot reach the backend on localhost")

- [x] T008 [P] [US2] Document in specs/003-containerize-backend/contracts/backend-runtime.md and specs/003-containerize-backend/quickstart.md that other services must use hostname `backend` and port `8080` (e.g. `http://backend:8080`) when connecting from another container; no localhost or host IP from inside containers

**Checkpoint**: User Story 2 complete — host cannot reach backend; same-network containers can

---

## Phase 5: User Story 3 - Start and Stop Backend via Container Lifecycle (Priority: P3)

**Goal**: Backend can be started and stopped via standard container lifecycle; when multiple services are in the same composition, backend is discoverable by other services on the same network.

**Independent Test**: Start backend container, verify it is running and reachable at backend:8080 from another container; stop backend container, verify it is no longer running or reachable. With db+backend up, confirm backend is discoverable by name.

### Implementation for User Story 3

- [x] T009 [US3] Document in specs/003-containerize-backend/quickstart.md the exact start command (e.g. `docker compose up -d backend` or `docker compose up -d`), stop command (`docker compose stop backend`), and full teardown (`docker compose down`); include that backend is typically ready within 30 seconds after start (once db is healthy)

- [x] T010 [US3] Verify backend discoverability: with `docker compose up -d` (db + backend), run a temporary container on the same network that resolves `backend` and connects to port 8080 (e.g. curl or wget); confirm backend responds; document result or add to quickstart "Connecting other services to the backend" with hostname `backend` and port 8080

**Checkpoint**: User Story 3 complete — start/stop documented and discoverability verified

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Full validation and edge-case documentation.

- [x] T011 Run full quickstart validation per specs/003-containerize-backend/quickstart.md: start db, start backend, verify backend reachable from another container at http://backend:8080, verify backend not reachable from host (e.g. curl localhost:8080 fails), stop backend and then full down; confirm no broken steps

- [x] T012 Ensure edge-case behavior is documented in specs/003-containerize-backend/quickstart.md or specs/003-containerize-backend/contracts/backend-runtime.md: (1) backend container fails to start — use `docker compose logs backend` and check MONGODB_URI / db connectivity; (2) no container runtime on host — Docker and Compose are required (quickstart Prerequisites); (3) other container cannot reach backend — use hostname `backend` and port 8080 on same Compose network

---

## Dependencies & Execution Order

### Phase Dependencies

- **Phase 1 (Setup)**: No dependencies — add Dockerfile and Compose backend service first.
- **Phase 2 (Foundational)**: Depends on Phase 1 — document commands and contract so US1–US3 can be verified.
- **Phase 3 (US1)**: Depends on Phase 2 — verify backend runs in container and responds from same network.
- **Phase 4 (US2)**: Depends on Phase 2 — verify no host exposure and document discovery.
- **Phase 5 (US3)**: Depends on Phase 2 — document start/stop and verify discoverability.
- **Phase 6 (Polish)**: Depends on Phases 3–5 — full quickstart and edge-case docs.

### User Story Dependencies

- **US1 (P1)**: No dependency on US2/US3 — MVP is "backend runs in container, reachable from same network."
- **US2 (P2)**: Relies on same Compose setup as US1; can be verified after Phase 2.
- **US3 (P3)**: Relies on same Compose setup; start/stop and discovery docs.

### Parallel Opportunities

- T003 and T004: Can run in parallel (quickstart vs contract doc).
- T006: Can run in parallel with T005 (docs vs verification).
- T007 and T008: Can run in parallel (no-ports verification vs discovery docs).
- T009 and T010: T010 (verify discoverability) can follow T009 (docs) or run after backend is up.

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Dockerfile + Compose backend service.
2. Complete Phase 2: Quickstart and contract (exact commands, verification).
3. Complete Phase 3: US1 — verify backend starts and responds from within Docker.
4. **STOP and VALIDATE**: Run independent test for US1 (no host Go, backend reachable from container).
5. Demo: `docker compose up -d`, then curl from another container to http://backend:8080.

### Incremental Delivery

1. Phase 1 + 2 → Foundation ready.
2. Phase 3 (US1) → Backend in container, same-network reachable (MVP).
3. Phase 4 (US2) → Confirm no host exposure, document discovery.
4. Phase 5 (US3) → Start/stop and discoverability documented and verified.
5. Phase 6 → Full quickstart and edge-case docs.

### Task Count Summary

| Phase | Task IDs | Count |
|-------|----------|-------|
| Phase 1 (Setup) | T001–T002 | 2 |
| Phase 2 (Foundational) | T003–T004 | 2 |
| Phase 3 (US1) | T005–T006 | 2 |
| Phase 4 (US2) | T007–T008 | 2 |
| Phase 5 (US3) | T009–T010 | 2 |
| Phase 6 (Polish) | T011–T012 | 2 |
| **Total** | | **12** |

---

## Notes

- [P] tasks = different files or sections, no dependencies on each other.
- [Story] label maps task to spec user story for traceability.
- Each user story is independently verifiable via the Independent Test criteria in the spec.
- No new application code; only backend/Dockerfile, docker-compose.yml changes, and docs in specs/003-containerize-backend/.
