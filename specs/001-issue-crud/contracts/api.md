# REST API Contract: Issues and Users

**Feature**: 001-issue-crud  
**Version**: 1.0  
**Base path**: `/api` (or as configured)

All request/response bodies are JSON. Errors use HTTP status + JSON body `{ "error": "<message>" }`.

---

## Issues

### List issues (paginated)

- **GET** `/api/issues`
- **Query**: `page` (1-based, default 1), `limit` (default 25, max 50)
- **Response 200**: `{ "items": [ Issue ], "total": number, "page": number, "limit": number }`
- **Issue** (in list): `{ "id", "title", "description", "status", "assigneeId", "assigneeName" (optional), "createdAt", "updatedAt" }`
- Order: `updatedAt` descending (newest first)

### Get one issue

- **GET** `/api/issues/:id`
- **Response 200**: single Issue object (full)
- **Response 404**: `{ "error": "issue not found" }`

### Create issue

- **POST** `/api/issues`
- **Body**: `{ "title": string (required), "description": string (optional), "status": string (optional), "assigneeId": string (optional) }`
- **Validation**: title required, 1–200 chars; description 0–10k chars; status one of "To Do" | "In Progress" | "Review" | "Done" or omit (default "To Do"); assigneeId must be valid user id or omit
- **Response 201**: created Issue object (with id, createdAt, updatedAt)
- **Response 400**: `{ "error": "<validation message>" }`

### Update issue

- **PUT** `/api/issues/:id`
- **Body**: same as create (all fields optional for PATCH semantics; for PUT, full replace: title required, others optional)
- **Semantics**: Full replace; title required; last-write-wins
- **Response 200**: updated Issue object
- **Response 400**: validation error
- **Response 404**: issue not found

---

## Users (assignee list)

### List users

- **GET** `/api/users`
- **Response 200**: `{ "items": [ { "id": string, "name": string } ] }`
- Used for assignee dropdown; no pagination required for initial scope.

---

## Common

- **Content-Type**: `application/json` for request and response bodies.
- **CORS**: Backend must allow frontend origin (configured in deployment).
- **Validation errors**: HTTP 400, body `{ "error": "<message>" }` (e.g. "title is required", "title must be at most 200 characters").
