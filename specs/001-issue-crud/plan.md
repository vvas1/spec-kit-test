# Implementation Plan: Create and Edit Issues

**Branch**: `001-issue-crud` | **Date**: 2025-03-04 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/001-issue-crud/spec.md`

## Summary

Implement create, edit, and view flows for issues with title, description, status, and assignee. Backend (Go) exposes a REST API with persisted storage (MongoDB); frontend (React + TypeScript + Material UI) consumes the API with paginated list, detail view, and create/edit forms. Status set: To Do, In Progress, Review, Done (free transitions). Assignees come from an in-app user list. Last-write-wins for concurrent edits; title max 200 chars, description max 10,000; list paginated (e.g. 25 per page).

## Technical Context

**Language/Version**: Go 1.21+ (backend), TypeScript 5.x (frontend)  
**Primary Dependencies**: React 18, Material UI (MUI), Go standard library + MongoDB driver; frontend build (e.g. Vite)  
**Storage**: MongoDB for issues and in-app user/assignee list  
**Testing**: Go: `go test`; frontend: Vitest or Jest + React Testing Library; integration tests for API and frontend–backend  
**Target Platform**: Web (browser); backend runs on Linux/macOS/Windows  
**Project Type**: web-application (frontend + backend)  
**Performance Goals**: Create/edit response under 2s; list page load under 1s; spec target: create issue in under 1 minute (user time)  
**Constraints**: Title ≤200 chars, description ≤10k chars; paginated list (page size 25)  
**Scale/Scope**: Single-tenant issue tracker; in-app user list; no auth in this feature

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Gate | Status |
|-----------|------|--------|
| I. Library-First | Backend: testable Go packages; frontend: TypeScript modules / React components with clear boundaries | Pass |
| II. API-First & Testable Services | REST API documented in contracts/; core logic testable without UI; frontend consumes API only | Pass |
| III. Test-First | TDD for backend and frontend; tests before implementation | Pass |
| IV. Integration Testing | API contracts, DB access, frontend–backend communication covered by integration tests | Pass |
| V. Observability & Simplicity | Structured logging and error responses in backend; frontend shows loading and error states | Pass |
| Technology Stack | React + TypeScript + Material UI (frontend), Go (backend), MongoDB (storage) | Pass |

## Project Structure

### Documentation (this feature)

```text
specs/001-issue-crud/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output (REST API)
└── tasks.md             # Phase 2 output (/speckit.tasks - not created by this command)
```

### Source Code (repository root)

```text
backend/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── model/
│   │   ├── issue.go
│   │   └── user.go
│   ├── store/
│   │   └── mongodb/
│   │       ├── issue.go
│   │       └── user.go
│   ├── api/
│   │   ├── handler/
│   │   │   ├── issues.go
│   │   │   └── users.go
│   │   └── router.go
│   └── service/
│       └── issue.go
└── tests/
    ├── integration/
    └── unit/

frontend/
├── src/
│   ├── components/
│   │   ├── IssueList.tsx
│   │   ├── IssueDetail.tsx
│   │   ├── IssueForm.tsx
│   │   └── ...
│   ├── pages/
│   │   ├── IssueListPage.tsx
│   │   ├── IssueDetailPage.tsx
│   │   └── ...
│   ├── services/
│   │   └── api.ts
│   ├── types/
│   │   └── issue.ts
│   └── App.tsx
└── tests/
    ├── integration/
    └── unit/
```

**Structure Decision**: Web application layout. Backend holds Go packages under `internal/` (model, store, api, service). Frontend holds React/TypeScript under `src/` with components, pages, services, and types. Tests live under `backend/tests/` and `frontend/tests/` for integration and unit.

## Complexity Tracking

> No constitution violations. Leave empty.
