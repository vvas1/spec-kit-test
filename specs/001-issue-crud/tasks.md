# Tasks: Create and Edit Issues

**Input**: Design documents from `/specs/001-issue-crud/`  
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/

**Organization**: Tasks are grouped by user story so each story can be implemented and tested independently. Constitution requires TDD: tests are written first and must fail before implementation.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: User story (US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- **Backend**: `backend/` (Go: cmd/, internal/, tests/)
- **Frontend**: `frontend/` (Vite + React + TypeScript: src/, tests/)

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure per plan.md

- [x] T001 Create backend directory structure (backend/cmd/server, backend/internal/model, backend/internal/store/mongodb, backend/internal/api/handler, backend/internal/api, backend/internal/service, backend/tests/unit, backend/tests/contract, backend/tests/integration)
- [x] T002 Create frontend directory structure (frontend/src/components, frontend/src/pages, frontend/src/services, frontend/src/types, frontend/tests)
- [x] T003 Initialize Go module and add MongoDB driver in backend/ (go mod init, go get go.mongodb.org/mongo-driver)
- [x] T004 Initialize Vite + React + TypeScript project in frontend/ (npm create vite@latest or equivalent)
- [x] T005 [P] Add Material UI and required dependencies in frontend/ (package.json)
- [x] T006 [P] Configure ESLint and Prettier in frontend/
- [x] T007 Implement backend entrypoint that starts HTTP server in backend/cmd/server/main.go (placeholder router OK)

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure required before any user story. No user story work until this phase is complete.

- [x] T008 Implement MongoDB connection and config from env (MONGODB_URI) in backend/internal/store/mongodb or backend/cmd/server
- [x] T009 Implement API router with /api prefix and JSON request/response middleware in backend/internal/api/router.go
- [x] T010 Implement CORS middleware allowing frontend origin in backend/internal/api/
- [x] T011 Implement centralized HTTP error handling and JSON error body `{ "error": "<message>" }` in backend/internal/api/
- [x] T012 Add structured logging for HTTP requests in backend/
- [x] T013 Create frontend API client base (fetch wrapper, base URL from env VITE_API_URL) in frontend/src/services/api.ts
- [x] T014 Create App shell with MUI theme and routing placeholder in frontend/src/App.tsx

**Checkpoint**: Foundation ready — user story implementation can begin.

---

## Phase 3: User Story 1 — Create a new issue (Priority: P1) — MVP

**Goal**: User can create an issue with title, description, optional status, and optional assignee; validation for required title and length.

**Independent Test**: Submit create form with title and description; verify API returns 201 and GET issue returns the created issue with default status.

### Tests for User Story 1 (TDD — write first, ensure they fail)

- [x] T015 [P] [US1] Write unit tests for Issue create validation (title required, title 1–200 chars, description 0–10k) in backend/tests/unit/issue_test.go
- [x] T016 [P] [US1] Write contract test for POST /api/issues (201 with body, 400 for missing title, 400 for title too long) in backend/tests/contract/issues_test.go
- [x] T017 [P] [US1] Write contract test for GET /api/users (200 with items array) in backend/tests/contract/users_test.go

### Implementation for User Story 1

- [x] T018 [P] [US1] Create Issue model (id, title, description, status, assigneeId, createdAt, updatedAt) in backend/internal/model/issue.go
- [x] T019 [P] [US1] Create User model (id, name) in backend/internal/model/user.go
- [x] T020 [US1] Implement Issue MongoDB store (Create, GetByID) in backend/internal/store/mongodb/issue.go
- [x] T021 [US1] Implement User MongoDB store (List) in backend/internal/store/mongodb/user.go
- [x] T022 [US1] Implement Issue service Create with validation and default status "To Do" in backend/internal/service/issue.go
- [x] T023 [US1] Implement POST /api/issues handler in backend/internal/api/handler/issues.go
- [x] T024 [US1] Implement GET /api/users handler in backend/internal/api/handler/users.go
- [x] T025 [US1] Register POST /api/issues and GET /api/users in backend/internal/api/router.go
- [x] T026 [P] [US1] Create frontend types (Issue, User, CreateIssueInput) in frontend/src/types/issue.ts
- [x] T027 [US1] Add createIssue and getUsers to frontend API client in frontend/src/services/api.ts
- [x] T028 [US1] Implement IssueForm component (title, description, status select, assignee select) with validation in frontend/src/components/IssueForm.tsx
- [x] T029 [US1] Implement CreateIssuePage and route in frontend/src/pages/CreateIssuePage.tsx
- [x] T030 [US1] Add navigation to create issue in frontend/src/App.tsx

**Checkpoint**: User Story 1 complete — create flow testable end-to-end.

---

## Phase 4: User Story 2 — Edit an existing issue (Priority: P2)

**Goal**: User can open an issue, change title/description/status/assignee, save; validation on save; last-write-wins.

**Independent Test**: Create issue via API, PUT update with new title/status; GET returns updated issue.

### Tests for User Story 2

- [ ] T031 [P] [US2] Write contract test for PUT /api/issues/:id (200 with updated body, 400 missing title, 404 not found) in backend/tests/contract/issues_test.go

### Implementation for User Story 2

- [ ] T032 [US2] Add Update method to Issue MongoDB store in backend/internal/store/mongodb/issue.go
- [ ] T033 [US2] Add Update to Issue service with validation in backend/internal/service/issue.go
- [ ] T034 [US2] Implement PUT /api/issues/:id and GET /api/issues/:id handlers in backend/internal/api/handler/issues.go
- [ ] T035 [US2] Register GET /api/issues/:id and PUT /api/issues/:id in backend/internal/api/router.go
- [ ] T036 [US2] Add getIssue and updateIssue to frontend API client in frontend/src/services/api.ts
- [ ] T037 [US2] Implement IssueDetailPage with load-by-id, edit form, and save in frontend/src/pages/IssueDetailPage.tsx
- [ ] T038 [US2] Add route and navigation for issue detail (by id) in frontend/src/App.tsx

**Checkpoint**: User Story 2 complete — edit flow testable.

---

## Phase 5: User Story 3 — View issue list and detail (Priority: P3)

**Goal**: User sees paginated list of issues (title, status, assignee) and can open one to view full detail.

**Independent Test**: GET /api/issues?page=1&limit=25 returns items and total; GET /api/issues/:id returns full issue; list page shows rows and pagination; detail page shows full issue.

### Tests for User Story 3

- [ ] T039 [P] [US3] Write contract test for GET /api/issues (query page, limit; response items, total, page, limit) in backend/tests/contract/issues_test.go

### Implementation for User Story 3

- [ ] T040 [US3] Add List with pagination (page, limit, total count) to Issue MongoDB store in backend/internal/store/mongodb/issue.go
- [ ] T041 [US3] Add List (paginated) to Issue service in backend/internal/service/issue.go
- [ ] T042 [US3] Implement GET /api/issues handler (query page, limit; default 25) in backend/internal/api/handler/issues.go
- [ ] T043 [US3] Register GET /api/issues in backend/internal/api/router.go
- [ ] T044 [US3] Add listIssues (page, limit) to frontend API client in frontend/src/services/api.ts
- [ ] T045 [US3] Implement IssueList component with table/cards and pagination controls in frontend/src/components/IssueList.tsx
- [ ] T046 [US3] Implement IssueListPage with pagination in frontend/src/pages/IssueListPage.tsx
- [ ] T047 [US3] Wire IssueDetailPage as view-only mode from list (link by id) and set list as default route in frontend/src/App.tsx
- [ ] T048 [US3] Show empty state and link to create when no issues in frontend/src/pages/IssueListPage.tsx

**Checkpoint**: User Story 3 complete — list and detail view testable.

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Constitution alignment and validation

- [ ] T049 [P] Ensure loading and error states are visible in frontend (per constitution V) in frontend/src/components and pages
- [ ] T050 [P] Run quickstart validation (backend + frontend start, create one issue, list, edit) per specs/001-issue-crud/quickstart.md
- [ ] T051 Add integration test for full flow (create → list → open → edit → save) in backend/tests/integration/ or frontend/tests/integration/

---

## Dependencies & Execution Order

### Phase Dependencies

- **Phase 1 (Setup)**: No dependencies — start immediately.
- **Phase 2 (Foundational)**: Depends on Phase 1 — blocks all user stories.
- **Phase 3 (US1)**: Depends on Phase 2 — no other story required.
- **Phase 4 (US2)**: Depends on Phase 2; reuses US1 Issue/User model and store.
- **Phase 5 (US3)**: Depends on Phase 2; reuses US1/US2 handlers (GET :id) and adds list.
- **Phase 6 (Polish)**: Depends on Phases 3–5 as needed.

### User Story Dependencies

- **US1 (P1)**: After Foundational only — independently testable (create issue via API or UI).
- **US2 (P2)**: After Foundational; uses US1 Issue/User; independently testable (PUT + GET :id).
- **US3 (P3)**: After Foundational; uses US1/US2; independently testable (GET list, GET :id, list/detail UI).

### Parallel Opportunities

- Phase 1: T005, T006 [P]; T018, T019 [P] in US1; T026 [P] in US1.
- Phase 2: None (sequential for clarity).
- After Phase 2: US1, US2, US3 can be staffed in parallel; within US1, tests T015–T017 [P], models T018–T019 [P].

---

## Implementation Strategy

### MVP First (User Story 1 only)

1. Complete Phase 1 (Setup) and Phase 2 (Foundational).
2. Complete Phase 3 (US1): write tests (T015–T017), then implement backend and frontend (T018–T030).
3. Validate: create issue via UI, verify via API or list later.
4. Stop and deploy/demo if desired.

### Incremental Delivery

1. Setup + Foundational → foundation ready.
2. US1 → create flow → deploy (MVP).
3. US2 → edit flow → deploy.
4. US3 → list + detail view → deploy.
5. Polish → done.

### Parallel Team Strategy

- One developer: Phases 1 → 2 → 3 → 4 → 5 → 6 in order.
- Multiple developers: After Phase 2, Dev A: US1, Dev B: US2, Dev C: US3 (with coordination on shared backend handlers).

---

## Notes

- Every task uses format `- [ ] Tnnn [P?] [USn?] Description with file path`.
- [P] = parallelizable; [US1/US2/US3] = user story for traceability.
- TDD: run tests after T015–T017, T031, T039 and ensure they fail before implementing corresponding handlers and stores.
- Commit after each task or logical group.
