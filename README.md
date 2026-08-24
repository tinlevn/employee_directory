# Employee Directory

Local employee-directory startup prototype.

- `api/` — Go 1.27 + Fiber + pgx API
- `web/` — Astro + React islands frontend
- PostgreSQL 16+ or newer

## Local startup

1. Ensure PostgreSQL is running and create the `employee_directory` database.
2. Set `api/.env` from `api/.env.example`.
3. Apply migrations:

```powershell
psql -h localhost -U postgres -d employee_directory -f api/migrations/000001_init.up.sql
psql -h localhost -U postgres -d employee_directory -f api/migrations/000002_integrity.up.sql
```

4. Start the API:

```powershell
Set-Location api
go run ./cmd/server
```

5. Optional development data:

```powershell
go run ./cmd/seed
```

6. Start the frontend in a second terminal:

```powershell
Set-Location web
npm install
npm run dev
```

Open `http://localhost:4321`.

## Docker

With Docker Desktop installed, `docker compose up --build` starts PostgreSQL, applies migrations, and starts the API. The Astro frontend can still run locally with `npm run dev`.

Authentication is intentionally not included in this prototype yet. Do not expose the API to the public internet.
