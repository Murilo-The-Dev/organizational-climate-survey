# Backend - Clima Organizacional

Go + mux + PostgreSQL REST API for the organizational climate survey system.

## Quick Start

### With Docker (recommended)
```bash
cp .env.example .env
# Edit .env with your values
docker compose up -d
```

API available at http://localhost:8080
Health check: http://localhost:8080/health
Docs static files: http://localhost:8080/docs/
Swagger UI: http://localhost:8080/swagger/index.html
OpenAPI JSON: http://localhost:8080/swagger/doc.json

## API Reference

- Frontend guide: `backend/docs/API_GUIDE.md`
- Generated OpenAPI files:
	- `backend/docs/swagger.yaml`
	- `backend/docs/swagger.json`
	- `backend/docs/docs.go`

### Development (hot reload)
```bash
docker compose -f docker-compose.dev.yml up -d
```

### Manual
```bash
cd backend
go mod download
go run cmd/api/main.go
```

## Tests
```bash
cd backend
go test ./... -count=1
go test ./internal/domain/usecase/... -cover
```

## CI/CD
- CI runs on every push to main/develop with backend changes: vet, tests, coverage check, Docker build.
- CD deploys to production on push to main via registry + SSH.
- Security scan runs weekly with govulncheck.

## Environment Variables
See .env.example at project root.

## Migrations
SQL migration files are in backend/migrations/.
They are mounted into the database container at /docker-entrypoint-initdb.d when using Docker Compose.