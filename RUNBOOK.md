For architectural decisions and technical reasoning, see [DECISIONS.md](./DECISIONS.md).

# Auth-Service Backend Runbook

## Purpose

This runbook describes operational procedures for running, monitoring, and troubleshooting the **auth and address API** service in the `backend/` directory.

---

## Service Overview

- **Language / Framework**: Go, Echo
- **Module**: `server`
- **Database**: PostgreSQL
- **Port**: `SERVER_PORT` (default `8080`)
- **Auth**: JWT in HTTP-only cookies + CSRF token
- **Containerization**: Docker image via `Dockerfile`

Configuration is environment-driven via `internal/config/config.go` and `.env` (when using `make` targets). See `backend/README.md` for a full variable list and example.

---

## Prerequisites

- PostgreSQL instance reachable from the service
- `migrate` CLI installed and available on `PATH`
- `.env` file in `backend/` with DB and app settings

At minimum, the following environment variables must be set:

- `DB_USER`, `DB_PWD`, `DB_NAME`
- `JWT_SECRET`, `SESSION_KEY`

---

## Starting the Service

All commands below assume `pwd` is `backend/`.

### Local (development)

Preferred (uses `.env` and compiled binary):

```bash
make run
```

This will:

- Build the `server` binary from `internal/cmd/main.go`
- Source `.env`
- Start the service as `./server`

Alternative (without `make`):

```bash
export $(grep -v '^#' .env | xargs)   # or set env vars manually
go run ./internal/cmd/main.go
```

---

## Running in Docker

Build and run using the provided `Dockerfile`:

```bash
docker build -t auth-service .
docker run --env-file .env -p 8080:8080 auth-service
```

Adjust `8080` and/or `SERVER_PORT` as required.

---

## Database Operations

Migrations are defined under `migrations/` and are executed via the `migrate` CLI wired through the `Makefile`.

### Run migrations

```bash
make migrate-up
```

### Roll back last migration

```bash
make migrate-down
```

To perform a full reset, run `make migrate-down` repeatedly until fully rolled back, then `make migrate-up` again.

---

## Health and Monitoring

There is currently no dedicated `/health` endpoint. To confirm basic service health:

- Verify the process/container is running.
- Send a request to a known route, for example:

```bash
curl -i http://localhost:8080/api/login
```

Database connectivity is logged at startup; failures to connect will cause the service to exit with an error.

---

## Logs

- **Local (`make run` / `go run`)**: logs are written to stdout.
- **Docker**: use:

```bash
docker logs <container_id>
```

Echo’s logger is configured with structured fields (time, method, URI, status, error).

---

## Common Issues

### 1. Database not reachable

- Ensure PostgreSQL is running and accessible from the container/host.
- Validate DB env vars (`DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PWD`, `DB_NAME`).
- Confirm migrations ran successfully: `make migrate-up`.

### 2. Port already in use

- Stop the conflicting process on `SERVER_PORT` (default `8080`), or
- Run on a different port:

```bash
SERVER_PORT=9090 make run
```

### 3. Invalid CSRF/JWT tokens

- Tokens expire; re-authenticate via the login endpoint (`POST /api/login`).
- Ensure cookies are enabled in the client.
- In development, confirm `ENV=development` in `.env` to relax cookie security where appropriate.

---

## Recovery

- If the service crashes locally, restart via:

```bash
make run
```

- If running in Docker, restart the container:

```bash
docker restart auth-service
```

- If DB schema or data is inconsistent, re-run migrations (`make migrate-down` then `make migrate-up`) and, if necessary, restore from a database backup.

---

## Deployment

Example registry-based deployment flow:

```bash
docker build -t registry.example.com/auth-service:latest .
docker push registry.example.com/auth-service:latest
```

Then deploy the image via your orchestration platform (Kubernetes, Docker Compose, etc.), wiring the same environment variables used in development.

---

## References

- `backend/README.md` – high-level overview and setup
- `backend/api_docs.md` – API documentation
- `backend/API.postman_collection.json` – Postman collection
