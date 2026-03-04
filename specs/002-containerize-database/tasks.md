# Tasks: Containerize Database

**Input**: Design documents from `/specs/002-containerize-database/`  
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/

**Tests**: No new test tasks requested by spec; validation is via existing backend/frontend tests run against the containerized DB and quickstart flow.

**Organization**: Tasks are grouped by user story so each story can be implemented and validated independently.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- Repository root: `docker-compose.yml`, optional `scripts/`
- Backend: `backend/` (existing)
- Feature docs: `specs/002-containerize-database/`

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Add the container definition so the database can be started with one command.

- [x] T001 Create docker-compose.yml at repository root with MongoDB service: image mongo:7, port 27017, bind 127.0.0.1, named volume for data, healthcheck using mongosh (e.g. `mongosh --eval "db.adminCommand('ping')"`) per research.md and contracts/database-runtime.md

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Document start/stop, readiness, and connection so all user stories can rely on the same contract.

**⚠️ CRITICAL**: No user story work can begin until this phase is complete.

- [x] T002 Update specs/002-containerize-database/quickstart.md with exact start/stop commands (docker compose up -d db, docker compose down, docker compose down -v) and readiness check steps matching the healthcheck in docker-compose.yml
- [x] T003 [P] Document connection parameters (MONGODB_URI default localhost:27017, optional MONGO_INITDB_ROOT_* and bind/port override) in specs/002-containerize-database/quickstart.md and/or specs/002-containerize-database/contracts/database-runtime.md

**Checkpoint**: Foundation ready — user story implementation can begin

---

## Phase 3: User Story 1 - Start Database for Local Development (Priority: P1) 🎯 MVP

**Goal**: Developer can start the database with minimal steps and the application connects successfully.

**Independent Test**: Run `docker compose up -d db`, wait for healthy, start backend with MONGODB_URI=mongodb://localhost:27017, confirm backend connects and can read/write; stop with `docker compose down`.

### Implementation for User Story 1

- [x] T004 [US1] Verify backend connects to containerized MongoDB: start DB via docker compose from repo root, ensure backend/internal/store/mongodb uses MONGODB_URI (default mongodb://localhost:27017), run backend and confirm connectivity (e.g. start server and one API call)
- [x] T005 [P] [US1] Add .env.example at repository root (or backend/.env.example) with MONGODB_URI and optional MONGO_INITDB_ROOT_USERNAME, MONGO_INITDB_ROOT_PASSWORD per specs/002-containerize-database/contracts/database-runtime.md

**Checkpoint**: User Story 1 complete — new developer can start DB and connect app within two minutes using docs only

---

## Phase 4: User Story 2 - Run Database in Automated Tests (Priority: P2)

**Goal**: Test suite can start a containerized DB before tests and tear it down after, with no manual setup; parallel or back-to-back runs get distinct instance or clean state.

**Independent Test**: Invoke test runner (or script) that starts DB, waits for readiness, runs backend tests, then runs docker compose down -v (or uses COMPOSE_PROJECT_NAME); repeat and confirm no leftover state.

### Implementation for User Story 2

- [x] T006 [US2] Add script or CI step that starts containerized DB (docker compose up -d db), waits for readiness (poll docker compose ps for healthy or documented probe), runs backend test suite (e.g. cd backend && go test ./...), then tears down (docker compose down -v) in scripts/ or backend/scripts/ or document in README/CI config
- [x] T007 [P] [US2] Document COMPOSE_PROJECT_NAME and ephemeral (docker compose down -v) for isolated test runs in specs/002-containerize-database/quickstart.md and specs/002-containerize-database/contracts/database-runtime.md

**Checkpoint**: User Story 2 complete — full automated test suite runs against containerized DB with no manual setup

---

## Phase 5: User Story 3 - Persist Data Across Restarts (Priority: P3)

**Goal**: Persistent mode keeps data across stop/start; ephemeral mode gives clean state for tests.

**Independent Test**: With default Compose (volume): write data, docker compose down, docker compose up -d db, confirm data present. With ephemeral (down -v then up): confirm data absent after restart.

### Implementation for User Story 3

- [x] T008 [P] [US3] Document persistent mode (default: named volume; docker compose down keeps data) in specs/002-containerize-database/quickstart.md
- [x] T009 [US3] Document ephemeral mode (docker compose down -v or no-volume run) and verify clean state after restart in specs/002-containerize-database/quickstart.md; optionally add Compose profile for ephemeral if desired

**Checkpoint**: User Story 3 complete — persistence and ephemeral both documented and verifiable

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Validation and edge-case behavior.

- [x] T010 Run full quickstart validation: start DB per specs/002-containerize-database/quickstart.md, start backend, start frontend, perform one read/write, stop DB; confirm no broken steps
- [x] T011 Ensure edge-case error behavior is documented or implemented: port-in-use and duplicate start yield clear, actionable messages (document in specs/002-containerize-database/contracts/database-runtime.md or verify Compose/scripts surface clear errors)

---

## Dependencies & Execution Order

### Phase Dependencies

- **Phase 1 (Setup)**: No dependencies — add docker-compose.yml first.
- **Phase 2 (Foundational)**: Depends on Phase 1 — document commands and params so US1–US3 can rely on them.
- **Phase 3 (US1)**: Depends on Phase 2 — verify and document local dev flow.
- **Phase 4 (US2)**: Depends on Phase 2 — add test script and isolation docs.
- **Phase 5 (US3)**: Depends on Phase 2 — document persistence vs ephemeral.
- **Phase 6 (Polish)**: Depends on Phases 3–5 — validate and edge cases.

### User Story Dependencies

- **US1 (P1)**: No dependency on US2/US3 — MVP is “start DB, connect app.”
- **US2 (P2)**: No dependency on US1 implementation (same Compose); can be done after Foundational.
- **US3 (P3)**: No dependency on US1/US2 — docs and volume behavior only.

### Parallel Opportunities

- T002 and T003: Can run in parallel (different doc sections).
- T005: Can run in parallel with T004 (different files).
- T007: Can run in parallel with T006 (docs vs script).
- T008 and T009: Both docs; T008 [P] with T009 (T009 may include verification step).

---

## Parallel Example: User Story 1

```bash
# After T001–T003 complete:
# T004 (verify backend) and T005 (.env.example) can be done in parallel
Task T004: Verify backend connects to containerized MongoDB
Task T005: Add .env.example with MONGODB_URI and optional MONGO_* vars
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Add docker-compose.yml.
2. Complete Phase 2: Document start/stop/readiness and connection params.
3. Complete Phase 3: Verify backend connection and add .env.example.
4. **STOP and VALIDATE**: Follow quickstart; new developer can start DB and connect in under two minutes.

### Incremental Delivery

1. Phase 1 + 2 → Foundation (single way to start/stop, documented).
2. Phase 3 (US1) → Local dev flow verified (MVP).
3. Phase 4 (US2) → CI/test runs against containerized DB.
4. Phase 5 (US3) → Persistence and ephemeral documented.
5. Phase 6 → Quickstart and edge-case validation.

### Format Validation

- All tasks use checklist format: `- [ ]`, Task ID (T001–T011), [P] where parallelizable, [USn] for story phases, and include file paths in descriptions.

---

## Notes

- [P] tasks = different files or doc sections, no ordering dependency.
- [USn] maps each task to the user story for traceability.
- No new application code; only docker-compose.yml, optional scripts, and docs.
- Commit after each task or logical group; stop at any checkpoint to validate that story independently.
