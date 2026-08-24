# Employee Directory — Session Handoff

Date: 2026-08-24
Purpose: handoff document for the next coding session. Records what was built, why, and what remains unfinished.

---

## 1. High-Level Summary

Rebuilt the existing `.NET 10 + Angular` employee directory as a **Go + Fiber API** with an **Astro + React** frontend, backed by **PostgreSQL**.

The old stack lives in `staff-api/` (ASP.NET Core 10, EF Core, SQL Server LocalDB) and `staff-client/` (Angular 22, Material). These are **deprecated** and kept only for reference — the new work is in `api/` and `web/`.

Key outcome: a working local prototype with 1,003 seeded employees, search, pagination, hire-date sorting, SCD2 employment history, an append-only event log, analytics, and a lean static frontend.

---

## 2. Current Architecture

```
employee_directory/
├── api/                      Go 1.27 + Fiber API
│   ├── cmd/server/main.go    entrypoint, routes, CORS, health, OpenAPI
│   ├── cmd/seed/main.go      dev data generator (1000 random + 3 known)
│   ├── internal/
│   │   ├── config/           env loading (.env, .env.local)
│   │   ├── domain/           entities + enums (Person, EmploymentRecord, Event, etc.)
│   │   ├── dto/              request/response contracts + validation tags
│   │   ├── handler/          HTTP handlers (person, employment, event, emergency, org, analytics)
│   │   ├── middleware/       error handling, RFC7807, DB error mapping
│   │   └── repository/       pgx queries (CRUD, SCD2, analytics, integrity helpers)
│   ├── migrations/           000001_init, 000002_integrity (+ down files)
│   ├── openapi.yaml          checked-in OpenAPI 3.0 contract
│   ├── Dockerfile
│   └── go.mod
├── web/                      Astro 7 + React islands + Tailwind 4
│   ├── src/pages/            index.astro, person.astro, analytics.astro
│   ├── src/components/       Directory.tsx, PersonDetail.tsx, AnalyticsView.tsx
│   ├── src/lib/api.ts        typed API client
│   └── src/layouts/          Layout.astro
├── docker-compose.yml        postgres + migrate + api
├── .github/workflows/ci.yml  gofmt/vet/test/build + astro check/build
└── README.md                 root quick-start
```

### Stack

| Layer | Choice | Notes |
|-------|--------|-------|
| Language | Go 1.27 | installed via scoop |
| HTTP framework | Fiber v2.52.15 | high-throughput, Express-like |
| DB driver | pgx/v5 + pgxpool | raw SQL, no ORM |
| Validation | go-playground/validator/v10 | mirrors Pydantic tags |
| IDs | UUID (google/uuid) | `gen_random_uuid()` in Postgres |
| Database | PostgreSQL 18.6 | local install (not Docker for dev) |
| Frontend | Astro 7 | static output, zero JS by default |
| Islands | React 19 | only where interactivity is needed |
| Styling | Tailwind 4 | via `@tailwindcss/vite` |
| Migrations | golang-migrate (`migrate/migrate` image) | local via psql, Docker via compose |

---

## 3. Database Schema

Source of truth: `person_directory_schema_employee_only.md` (in repo parent), adapted and trimmed.

### Tables

- `organizations` — multi-tenant anchor
- `persons` — identity anchor (trimmed, see §5)
- `emergency_contacts` — separate entity, `UNIQUE(person_id)` (0..1 per person)
- `employment_records` — **SCD2** versioning (`valid_from` / `valid_to` / `is_current`)
- `status_change_events` — **append-only** event log
- `transfer_records` — lateral movement history
- `headcount_snapshots` — pre-aggregated analytics
- Supporting (created, not yet wired to endpoints): `person_skills`, `person_languages`, `person_documents`, `person_certifications`, `employee_benefits`, `person_accounts`

### Views

- `active_employees` — current headcount view
- `person_event_timeline` — event timeline with initiator/approver names

### Migrations

- `api/migrations/000001_init.up.sql` — full trimmed schema, triggers, views, default org seed
- `api/migrations/000002_integrity.up.sql` — relational integrity (see §7)

---

## 4. API Surface

All under `/api/v1`:

| Method | Route | Purpose |
|--------|-------|---------|
| GET/POST | `/organizations` | list/create |
| GET | `/organizations/{id}` | get org |
| GET | `/persons` | list, filter, paginate, sort |
| POST | `/persons` | create |
| GET | `/persons/{id}` | get (enriched with current employment) |
| PATCH | `/persons/{id}` | partial update |
| DELETE | `/persons/{id}?reason=` | soft delete |
| GET/POST/PATCH/DELETE | `/persons/{id}/emergency-contact` | emergency contact CRUD |
| GET | `/persons/{id}/employment` | all versions |
| GET | `/persons/{id}/employment/current` | current SCD2 record |
| POST | `/persons/{id}/employment` | new versioned record (closes previous) |
| GET/POST | `/persons/{id}/events` | immutable event log |
| GET/POST | `/persons/{id}/transfers` | transfer log |
| GET | `/analytics/headcount` | headcount by dept |
| GET | `/analytics/attrition` | attrition stats |
| GET | `/analytics/movements` | hires/exits over time |
| GET | `/analytics/snapshot/{date}?org_id=` | point-in-time snapshot |

Non-versioned routes: `GET /health`, `GET /openapi.yaml`, `GET /swagger` (redirect to OpenAPI).

### List query params (`/persons`)

`q` (name/email/job-title ILIKE search), `department`, `team`, `office_location`, `city`, `country`, `tag`, `org_id`, `is_active`, `page`, `page_size` (max 100), `sort`.

Sort whitelist: `last_name`, `first_name`, `created_at`, `city`, `hire_date` (with `-` prefix for DESC). **Default sort is latest hire first** (`hire_date DESC NULLS LAST`).

### Response envelope

- Paginated lists: `{ data, page, page_size, total, total_pages }`
- Errors: RFC 7807 style `{ type, title, status, detail, errors?, trace_id }`

---

## 5. Domain Trimming Decisions

**Removed from `persons`** (deemed unnecessary or sensitive):
- `ethnicity`, `religion`, `blood_type` (sensitive/GDPR)
- `maiden_name`, `pronouns`, `nationality` (rarely used)
- `phone_secondary`, `linkedin_url`, `personal_website_url`

**Removed from `employment_records`:**
- `division`, `business_unit`, `cost_center` (org bloat; easy to add back)

**Emergency contact** was split out of `persons` into its own table linked by FK, with `phone`, `email`, `relationship` (string) fields. Optional (0..1 per person).

**Kept in `persons`:** first/middle/last/preferred name, DOB, gender, profile photo, personal/org email, primary phone, full address block, `is_international`, `source`, `notes`, `tags`, soft-delete fields (`is_active`, `archived_at`, `archive_reason`), audit fields.

**Kept in `employment_records`:** employee_id, job_title, job_level, status/type/arrangement, department, team, office_location, desk_number, reports_to, salary fields, hire/probation/contract dates, SCD2 fields.

> Note: `salary_amount`/`salary_currency`/`hourly_rate` and personal contact details are still returned in directory/detail responses. There is no field-level access control yet — see §8.

---

## 6. Frontend Decisions

- **Astro static output** — directory page ships minimal HTML; React islands hydrate only where needed (`client:load`).
- **Routing** — employee detail uses query string `GET /person?id={uuid}` instead of `GET /persons/{uuid}` because static output can't pre-generate dynamic UUID paths. This was a deliberate fix for the production build (the old `[id].astro` dynamic route only generated a `preview` placeholder and would 404 in static hosting).
- **Directory table columns**: Name, Position (`job_title`), Department (badge), Email, City.
- **Zebra striping** — even rows white, odd rows `bg-zinc-50/70`, hover highlight.
- **URL state sync** — search/pagination/page-size are reflected in the query string via `history.pushState`, with `popstate` handling for back/forward and refresh persistence.
- **Race protection** — overlapping search requests are discarded via a `requestVersion` ref so stale responses can't overwrite newer ones.
- **Typed API client** — `web/src/lib/api.ts` declares `Person`, `EmploymentRecord`, `StatusChangeEvent`, `EmergencyContact`, `Paginated<T>`, etc.
- **No `innerHTML`** — detail and analytics pages were rewritten as React components to eliminate stored-XSS risk.

---

## 7. Hardening Pass Completed (This Session)

A staff-level code review identified release blockers. All **except authentication** were addressed:

1. **Relational integrity** (`000002_integrity.up.sql`):
   - `(person_id, org_id)` composite FK from `employment_records` → `persons` (same-org enforcement)
   - `reports_to` same-org FK
   - `status_change_events` and `transfer_records` same-org FKs
   - Unique `(org_id, lower(org_email))`
   - Unique current `(org_id, employee_id)`
   - App-level `ensurePersonInOrg` checks for `reports_to`, event actors (`initiated_by`/`approved_by`/`witnessed_by`), and transfer managers.
2. **Event immutability** — DB trigger `reject_event_mutation()` blocks `UPDATE`/`DELETE` on `status_change_events` (not just convention).
3. **XSS** — removed all `innerHTML` interpolation from Astro pages (now React JSX).
4. **Static routing** — `/person?id=` replaces broken dynamic route.
5. **Migrations** — Docker Compose now runs a `migrate` service before the API starts.
6. **Persisted timestamps** — create operations use `INSERT ... RETURNING` so `created_at`/`updated_at`/`recorded_at` are correct in responses.
7. **Error handling** — `RepositoryError()` maps `pgx.ErrNoRows` → 404, unique violation → 409, FK/check violation → 422, bad value → 400. Error `title` now matches the actual status code.
8. **Analytics** — implemented point-in-time `snapshot` (was `501 Not Implemented`); seeded events + snapshots so analytics returns real data.
9. **Seed data** — 1003 persons, 1003 current employment records, 1003 HIRED events, 291 emergency contacts, 84 headcount snapshots.
10. **Secrets** — `api/.env` is gitignored; seed reads `DATABASE_URL` from config instead of hardcoding.
11. **Tests** — `internal/dto/request_test.go` (pagination) and `internal/repository/integration_test.go` (person/employment SCD2, analytics; requires `TEST_DATABASE_URL`).
12. **CI** — `.github/workflows/ci.yml` (gofmt, vet, test, build for Go; check, build for web).
13. **OpenAPI** — checked-in `api/openapi.yaml`, served at `/openapi.yaml`.
14. **Formatting** — `gofmt` applied to all Go files.

### Verification (all passing)

- `go build ./...`
- `go vet ./...`
- `go test ./...` (with `TEST_DATABASE_URL` set for integration)
- `npm run check` (astro check: 0 errors/warnings/hints)
- `npm run build`
- `npm audit --omit=dev` — 0 vulnerabilities
- Fresh DB migration sequence (000001 then 000002) verified end-to-end

---

## 8. NOT DONE — Security (Critical)

Authentication and authorization were **explicitly deferred** and are the biggest remaining gap. The API must not be exposed to the public internet.

### What is missing

- **No authentication** — all endpoints are anonymous.
- **No authorization / RBAC** — no role checks despite `person_accounts.role` existing in the schema (unused).
- **No JWT / session / API-key middleware**.
- **Client-controlled tenancy** — `org_id` is taken from request bodies/query strings (`dto/request.go`). A caller can read/write arbitrary organizations. There is no token-derived org context.
- **No rate limiting**.
- **No CSRF protection** (irrelevant for token-auth APIs, but relevant if cookies are later used).
- **No field-level access control** — salary, personal email, phone, address, DOB, notes are returned to any caller.
- **No TLS termination** at the app (expected at proxy/load balancer).
- **`person_accounts`** table exists but has no login/registration endpoints and `password_hash` is unpopulated.

### What IS in place (defense-in-depth only)

- Parameterized SQL everywhere (no injection).
- Sort-field whitelist (no ORDER BY injection).
- 1 MB request body limit.
- CORS restricted to configured origins.
- `requestid`, logger, recover middleware.
- DB-level same-org FKs and event immutability trigger.

### Recommended auth work (next session)

1. Introduce JWT auth middleware; derive `org_id` from the authenticated principal, not the request body.
2. Add RBAC middleware using `person_accounts.role` (`admin`/`manager`/`staff`/`read-only`).
3. Split public directory DTOs from HR/admin DTOs (hide salary + private contact data from general listing).
4. Add rate limiting and request logging with trace IDs.
5. Add login/registration endpoints + password hashing (bcrypt/argon2).
6. Re-evaluate CORS for production origins.

---

## 9. How to Run

### Prerequisites

- Go 1.27
- Node 22+
- PostgreSQL running locally (database `employee_directory`)

### API

```powershell
cd api
cp .env.example .env          # set DATABASE_URL (see .env.example)
go run ./cmd/server           # :8080
```

Apply migrations (fresh DB):

```powershell
psql -h localhost -U postgres -d employee_directory -f api/migrations/000001_init.up.sql
psql -h localhost -U postgres -d employee_directory -f api/migrations/000002_integrity.up.sql
```

Seed dev data:

```powershell
cd api
go run ./cmd/seed
```

### Frontend

```powershell
cd web
npm install
npm run dev                   # :4321
```

### Docker (full stack)

```powershell
docker compose up --build
```

---

## 10. Known Limitations / Pitfalls

- **No auth** (see §8) — biggest blocker to any real deployment.
- **No Go unit tests for handlers/middleware** — only repository + DTO tests exist. Handler-level HTTP tests still needed.
- **Legacy compat routes** (`/api/employees`) are a placeholder and do **not** preserve the old Angular response shape; they should be removed once the old client is gone.
- **`person_accounts`, skills, languages, documents, certifications, benefits** tables exist but have no API endpoints.
- **Attrition** depends on `status_change_events` (RESIGNED/TERMINATED/LAID_OFF) and `headcount_snapshots`; seed data only generates HIRED events, so attrition is empty until those event types are logged.
- **Salary** is stored as minor units (cents) by convention, but there's no currency consistency enforcement.
- **Date fields** serialize as full ISO timestamps (`2024-12-24T00:00:00Z`) rather than `YYYY-MM-DD`; a date wrapper may be desirable.
- **Startup DB reconnect** — the API connects once at startup and stays degraded if Postgres is unavailable; no retry loop.
- **Old .NET/Angular projects** (`staff-api/`, `staff-client/`, `tests/`) remain on disk and are deprecated; delete them once migration is complete.

---

## 11. Suggested Next Steps (Priority Order)

1. Authentication + tenant derivation (P0).
2. RBAC + field-level DTO separation (P0).
3. Handler-level HTTP tests (P1).
4. Remove legacy `/api/employees` compat routes.
5. Wire up supporting tables (skills/documents/certifications/benefits) to endpoints.
6. Add password hashing + login/registration.
7. Delete deprecated `staff-api`/`staff-client`.
