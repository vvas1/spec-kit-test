#!/usr/bin/env bash
# Run backend tests against the containerized MongoDB.
# Usage: from repo root, ./scripts/test-with-db.sh
# Starts DB, waits for healthy, runs "go test ./..." in backend, then tears down (docker compose down -v).

set -e

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$REPO_ROOT"

echo "[test-with-db] Starting MongoDB..."
docker compose up -d db

echo "[test-with-db] Waiting for DB to be healthy..."
for i in 1 2 3 4 5 6 7 8 9 10; do
  if docker compose ps 2>/dev/null | grep -q "healthy"; then
    echo "[test-with-db] DB is healthy."
    break
  fi
  if [ "$i" -eq 10 ]; then
    echo "[test-with-db] Timeout waiting for healthy DB." >&2
    docker compose down -v
    exit 1
  fi
  sleep 3
done
sleep 2

echo "[test-with-db] Running backend tests..."
(cd backend && go test ./...)
TEST_EXIT=$?

echo "[test-with-db] Tearing down MongoDB (docker compose down -v)..."
docker compose down -v

exit "$TEST_EXIT"
