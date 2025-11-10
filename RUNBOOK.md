## Purpose

This runbook describes operational procedures for running, maintaining, and troubleshooting the **auth and address API** service.

---

## Service Overview

* **Language**: Go (Echo framework)
* **Database**: PostgreSQL
* **Containerized**: Docker image available
* **Port**: Default `8080`

---

## Starting the Service

### Local (dev)

```bash
go run ./internal/cmd/main.go
```

### Docker (prod/dev)

```bash
docker build -t auth-api .
docker run -p 8080:8080 auth-api
```

---

## Database Operations

### Run migrations

```bash
make migrate-up
```

### Rollback last migration

```bash
make migrate-down
```

### Reset DB

```bash
make migrate-reset
```

---

## Health Checks

* API root: `curl http://localhost:8080/health`
* Database connectivity is logged at startup.

---

## Logs

* Local: stdout from `go run`.
* Docker: `docker logs <container_id>`.

---

## Common Issues

### 1. Database not reachable

* Ensure PostgreSQL is running.
* Check credentials in `internal/config/config.go`.
* Run `make migrate-up` before first start.

### 2. Port already in use

* Stop the conflicting process.
* Or run with custom port:

  ```bash
  PORT=9090 go run ./internal/cmd/main.go
  ```

### 3. Invalid CSRF/JWT tokens

* Tokens expire. Re-login via `/login`.
* Ensure cookies are enabled in client.

---

## Recovery

* If service crashes, restart container:

  ```bash
  docker restart auth-api
  ```
* If DB schema mismatch, run migrations again.

---

## Deployment

* Build and push Docker image to registry:

  ```bash
  docker build -t registry.example.com/auth-api:latest .
  docker push registry.example.com/auth-api:latest
  ```
* Deploy with orchestration (Kubernetes, Docker Compose).

---

## Contacts

* **Maintainer**: DevOps / Backend team
* **Docs**: [api_docs.md](./api_docs.md)
