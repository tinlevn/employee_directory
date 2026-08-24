# Employee Directory API — Go + Fiber

Robust replacement for the old `staff-api` (.NET) implementing the full schema from `person_directory_schema_employee_only.md`.

## Stack
- **Fiber v2** — 3x faster than net/http, Express-like
- **pgx v5** — native Postgres driver + `pgxpool`
- **validator v10** — request validation (mirrors Pydantic)
- **golang-migrate** — SQL migrations (see `migrations/`)
- Postgres 16+ with `pgcrypto` for `gen_random_uuid()`

## Schema principles
- Immutable event log (`status_change_events` append-only)
- SCD2 versioning (`employment_records.valid_from/valid_to/is_current` — only one current per person)
- Soft deletes (`persons.is_active`, `archived_at`)
- Multi-tenant (`org_id` everywhere)
- Audit (`created_at/updated_at/created_by`)

## Run

```powershell
cp .env.example .env
# edit DATABASE_URL
go run ./cmd/server              # :8080
# or
go build -o bin/server ./cmd/server && ./bin/server
```

Migrate manually:

```powershell
migrate -path migrations -database $env:DATABASE_URL up
# or psql -f migrations/000001_init.up.sql
```

Docker (when you install Docker Desktop):

```powershell
docker compose up --build
```

The Compose `migrate` service applies every migration before the API starts. For a local PostgreSQL install, apply migrations explicitly with `psql -f migrations/000001_init.up.sql` followed by `psql -f migrations/000002_integrity.up.sql`.

## Endpoints

| Method | Route | Description |
|--------|-------|-------------|
| GET | `/health` | DB-aware health |
| GET | `/api/v1/organizations` | List orgs |
| POST | `/api/v1/organizations` | Create org |
| GET | `/api/v1/persons?q=&department=&org_id=&page=&page_size=&sort=` | List persons (paginated, searchable) |
| POST | `/api/v1/persons` | Create person |
| GET | `/api/v1/persons/{id}` | Get person |
| PATCH | `/api/v1/persons/{id}` | Update person |
| DELETE | `/api/v1/persons/{id}?reason=` | Soft-delete |
| GET | `/api/v1/persons/{id}/employment` | All versions |
| GET | `/api/v1/persons/{id}/employment/current` | Current SCD2 |
| POST | `/api/v1/persons/{id}/employment` | New versioned record (closes previous) |
| GET | `/api/v1/persons/{id}/events?event_type=&from_date=&to_date=` | Timeline (paginated) |
| POST | `/api/v1/persons/{id}/events` | Append event (immutable) |
| GET | `/api/v1/persons/{id}/transfers` | Transfer history |
| POST | `/api/v1/persons/{id}/transfers` | Log transfer |
| GET | `/api/v1/analytics/headcount?org_id=` | Headcount by dept |
| GET | `/api/v1/analytics/attrition?from_date=&to_date=&org_id=` | Attrition |
| GET | `/api/v1/analytics/movements?from_date=&to_date=&org_id=` | Hires/exits over time |
| GET | `/api/v1/analytics/snapshot/{date}?org_id=` | Point-in-time headcount snapshot |

Errors are RFC 7807 `ProblemDetails` compatible. Validation 400s include `errors` map.

## Legacy compat
The old `/api/employees` routes remain only as a temporary compatibility shim. They do not preserve the old Angular response shape; migrate clients to `/api/v1/persons`.

## Old .NET project
`../staff-api` is kept for reference but is **deprecated**. Delete it once Astro client points to `/api/v1`.
