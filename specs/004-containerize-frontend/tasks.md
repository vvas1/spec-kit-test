# Tasks: Containerize Frontend (Host Port 3000, Container Port 5137)

**Input**: Design documents from `/specs/004-containerize-frontend/`  
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/

**Tests**: No new test tasks requested by spec; validation is via existing frontend tests and quickstart/verification steps (start container, reach http://localhost:3000, verify API URL config).

**Organization**: Tasks are grouped by user story so each story can be implemented and validated independently.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- Repository root: `docker-compose.yml`, optional `scripts/`
- Frontend: `frontend/` (existing), `frontend/Dockerfile`, `frontend/.dockerignore` (new)
- Feature docs: `specs/004-containerize-frontend/`

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Add the frontend container definition so the frontend can be started with Compose, reachable at host port 3000 with dev server on container port 5137.

- [x] T001 Create `frontend/Dockerfile`: Node LTS image (e.g. node:20-alpine), WORKDIR, copy package.json and package-lock.json, run `npm ci`, copy source, run dev server on port 5137 (e.g. `npm run dev -- --port 5137` or configure Vite to use 5137) per specs/004-containerize-frontend/research.md. Ensure dev server listens on 5137 inside the container.

- [x] T002 [P] Create `frontend/.dockerignore` to exclude node_modules, .env*, .git, coverage, dist, build, *.log so the image build stays minimal and repeatable

- [x] T003 Add `frontend` service to `docker-compose.yml` at repository root: build context `./frontend`, Dockerfile `frontend/Dockerfile`; ports `3000:5137`; environment `VITE_API_URL=http://backend:8080`; same default network as backend per specs/004-containerize-frontend/research.md and specs/004-containerize-frontend/contracts/frontend-runtime.md

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Document start/stop, port mapping, and API URL so all user stories can rely on the same contract.

**⚠️ CRITICAL**: No user story work can begin until this phase is complete.

- [x] T004 Update specs/004-containerize-frontend/quickstart.md with exact frontend start/stop commands: `docker compose up -d frontend` (or `docker compose up -d` for full stack), `docker compose stop frontend`, `docker compose down`; add verification step (open http://localhost:3000 in browser) and note that VITE_API_URL defaults to http://backend:8080 when in same composition

- [x] T005 [P] Ensure specs/004-containerize-frontend/contracts/frontend-runtime.md documents VITE_API_URL (default http://backend:8080), port mapping 3000:5137, and readiness within 60 seconds; align with quickstart commands

**Checkpoint**: Foundation ready — user story implementation and verification can begin

---

## Phase 3: User Story 1 - Run Frontend in a Container (Priority: P1) 🎯 MVP

**Goal**: Developer can run the frontend as a container without installing Node on the host; frontend starts and is reachable at http://localhost:3000.

**Independent Test**: Run `docker compose up -d frontend` (with db and backend up if testing API); open http://localhost:3000 in browser and confirm the app loads; stop with `docker compose down`. No Node installed on host.

### Implementation for User Story 1

- [x] T006 [US1] Verify frontend starts and is reachable at host port 3000: after T001–T003, run `docker compose up -d frontend` (or `docker compose up -d`); confirm frontend container is running and that http://localhost:3000 serves the application; document or add a one-line verification to specs/004-containerize-frontend/quickstart.md

- [x] T007 [P] [US1] Document in specs/004-containerize-frontend/quickstart.md that the frontend container workflow requires Docker/Compose and that no frontend runtime (Node/npm) need be installed on the host; add to Prerequisites section

**Checkpoint**: User Story 1 complete — frontend runs in container and is reachable at http://localhost:3000; host does not need Node

---

## Phase 4: User Story 2 - Access Frontend on Host Port 3000 (Priority: P2)

**Goal**: Frontend is reachable from the host at port 3000; port mapping 3000:5137 is confirmed and documented.

**Independent Test**: With frontend running via Compose, open http://localhost:3000 and confirm the app loads; confirm docker-compose.yml has ports 3000:5137.

### Implementation for User Story 2

- [x] T008 [US2] Confirm frontend service in docker-compose.yml has ports `3000:5137`; add a short verification step in specs/004-containerize-frontend/quickstart.md that users access the app at http://localhost:3000 (host port 3000 maps to container port 5137)

- [x] T009 [P] [US2] Document in specs/004-containerize-frontend/contracts/frontend-runtime.md and specs/004-containerize-frontend/quickstart.md that VITE_API_URL defaults to http://backend:8080 when frontend and backend are in the same composition, and is overridable via environment variable; document the variable name (VITE_API_URL)

**Checkpoint**: User Story 2 complete — host port 3000 and container port 5137 confirmed; API URL documented

---

## Phase 5: User Story 3 - Start and Stop Frontend via Container Lifecycle (Priority: P3)

**Goal**: Frontend can be started and stopped via standard container lifecycle; when full stack (db + backend + frontend) is in the same composition, frontend starts and is reachable at port 3000.

**Independent Test**: Start frontend container, verify reachable at http://localhost:3000; stop frontend container, verify no longer reachable. With `docker compose up -d`, confirm frontend is reachable at port 3000.

### Implementation for User Story 3

- [x] T010 [US3] Document in specs/004-containerize-frontend/quickstart.md the exact start command (e.g. `docker compose up -d frontend` or `docker compose up -d`), stop command (`docker compose stop frontend`), and full teardown (`docker compose down`); include that frontend is typically ready within 60 seconds after start

- [x] T011 [US3] Verify full-stack composition: with `docker compose up -d` (db + backend + frontend), open http://localhost:3000 and confirm the app loads and can call the backend (e.g. load issues or users); document result or add to quickstart "Configuring the API base URL" with default http://backend:8080

**Checkpoint**: User Story 3 complete — start/stop documented and full-stack reachability verified

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Full validation and edge-case documentation.

- [x] T012 Run full quickstart validation per specs/004-containerize-frontend/quickstart.md: start full stack, open http://localhost:3000, verify app loads and API calls work (if backend is up), stop frontend and then full down; confirm no broken steps

- [x] T013 Ensure edge-case behavior is documented in specs/004-containerize-frontend/quickstart.md or specs/004-containerize-frontend/contracts/frontend-runtime.md: (1) frontend container fails to start (e.g. port 3000 in use, build failure) — use `docker compose logs frontend`; (2) no container runtime on host — Docker and Compose are required (quickstart Prerequisites); (3) API calls fail — ensure backend is running and VITE_API_URL is correct (default http://backend:8080 when in same composition)

---

## Dependencies & Execution Order

### Phase Dependencies

- **Phase 1 (Setup)**: No dependencies — add Dockerfile, .dockerignore, and Compose frontend service first.
- **Phase 2 (Foundational)**: Depends on Phase 1 — document commands and contract so US1–US3 can be verified.
- **Phase 3 (US1)**: Depends on Phase 2 — verify frontend runs in container and is reachable at port 3000.
- **Phase 4 (US2)**: Depends on Phase 2 — verify port mapping 3000:5137 and document VITE_API_URL.
- **Phase 5 (US3)**: Depends on Phase 2 — document start/stop and verify full-stack.
- **Phase 6 (Polish)**: Depends on Phases 3–5 — full quickstart and edge-case docs.

### User Story Dependencies

- **US1 (P1)**: No dependency on US2/US3 — MVP is "frontend runs in container, reachable at http://localhost:3000."
- **US2 (P2)**: Relies on same Compose setup as US1; port mapping and API URL docs.
- **US3 (P3)**: Relies on same Compose setup; start/stop and full-stack verification.

### Parallel Opportunities

- T002 can run in parallel with T001 (.dockerignore vs Dockerfile).
- T004 and T005: Can run in parallel (quickstart vs contract doc).
- T007 can run in parallel with T006 (docs vs verification).
- T008 and T009: Can run in parallel (port verification vs API URL docs).

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Dockerfile (dev server on 5137) + .dockerignore + Compose frontend service (3000:5137, VITE_API_URL).
2. Complete Phase 2: Quickstart and contract (exact commands, verification).
3. Complete Phase 3: US1 — verify frontend starts and is reachable at http://localhost:3000.
4. **STOP and VALIDATE**: Run independent test for US1 (no host Node, frontend at port 3000).
5. Demo: `docker compose up -d frontend`, then open http://localhost:3000.

### Incremental Delivery

1. Phase 1 + 2 → Foundation ready.
2. Phase 3 (US1) → Frontend in container, reachable at 3000 (MVP).
3. Phase 4 (US2) → Port mapping and VITE_API_URL documented.
4. Phase 5 (US3) → Start/stop and full-stack documented and verified.
5. Phase 6 → Full quickstart and edge-case docs.

### Task Count Summary

| Phase | Task IDs | Count |
|-------|----------|-------|
| Phase 1 (Setup) | T001–T003 | 3 |
| Phase 2 (Foundational) | T004–T005 | 2 |
| Phase 3 (US1) | T006–T007 | 2 |
| Phase 4 (US2) | T008–T009 | 2 |
| Phase 5 (US3) | T010–T011 | 2 |
| Phase 6 (Polish) | T012–T013 | 2 |
| **Total** | | **13** |

---

## Notes

- [P] tasks = different files or sections, no dependencies on each other.
- [Story] label maps task to spec user story for traceability.
- Each user story is independently verifiable via the Independent Test criteria in the spec.
- No new application logic; only frontend/Dockerfile, frontend/.dockerignore, docker-compose.yml changes, and docs in specs/004-containerize-frontend/. Dev server must listen on **5137** inside the container (current Vite config uses 5173; override via CLI `--port 5137` or config).
