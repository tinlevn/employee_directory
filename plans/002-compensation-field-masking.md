# Plan 002: Role-Based Compensation and Private Detail Field Masking

> **Executor instructions**: Follow this plan step by step. Run every verification command and confirm the expected result before moving to the next step. If anything in the "STOP conditions" section occurs, stop and report — do not improvise. When done, update the status row for this plan in `plans/README.md`.
>
> **Drift check (run first)**: `git diff --stat e2eb40e..HEAD -- api/internal/domain/employment.go api/internal/handler/employment_handler.go api/internal/handler/person_handler.go`
> If any in-scope file changed since this plan was written, compare the "Current state" excerpts against the live code before proceeding; on a mismatch, treat it as a STOP condition.

## Status

- **Priority**: P0
- **Effort**: M
- **Risk**: LOW
- **Depends on**: plans/001-tenant-isolation-idor.md
- **Category**: security
- **Planned at**: commit `e2eb40e`, 2026-09-06

## Why this matters

`EmploymentRecord` currently includes sensitive financial compensation data: `SalaryAmount`, `SalaryCurrency`, `PayFrequency`, and `HourlyRate`. In `EmploymentHandler.List` and `EmploymentHandler.GetCurrent`, raw records are returned directly to any authenticated caller. Any user with role `staff` or `read-only` can inspect executive and peer salaries across the company. This plan introduces field-level access control: only users with `admin` or `manager` roles, or an employee inspecting their own record (`callerPersonID == targetPersonID`), may view compensation data; for others, salary fields are omitted.

## Current state

- In `api/internal/domain/employment.go:24-27`:
  ```go
  SalaryAmount   *int64   `json:"salary_amount,omitempty"`
  SalaryCurrency *string  `json:"salary_currency,omitempty"`
  PayFrequency   *string  `json:"pay_frequency,omitempty"`
  HourlyRate     *float64 `json:"hourly_rate,omitempty"`
  ```
- In `api/internal/handler/employment_handler.go:27-56`, `List` and `GetCurrent` serialize `domain.EmploymentRecord` directly without evaluating `middleware.GetRole(c)` or `middleware.GetPersonID(c)`.
- In `api/internal/middleware/auth.go:72-84`, helper functions `GetRole(c)` and `GetPersonID(c)` are already defined.

## Commands you will need

| Purpose | Command | Expected on success |
|---------|---------|---------------------|
| Build API | `go build ./...` (in `api/`) | exit 0, no errors |
| Vet API | `go vet ./...` (in `api/`) | exit 0 |
| Run Tests | `go test ./...` (in `api/`) | exit 0 |

## Scope

**In scope**:
- `api/internal/domain/employment.go`
- `api/internal/handler/employment_handler.go`
- `api/internal/dto/response.go`

**Out of scope**:
- Database schema / migrations (salary storage in database remains unchanged).
- Modifying write operations (managers/admins still supply salary on POST).

## Git workflow

- Branch: `advisor/002-compensation-field-masking`
- Commit style: Conventional Commits, e.g. `feat(api): mask employee compensation fields for staff and non-owner roles`

## Steps

### Step 1: Add a sanitization / masking method to `EmploymentRecord`
In `api/internal/domain/employment.go`, add a method `MaskSensitiveDetails()` on `EmploymentRecord` or pointer:
```go
func (e *EmploymentRecord) MaskSensitiveDetails() {
	e.SalaryAmount = nil
	e.SalaryCurrency = nil
	e.PayFrequency = nil
	e.HourlyRate = nil
}
```

**Verify**: Run `go build ./...` in `api/` → exit 0.

### Step 2: Implement permission check in `EmploymentHandler`
In `api/internal/handler/employment_handler.go`, create a helper `canViewCompensation(c *fiber.Ctx, targetPersonID uuid.UUID) bool`:
```go
func canViewCompensation(c *fiber.Ctx, targetPersonID uuid.UUID) bool {
	role := middleware.GetRole(c)
	if role == "admin" || role == "manager" {
		return true
	}
	callerPersonID := middleware.GetPersonID(c)
	return callerPersonID != uuid.Nil && callerPersonID == targetPersonID
}
```

### Step 3: Apply masking in `EmploymentHandler.List` and `GetCurrent`
1. In `EmploymentHandler.List`:
   ```go
   if !canViewCompensation(c, pid) {
       for i := range list {
           list[i].MaskSensitiveDetails()
       }
   }
   ```
2. In `EmploymentHandler.GetCurrent`:
   ```go
   if !canViewCompensation(c, pid) {
       cur.MaskSensitiveDetails()
   }
   ```

**Verify**:
Run `go build ./...` in `api/` → exit 0.
Run `go vet ./...` in `api/` → exit 0.

## Test plan

- Test that an authenticated caller with role `admin` receives `salary_amount` populated.
- Test that an authenticated caller with role `staff` inspecting *their own* record (`callerPersonID == targetPersonID`) receives `salary_amount`.
- Test that an authenticated caller with role `staff` or `read-only` inspecting *another* employee receives `salary_amount: null` (omitted from JSON response).

## Done criteria

- [ ] `go build ./...` in `api/` exits 0.
- [ ] `go vet ./...` in `api/` exits 0.
- [ ] No compensation fields are leaked to unauthorized staff users in API responses.
- [ ] `plans/README.md` status row updated.

## STOP conditions

- If `domain.EmploymentRecord` is used by repository scanner, ensure `MaskSensitiveDetails` modifies only the in-memory response struct, never mutating database values.

## Maintenance notes

- If additional sensitive HR attributes (e.g. performance ratings or disciplinary notes) are added to `employment_records`, include them in `MaskSensitiveDetails()`.

