# Auth-Service Backend

## Live Demo
🚀 **Backend API**: https://auth-service-backend-go-production.up.railway.app
🌐 **Frontend App**: https://auth-frontend-react-ts.vercel.app

<<<<<<< HEAD
A Go-based RESTful API that provides user authentication and address management. Built with **Echo**, structured into clean layers (handler, service, repo, model, validator), and containerized with Docker. Includes database migrations and Postman collection for testing.

---

## Related Repositories

- **Frontend Client**: [auth-frontend-react-ts](https://github.com/nawwardiab/auth-frontend-react-ts)
=======
Go-based REST API for user authentication and address management with production deployment on Railway.
>>>>>>> 4724c0c991b7d2c9dda6ecf54b0be81fd1dda809

---

## Features

<<<<<<< HEAD
- **Auth**: Register, login (JWT cookies), logout.
- **Address Management**: Create, read, update, delete addresses linked to users.
- **Database Migrations**: Versioned SQL migrations.
- **Validation**: Input validation for all major endpoints.
- **Dockerized**: Production-ready Dockerfile.
- **Postman Collection**: Predefined requests for easy testing.
=======
- **JWT Authentication** - Secure token-based auth with HTTP-only cookies
- **CSRF Protection** - Double-submit cookie pattern for cross-site request forgery prevention
- **Bcrypt Password Hashing** - Industry-standard password security
- **PostgreSQL Database** - Reliable persistence with Railway-hosted PostgreSQL
- **Address CRUD** - Complete address management with user-scoped access control
- **Clean 4-Layer Architecture** - Handlers → Services → Repositories → Models
- **Docker Multi-Stage Builds** - Optimized 5MB production image using Alpine Linux
- **Standardized Error Handling** - Consistent error response envelope across all endpoints

---

## Tech Stack

- **Go 1.23** - Modern Go with latest features
- **Echo v4** - High-performance web framework
- **PostgreSQL** - Production database with pgx driver
- **JWT + CSRF** - Secure authentication and session management
- **Docker** - Containerized deployment
- **Railway** - Cloud platform for backend and database hosting

---

## API Endpoints

### Public Endpoints
- `POST /api/register` - Register new user
- `POST /api/login` - Login and receive JWT cookie

### Protected Endpoints (require JWT + CSRF token)
- `POST /api/v1/logout` - Logout and clear session
- `GET  /api/v1/profile` - Get current user profile
- `GET  /api/v1/users/addresses` - List all user addresses
- `POST /api/v1/users/address/add` - Create new address
- `GET  /api/v1/users/address/:id` - Get specific address
- `PATCH /api/v1/users/address/:id` - Update address
- `DELETE /api/v1/users/address/:id` - Delete address

---

## Architecture

The backend follows a clean 4-layer architecture pattern:

```
Handlers (HTTP layer)
    ↓
Services (Business logic)
    ↓
Repositories (Data access)
    ↓
Models (Domain entities)
```

**Benefits:**
- Clear separation of concerns
- Easy to test each layer independently
- Maintainable and scalable codebase
- Framework-agnostic business logic

---

## Security

- **HTTP-Only Cookies** - JWT tokens stored in HTTP-only cookies to prevent XSS attacks
- **SameSite=None** - Configured for secure cross-site requests (production)
- **Bcrypt Password Hashing** - All passwords hashed with bcrypt (cost factor 10)
- **CSRF Token Validation** - All protected routes validate CSRF tokens via `X-CSRF-Token` header
- **Environment-Based Configuration** - Secure/SameSite settings adjust based on `ENV` variable

---

## Quick Test

Test the registration endpoint with curl:

```bash
curl -X POST https://auth-service-backend-go-production.up.railway.app/api/register \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "John",
    "last_name": "Doe",
    "email": "john.doe@example.com",
    "password": "SecurePass123!"
  }'
```

Expected response:
```json
{
  "id": "uuid-here",
  "first_name": "John",
  "last_name": "Doe",
  "email": "john.doe@example.com",
  "created_at": "2025-12-15T10:00:00Z"
}
```

---

## Deployment

**Platform**: Railway
**Deployment Date**: December 15, 2025
**Region**: EU West
**Features**:
- Automatic HTTPS via Railway's edge network
- Docker containerization with multi-stage builds
- Connected to Railway PostgreSQL database
- Environment variables managed through Railway dashboard
- Auto-deploy on git push (optional)

**Production URL**: https://auth-service-backend-go-production.up.railway.app
>>>>>>> 4724c0c991b7d2c9dda6ecf54b0be81fd1dda809

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

## Local Development

<<<<<<< HEAD
- Go 1.21+
- PostgreSQL (or configured DB)
- Docker & Docker Compose (optional)
=======
All commands assume you are in the `backend/` directory.
>>>>>>> 4724c0c991b7d2c9dda6ecf54b0be81fd1dda809

### Prerequisites

- Go 1.21+ (Go module currently targets Go 1.24)
- PostgreSQL instance (local or remote)
- `migrate` CLI installed (for DB migrations)
- Docker (optional, for containerized runs)

### Configuration

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

<<<<<<< HEAD
- Markdown doc: [api_docs.md](./api_docs.md)
- Postman collection: [API.postman_collection.json](./API.postman_collection.json)
=======
From `backend/`:
>>>>>>> 4724c0c991b7d2c9dda6ecf54b0be81fd1dda809

```bash
docker build -t auth-service .
docker run --env-file .env -p 8080:8080 auth-service
```

<<<<<<< HEAD
- **Auth**: `/register`, `/login`, `/api/v1/logout`
- **Address**: `/api/v1/users/address/add`, `/api/v1/users/address/{id}`
=======
Adjust the port mapping if you change `SERVER_PORT` in your environment.
>>>>>>> 4724c0c991b7d2c9dda6ecf54b0be81fd1dda809

---

## Error Response Format

All API errors return a standardized JSON envelope:
```json
{
  "error": {
    "code": "AUTH_INVALID_CREDENTIALS",
    "message": "Invalid email or password"
  }
}
```

Error codes include:
- Authentication: `AUTH_INVALID_CREDENTIALS`, `AUTH_USER_EXISTS`, `AUTH_MISSING_TOKEN`
- Address: `ADDRESS_NOT_FOUND`, `ADDRESS_FORBIDDEN`, `ADDRESS_CANNOT_DELETE_DEFAULT`
- Validation: `VALIDATION_ERROR`, `INVALID_PAYLOAD`

See `internal/response/error.go` for the complete list of error codes.

---

## Documentation

- **[DECISIONS.md](./DECISIONS.md)** - Architectural decisions and technical reasoning
- **[RUNBOOK.md](./RUNBOOK.md)** - Setup, operations, and deployment instructions
- **[api_docs.md](./api_docs.md)** - Human-readable API documentation
- **[API Collection](./API.postman_collection.json)** - Postman collection for testing the API

---

## License

MIT License. See `LICENSE` for details.

---

**Built by Nawar Diab** | [GitHub](https://github.com/nawwardiab) | [LinkedIn](https://www.linkedin.com/in/nawar-diab)
