# Auth-Service Backend

Go-based REST API for user authentication and address management.  
Built with **Echo**, structured in layered architecture (handler, service, repository, model, validator), and designed for containerized deployment.

---

## Overview

- **Domain**: Auth + user address book
- **Framework**: Echo (`github.com/labstack/echo/v4`)
- **Persistence**: PostgreSQL (via `pgx`)
- **Security**:
  - JWT access tokens stored in **HTTP-only cookies**
  - CSRF protection using a separate CSRF token cookie + header
  - Environment-aware cookie settings (secure/same-site based on `ENV`)
- **Docs & tooling**:
  - `api_docs.md` – human-readable API documentation
  - `API.postman_collection.json` – Postman collection
  - `RUNBOOK.md` – operational procedures

---

## Features

- **Authentication**
  - User registration and login
  - JWT-based session via `access_token` cookie
  - Logout endpoint that invalidates client session
- **Address management**
  - CRUD for user addresses
  - User-scoped access to address records
- **Validation & error handling**
  - Centralized request validation using `go-playground/validator`
  - Standardized error response envelope from handler layer
- **Database migrations**
  - SQL migrations under `migrations/`
  - Make targets for up/down/reset flows
- **Containerization**
  - Production-oriented `Dockerfile`
  - Make target `docker-build` for image creation

---

## Project Structure

From `backend/`:

```text
.
├── api_docs.md                   # API documentation
├── API.postman_collection.json   # Postman collection
├── Dockerfile                    # Docker image build
├── go.mod / go.sum               # Go module definition and dependencies
├── internal                      # Application code
│   ├── cmd/main.go               # Service entry point (Echo server)
│   ├── config/config.go          # Configuration loading (env-based)
│   ├── cookie/cookie.go          # Cookie helpers (SameSite / Secure)
│   ├── db/db.go                  # Database connection setup
│   ├── handler/                  # HTTP handlers (auth, address, errors)
│   ├── model/                    # Domain models (user, address)
│   ├── repo/                     # Repositories (DB access)
│   ├── response/                 # Response + error envelope helpers
│   ├── service/                  # Business logic
│   └── validator/validator.go    # Request validation integration
├── migrations                    # SQL migration files
├── Makefile                      # Build, run, migration, Docker helpers
├── LICENSE                       # License file
└── server/                       # Binary output (via Makefile `server`)
```

---

## Prerequisites

- Go 1.21+ (Go module currently targets Go 1.24)
- PostgreSQL instance (local or remote)
- `migrate` CLI installed (for DB migrations)
- Docker (optional, for containerized runs)

---

## Configuration

Configuration is loaded from environment variables (or `.env` when using `make run` / migration targets).

Key variables (see `internal/config/config.go`):

- `DB_HOST` (default: `localhost`)
- `DB_PORT` (default: `5432`)
- `DB_USER` (required)
- `DB_PWD` (required)
- `DB_NAME` (required)
- `ENV` (default: `development`)
- `JWT_SECRET` (required)
- `SESSION_KEY` (required)
- `SERVER_HOST` (default: `0.0.0.0`)
- `SERVER_PORT` (default: `8080`)

Example `.env` (for local development):

```bash
DB_HOST=localhost
DB_PORT=5432
DB_USER=auth_user
DB_PWD=change-me
DB_NAME=auth_db
ENV=development
JWT_SECRET=local-jwt-secret
SESSION_KEY=local-session-key
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
```

---

## Local Development

All commands assume you are in the `backend/` directory.

### Install dependencies

```bash
go mod tidy
```

### Build the binary

```bash
make compile
```

### Run database migrations

```bash
make migrate-up
```

### Start the server (with `.env`)

```bash
make run
```

The API listens on `http://SERVER_HOST:SERVER_PORT` as defined in your environment (by default `http://localhost:8080`).

---

## Running with Docker

From `backend/`:

```bash
docker build -t auth-service .
docker run --env-file .env -p 8080:8080 auth-service
```

Adjust the port mapping if you change `SERVER_PORT` in your environment.

---

## API Surface

See `api_docs.md` and the Postman collection for full details.

Core routes (from `internal/cmd/main.go`):

- **Auth**
  - `POST /api/register`
  - `POST /api/login`
  - `POST /api/v1/logout`
  - `GET  /api/v1/profile`
- **Addresses**
  - `GET    /api/v1/users/addresses`
  - `POST   /api/v1/users/address/add`
  - `GET    /api/v1/users/address/:id`
  - `PATCH  /api/v1/users/address/:id`
  - `DELETE /api/v1/users/address/:id`

Protected routes under `/api/v1` require:

- Valid JWT in `access_token` HTTP-only cookie
- Valid CSRF token provided via `X-CSRF-Token` header

---

## Development Notes

- Follows handler → service → repository layering for clear separation of concerns.
- Echo middleware is used for:
  - Structured request logging
  - JWT auth (`echo-jwt`)
  - CSRF protection
  - CORS (configured for the React frontend at `http://localhost:5173`).
- Error responses are standardized through the `internal/response` package.

---

## Documentation

- **[DECISIONS.md](./DECISIONS.md)** – Architectural decisions and technical reasoning
- **[RUNBOOK.md](./RUNBOOK.md)** – Setup, operations, and deployment instructions
- **[API Collection](./API.postman_collection.json)** – Postman collection for testing the API

## Quick Links

- [Why these architectural choices?](./DECISIONS.md)
- [How to run this project?](./RUNBOOK.md)

---

## License

MIT License. See `LICENSE` for details.
