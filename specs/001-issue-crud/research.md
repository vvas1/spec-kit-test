# Research: Create and Edit Issues

**Feature**: 001-issue-crud  
**Date**: 2025-03-04

## API Style

- **Decision**: REST API for backend–frontend communication.
- **Rationale**: Constitution allows "e.g. REST, GraphQL"; REST is simpler for CRUD-heavy issue tracker, easier to test and document, and aligns with API-First principle.
- **Alternatives considered**: GraphQL (deferred; add later if query flexibility becomes a requirement).

## Page Size (List Pagination)

- **Decision**: 25 issues per page.
- **Rationale**: Spec allows "e.g. 20–50"; 25 is a common default, balances load time and scroll depth.
- **Alternatives considered**: 20 (more requests), 50 (heavier first load).

## Status Set

- **Decision**: Four statuses: **To Do**, **In Progress**, **Review**, **Done**. Free transitions (any status → any other).
- **Rationale**: Spec clarification chose "different set (e.g. To Do, In Progress, Review, Done)" with exact list in implementation plan; free transitions assumed unless specified.
- **Alternatives considered**: Three statuses (e.g. Open, In Progress, Done); restricted transitions (deferred).

## In-App User List (Assignees)

- **Decision**: Minimal **User** entity stored in MongoDB: id, display name (and optional email if needed later). No authentication or roles in this feature. Backend exposes `GET /users` for assignee dropdown; users can be seeded or added via a minimal admin path (out-of-scope for this feature: assume seeded or single endpoint to add a user).
- **Rationale**: Spec: "in-app list only; no external directory or authentication required." Storing users in the same DB keeps the feature self-contained.
- **Alternatives considered**: External LDAP/IdP (rejected per spec); free-text assignee (rejected per spec).

## Error and Validation Responses

- **Decision**: API returns HTTP 4xx with JSON body `{ "error": "<message>" }` for validation and business errors; 400 for validation (e.g. missing title, length exceeded), 404 for missing issue/user.
- **Rationale**: Consistent with REST and constitution (clear error responses); frontend can show messages in UI.
- **Alternatives considered**: HTTP status only with no body (worse DX); custom error codes (unnecessary for first version).
