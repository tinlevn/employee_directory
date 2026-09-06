# Plan 001: Enforce Multi-Tenant Isolation and Fix Auth Claims Parsing

> **Executor instructions**: Follow this plan step by step. Run every verification command and confirm the expected result before moving to the next step. If anything in the "STOP conditions" section occurs, stop and report — do not improvise. When done, update the status row for this plan in `plans/README.md`.
>
> **Drift check (run first)**: `git diff --stat e2eb40e..HEAD -- api/internal/handler/person_handler.go api/internal/handler/emergency_handler.go api/internal/handler/employment_handler.go api/internal/handler/event_handler.go api/internal/repository/person_repo.go api/internal/middleware/auth.go`
> If any in-scope file changed since this plan was written, compare the "Current state" excerpts against the live code before proceeding; on a mismatch, treat it as a STOP condition.

## Status

- **Priority**: P0
- **Effort**: M
- **Risk**: LOW
- **Depends on**: none
- **Category**: security
- **Planned at**: commit `e2eb40e`, 2026-09-06

## Why this matters

The API enforces tenant isolation on `List` and `Create`, but every handler querying a sub-resource by person ID (`GET /persons/:id`, `PATCH /persons/:id`, `DELETE /persons/:id`, `/persons/:id/emergency-contact`, `/persons/:id/employment`, and `/persons/:id/events`) only passes the unvalidated URL parameter `:id`. An authenticated user belonging to Tenant A can read, alter, or soft-delete employee profiles and private logs belonging to Tenant B simply by guessing or knowing the UUID (IDOR). Additionally, `middleware.GetAccountID` uses `uuid.MustParse` on the JWT subject claim, creating a panic vector on malformed claims. This plan binds all single-entity operations to `middleware.GetOrgID(c)`.

## Current state

- In `api/internal/handler/person_handler.go:46-59`, `Get` executes:
  ```go
  id, err := uuid.Parse(c.Params("id"))
  ...
  p, err := h.repo.GetByID(c.Context(), id)
  ```
  And in `api/internal/repository/person_repo.go:88`:
  ```sql
  WHERE p.id=$1
  ```
  There is no check that `p.OrgID == middleware.GetOrgID(c)`.
- In `api/internal/handler/emergency_handler.go:28,60`, `employment_handler.go:32,47,60`, and `event_handler.go:41,115`, handlers fetch records by `person_id` without confirming the person belongs to the caller's organization.
- In `api/internal/middleware/auth.go:67`:
  ```go
  func GetAccountID(c *fiber.Ctx) uuid.UUID {
      if cl := claimsFrom(c); cl != nil {
          return uuid.MustParse(cl.Subject)
      }
      return uuid.Nil
  }
  ```
  `uuid.MustParse` panics if `cl.Subject` is not a valid UUID string.

## Commands you will need

| Purpose | Command | Expected on success |
|---------|---------|---------------------|
| Build API | `go build ./...` (in `api/`) | exit 0, no errors |
| Vet API | `go vet ./...` (in `api/`) | exit 0 |
| Run Tests | `go test ./...` (in `api/`) | exit 0 |
| Format Check | `gofmt -l .` (in `api/`) | empty output |

## Scope

**In scope**:
- `api/internal/middleware/auth.go`
- `api/internal/repository/person_repo.go`
- `api/internal/repository/emergency_repo.go`
- `api/internal/repository/employment_repo.go`
- `api/internal/repository/event_repo.go`
- `api/internal/handler/person_handler.go`
- `api/internal/handler/emergency_handler.go`
- `api/internal/handler/employment_handler.go`
- `api/internal/handler/event_handler.go`

**Out of scope**:
- Database schema changes (migrations). The `org_id` column already exists on all relevant tables (`persons`, `emergency_contacts` via `person_id`, `employment_records`, `status_change_events`).
- Frontend components.

## Git workflow

- Branch: `advisor/001-tenant-isolation-idor`
- Commit style: Conventional Commits, e.g. `fix(api): enforce org_id tenant isolation across person and sub-resource handlers`

## Steps

### Step 1: Safely parse `Subject` in `middleware.GetAccountID`
In `api/internal/middleware/auth.go`, change `uuid.MustParse(cl.Subject)` to use `uuid.Parse`.
```go
func GetAccountID(c *fiber.Ctx) uuid.UUID {
	if cl := claimsFrom(c); cl != nil {
		if id, err := uuid.Parse(cl.Subject); err == nil {
			return id
		}
	}
	return uuid.Nil
}
```

**Verify**: Run `go test ./...` in `api/` → exit 0.

### Step 2: Add `org_id` scoped repository methods for Person
In `api/internal/repository/person_repo.go`:
1. Update `GetByID(ctx context.Context, id, orgID uuid.UUID) (*domain.Person, error)` to append `AND p.org_id = $2`.
2. Update `Update(ctx context.Context, id, orgID uuid.UUID, fields map[string]any) (*domain.Person, error)` to include `AND org_id = $%d` in the update query.
3. Update `SoftDelete(ctx context.Context, id, orgID uuid.UUID, reason string) (bool, error)` to include `AND org_id = $3`.

**Verify**: Run `go build ./...` in `api/` (will show compiler errors for callers until Step 3).

### Step 3: Update `PersonHandler` to pass authenticated OrgID
In `api/internal/handler/person_handler.go`:
1. In `Get(c *fiber.Ctx)`:
   Pass `middleware.GetOrgID(c)` to `h.repo.GetByID(c.Context(), id, orgID)`.
2. In `Update(c *fiber.Ctx)`:
   Pass `orgID` to `h.repo.Update(c.Context(), id, orgID, fields)`.
3. In `Delete(c *fiber.Ctx)`:
   Pass `orgID` to `h.repo.SoftDelete(c.Context(), id, orgID, reason)`.

**Verify**: Run `go build ./...` in `api/` → compiles cleanly for person handler.

### Step 4: Add tenant verification helper for sub-resource handlers
In `api/internal/handler/emergency_handler.go`, `employment_handler.go`, and `event_handler.go`:
Before querying or modifying sub-resources by person ID (`c.Params("id")`), verify that the target person belongs to the caller's tenant:
```go
orgID := middleware.GetOrgID(c)
person, err := h.personRepo.GetByID(c.Context(), pid, orgID)
if err != nil {
    return middleware.RepositoryError(err)
}
if person == nil {
    return fiber.NewError(fiber.StatusNotFound, "person not found")
}
```
Inject `*repository.PersonRepository` into `NewEmergencyHandler`, `NewEmploymentHandler`, and `NewEventHandler` (and wire in `api/cmd/server/main.go`).

**Verify**:
Run `go build ./...` in `api/` → exit 0.
Run `gofmt -l .` in `api/` → no unformatted files.

## Test plan

- Test that calling `GetByID` with a valid person ID from a different `org_id` returns `nil, nil` (yielding HTTP 404).
- Test that calling `SoftDelete` on a person with mismatched `org_id` returns `false` (no rows affected).
- Test `GetAccountID` with malformed subject (e.g. `"not-a-uuid"`) does not panic and returns `uuid.Nil`.

## Done criteria

- [ ] `go build ./...` in `api/` exits 0.
- [ ] `go vet ./...` in `api/` exits 0.
- [ ] `go test ./...` in `api/` exits 0.
- [ ] No handler can retrieve or mutate a person record or associated sub-resources without validating `org_id`.
- [ ] `plans/README.md` row updated to `DONE`.

## STOP conditions

- If `integration_test.go` fails due to changed method signatures, update `integration_test.go` helper calls to pass the test `orgID`. Do not remove test assertions.

## Maintenance notes

- Any future endpoints added under `/persons/:id/*` must validate that the parent person belongs to `middleware.GetOrgID(c)`.

