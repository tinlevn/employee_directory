# Staff Search API

Backend for the Employee Directory web app — an ASP.NET Core **10** Web API with Entity Framework Core 10 (SQL Server).

## Architecture

```
Controllers/            -> HTTP layer (EmployeesController: REST endpoints only)
Dtos/                   -> API contracts (EmployeeDto out, EmployeeQuery in + validation)
Entities/               -> Persistence model (Employee, maps to dbo.Employees)
Data/                   -> StaffSearchDbContext + Repositories (EF Core data access)
Interfaces/             -> IEmployeeRepository abstraction (DI + testability)
Extensions/             -> Entity -> DTO mapping
Migrations/             -> EF Core code-first migrations
```

Requests flow: `EmployeesController -> IEmployeeRepository -> StaffSearchDbContext -> SQL Server`.
The controller never sees entities; the repository never sees HTTP.

## Endpoints

| Method | Route | Description |
| ------ | ----- | ----------- |
| GET | `/api/employees` | All employees; optional query-string filters `firstName`, `lastName`, `jobTitle`, `department`, `location` (case-insensitive contains for text, exact for dropdowns) |
| GET | `/api/employees/{id}` | Single employee or `404` |
| GET | `/health` | Health check (includes database connectivity) |
| GET | `/swagger` | OpenAPI UI (Development only) |

Errors are returned as RFC 7807 `ProblemDetails`.

## Run locally

```powershell
dotnet restore
dotnet tool restore                 # restores the pinned dotnet-ef tool
dotnet ef database update           # creates/updates the LocalDB agu_staff database
dotnet run                          # http://localhost:51828 / https://localhost:51827
```

CORS origins are configured in `appsettings.json` under `Cors:AllowedOrigins` (Angular dev server: `http://localhost:4200`).

## Frontend

The Angular client lives in [`../staff-client`](../staff-client) (standalone components, signals, zoneless change detection).
