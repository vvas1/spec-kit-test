<!--
Sync Impact Report (2025-03-04)
- Version change: 1.1.0 → 1.2.0 (frontend UI library)
- Modified principles: none
- Added sections: none
- Removed sections: none
- Templates: plan-template.md ✅; spec-template.md ✅; tasks-template.md ✅
- Follow-up TODOs: none
-->

# Issue Tracker Constitution

## Technology Stack

The project is an **issue tracker** built with:

- **React + TypeScript + Material UI**: Frontend UI; all frontend code in TypeScript; UI components and theming from Material UI (MUI).
- **Go**: Backend services, APIs, and business logic.
- **MongoDB**: Primary data store for issues, users, and related entities.

No other frontend framework, UI component library, backend language, or primary database may be introduced without an amendment to this constitution. Build, tooling, and API style (e.g. REST, GraphQL) are chosen in the implementation plan.

## Core Principles

### I. Library-First

Every feature starts as a standalone library or module. Backend logic MUST be in testable Go packages; frontend logic in TypeScript modules or React components with clear boundaries. Each module MUST have a clear purpose; organizational-only or placeholder code is not permitted.

*Rationale*: Enables reuse, clear boundaries, and incremental delivery.

### II. API-First & Testable Services

Backend services MUST expose functionality via a defined API (e.g. REST or GraphQL). Contracts MUST be documented and versioned. Core logic MUST be testable without the UI; Go packages MUST support unit tests. Where scripting or tooling is needed, a CLI or scriptable entry point MAY be provided (stdin/args → stdout, stderr; JSON where applicable). Frontend consumes only the public API.

*Rationale*: Enables testing, clear frontend/backend separation, and consistent behavior.

### III. Test-First (NON-NEGOTIABLE)

TDD is mandatory for backend and frontend behavior. Flow: tests written → user or spec approval → tests fail → then implement. Red–Green–Refactor cycle MUST be followed; implementation MUST NOT precede failing tests for the behavior in scope.

*Rationale*: Ensures requirements are executable and prevents untested code paths.

### IV. Integration Testing

Integration tests are REQUIRED for: API contract boundaries, database access and schemas, frontend–backend communication, and shared models. Unit tests alone do not satisfy these areas.

*Rationale*: Catches contract and integration failures before production.

### V. Observability & Simplicity

Structured logging and clear error responses MUST be used in the backend; frontend MUST handle errors and loading states visibly. Prefer simplicity: YAGNI applies; complexity MUST be justified against this constitution and the current feature spec.

*Rationale*: Keeps the system operable and limits unnecessary scope creep.

## Additional Constraints

Technology stack, compliance standards, and deployment policies MUST align with the feature spec and plan. No technology or dependency may be introduced without coverage in the implementation plan or an approved amendment. The stack (React + TypeScript + Material UI, Go, MongoDB) is fixed unless amended.

## Development Workflow

Code review MUST verify compliance with this constitution. Quality gates MUST include: passing tests (including integration where required), constitution check (see plan template), and completion of mandatory spec sections. Deployment or release approval assumes these gates have been met.

## Governance

This constitution supersedes ad-hoc practices for the scope of this project. Amendments require: documented rationale, explicit approval (e.g., maintainer or team), and a migration or compatibility note where existing work is affected. All PRs and reviews MUST verify compliance; deviations require explicit justification and, if accepted, an update to the constitution or a documented exception.

**Version**: 1.2.0 | **Ratified**: 2025-03-04 | **Last Amended**: 2025-03-04
