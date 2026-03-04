# Data Model: Create and Edit Issues

**Feature**: 001-issue-crud  
**Date**: 2025-03-04

## Entities

### Issue

| Field       | Type     | Required | Constraints                    | Notes                          |
|------------|----------|----------|--------------------------------|--------------------------------|
| id         | string   | yes      | stable, unique                 | Generated (e.g. MongoDB ObjectID or UUID) |
| title      | string   | yes      | length 1–200                    | User-facing title               |
| description| string   | no       | length 0–10,000                | Optional long text              |
| status     | string   | yes      | one of enum below              | Default: "To Do"                |
| assigneeId | string   | no       | references User.id or empty    | Empty = unassigned              |
| createdAt  | datetime | yes      | set on create                  | Immutable                       |
| updatedAt  | datetime | yes      | set on create, update on save  | Last-write-wins timestamp       |

**Status enum**: `To Do` | `In Progress` | `Review` | `Done`. Transitions: free (any → any).

**Validation (from spec)**:
- Title required on create and update; reject empty or whitespace-only.
- Title max 200 characters; description max 10,000; reject and return validation error if exceeded.
- AssigneeId must be a valid User id or empty; if invalid, allow save but treat as unassigned (or reject per product choice; spec says invalid assignee must not block saving other fields — here we allow empty and valid id only; invalid id can be treated as "unassign" or validation error).

**Persistence**: MongoDB collection `issues`. Indexes: `_id` (primary), `updatedAt` (for list ordering/pagination).

---

### User (assignee list)

| Field   | Type   | Required | Constraints   | Notes                |
|--------|--------|----------|---------------|----------------------|
| id     | string | yes      | stable, unique| Generated             |
| name   | string | yes      | non-empty     | Display name in UI   |

Used only as assignee list for this feature. No authentication. Stored in MongoDB collection `users`. How users are added is out of scope (seed data or minimal admin endpoint later).

---

## Relationships

- **Issue.assigneeId** → **User.id** (optional). Many issues can reference one user. Display "Unassigned" when assigneeId is empty or user not found.

---

## Lifecycle / State

- **Issue**: Created with default status "To Do" if not provided. Updated in place; last-write-wins; no soft delete in scope.
- **User**: Out of scope for lifecycle in this feature (read-only list for assignee dropdown).
