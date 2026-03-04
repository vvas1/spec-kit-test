# Quickstart: Issue Tracker (001-issue-crud)

**Feature**: Create and edit issues  
**Stack**: React + TypeScript + Material UI (frontend), Go (backend), MongoDB (storage)

## Prerequisites

- Go 1.26
- Node.js 22 and npm/pnpm
- MongoDB running locally (or connection string for remote)

## Backend

```bash
cd backend
# Install dependencies (if using go mod)
go mod tidy
# Run server (set MONGODB_URI if not default localhost)
go run ./cmd/server
```

Default: API listens on `http://localhost:8080` (or as set in env). Expects MongoDB at `localhost:27017` unless `MONGODB_URI` is set.

## Frontend

```bash
cd frontend
npm install   # or pnpm install
npm run dev   # or pnpm dev (Vite default: http://localhost:5173)
```

Set API base URL via env (e.g. `VITE_API_URL=http://localhost:8080/api`) if backend is not at default.

## First Run

1. Start MongoDB.
2. Seed users (if required): add at least one user to `users` collection for assignee dropdown, or implement a minimal `POST /api/users` for seeding.
3. Start backend, then frontend.
4. Open frontend URL; create an issue (title + optional description, status, assignee); list and edit from the UI.

## Tests

- **Backend**: `cd backend && go test ./...`
- **Frontend**: `cd frontend && npm run test` (or pnpm test)
- **Integration**: Run API tests against a test MongoDB instance; run frontend e2e or integration tests against local backend.

## Key Files

- **Spec**: [spec.md](./spec.md)
- **Plan**: [plan.md](./plan.md)
- **Data model**: [data-model.md](./data-model.md)
- **API contract**: [contracts/api.md](./contracts/api.md)
