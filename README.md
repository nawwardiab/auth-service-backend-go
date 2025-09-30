# Auth-Service Backend API

A Go-based RESTful API that provides user authentication and address management. Built with **Echo**, structured into clean layers (handler, service, repo, model, validator), and containerized with Docker. Includes database migrations and Postman collection for testing.

---

## Related Repositories
- **Frontend Client**: [auth-frontend-react-ts](https://github.com/nawwardiab/auth-frontend-react-ts)

---

## Features

* **Auth**: Register, login (JWT cookies), logout.
* **Address Management**: Create, read, update, delete addresses linked to users.
* **Database Migrations**: Versioned SQL migrations.
* **Validation**: Input validation for all major endpoints.
* **Dockerized**: Production-ready Dockerfile.
* **Postman Collection**: Predefined requests for easy testing.

---

## Project Structure

```
.
├── api_docs.md                   # API documentation
├── API.postman_collection.json   # Postman collection
├── Dockerfile                    # Docker build file
├── go.mod / go.sum               # Go modules
├── internal                      # Application code
│   ├── cmd/main.go               # Entry point
│   ├── config/config.go          # Config management
│   ├── db/db.go                  # DB connection setup
│   ├── handler/                  # HTTP handlers
│   ├── model/                    # Domain models
│   ├── repo/                     # Database repositories
│   ├── service/                  # Business logic
│   └── validator/validator.go    # Input validation
├── migrations                    # SQL migrations
├── Makefile                      # Development tasks
├── LICENSE                       # License file
└── server/                       # (Reserved for server setup)
```

---

## Requirements

* Go 1.21+
* PostgreSQL (or configured DB)
* Docker & Docker Compose (optional)

---

## Setup

### 1. Clone the repo

```bash
git clone https://github.com/nawwardiab/auth-service-backend-go.git
cd auth-service-backend-go
```

### 2. Run migrations

```bash
make migrate-up
```

### 3. Start server (local)

```bash
go run ./internal/cmd/main.go
```

Server will start on `http://localhost:8080`.

### 4. Run with Docker

```bash
docker build -t auth-api .
docker run -p 8080:8080 auth-api
```

---

## API Documentation

* Markdown doc: [api_docs.md](./api_docs.md)
* Postman collection: [API.postman_collection.json](./API.postman_collection.json)

### Key Endpoints

* **Auth**: `/register`, `/login`, `/api/v1/logout`
* **Address**: `/api/v1/users/address/add`, `/api/v1/users/address/{id}`

---

## Development

Common Makefile commands:

```bash
make run           # Run server
make test          # Run tests
make migrate-up    # Apply migrations
make migrate-down  # Rollback migrations
```

---

## License

MIT License. See [LICENSE](./LICENSE).
