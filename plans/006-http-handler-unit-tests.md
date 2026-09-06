# Plan 006: Standalone Unit Tests for Handlers, Auth Gates, and Middleware

> **Executor instructions**: Follow this plan step by step. Run every verification command and confirm the expected result before moving to the next step. If anything in the "STOP conditions" section occurs, stop and report — do not improvise. When done, update the status row for this plan in `plans/README.md`.
>
> **Drift check (run first)**: `git diff --stat e2eb40e..HEAD -- api/internal/auth/ api/internal/middleware/`
> If any in-scope file changed since this plan was written, compare the "Current state" excerpts against the live code before proceeding; on a mismatch, treat it as a STOP condition.

## Status

- **Priority**: P1
- **Effort**: M
- **Risk**: LOW
- **Depends on**: plans/001-tenant-isolation-idor.md
- **Category**: test-coverage
- **Planned at**: commit `e2eb40e`, 2026-09-06

## Why this matters

The Go backend currently possesses only two unit tests in `request_test.go` (testing pagination math). The entire `integration_test.go` suite requires `TEST_DATABASE_URL` and skips in standard CI when a Postgres container is unavailable. As a consequence, security-critical authentication logic (`auth.Service`), role verification (`RequireRole`), tenant retrieval (`GetOrgID`), and error handling (`ErrorHandler`) have 0% automated test coverage. This plan implements isolated unit tests that execute in <1 second without requiring a live PostgreSQL instance.

## Current state

- In `api/internal/auth/auth.go`, `IssueToken`, `ParseToken`, `HashPassword`, and `CheckPassword` have no tests.
- In `api/internal/middleware/auth.go`, `RequireAuth` and `RequireRole` have no tests.
- In `api/internal/middleware/error.go`, RFC 7807 problem formatting has no tests.
- In `.github/workflows/ci.yml:21`, `go test ./...` skips integration tests and only verifies 2 lines in `request_test.go`.

## Commands you will need

| Purpose | Command | Expected on success |
|---------|---------|---------------------|
| Run Unit Tests | `go test ./... -v` (in `api/`) | all tests pass, exit 0 |
| Check Coverage | `go test -cover ./...` (in `api/`) | coverage metrics reported |

## Scope

**In scope**:
- `api/internal/auth/auth_test.go` (new)
- `api/internal/middleware/auth_test.go` (new)
- `api/internal/middleware/error_test.go` (new)

**Out of scope**:
- Modifying database schemas or external integration test setup.

## Git workflow

- Branch: `advisor/006-http-handler-unit-tests`
- Commit style: Conventional Commits, e.g. `test(api): add unit tests for auth tokens, role middleware, and error handling`

## Steps

### Step 1: Create `api/internal/auth/auth_test.go`
Create unit tests verifying:
1. `HashPassword` generates a valid bcrypt hash and `CheckPassword` verifies correct vs incorrect password.
2. `IssueToken` generates a JWT that parses back claims (`Username`, `PersonID`, `OrgID`, `Role`).
3. `ParseToken` rejects tokens signed with a different secret or expired tokens.

```go
package auth

import (
	"testing"
	"time"

	"employee-directory-api/internal/domain"

	"github.com/google/uuid"
)

func TestHashAndCheckPassword(t *testing.T) {
	pw := "secretPassword123!"
	hash, err := HashPassword(pw)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}
	if !CheckPassword(hash, pw) {
		t.Fatalf("expected valid password check")
	}
	if CheckPassword(hash, "wrong") {
		t.Fatalf("expected invalid password check to fail")
	}
}

func TestIssueAndParseToken(t *testing.T) {
	svc, err := NewService("test-secret-key-12345", "1h")
	if err != nil {
		t.Fatalf("NewService failed: %v", err)
	}

	personID := uuid.New()
	orgID := uuid.New()
	acc := &domain.PersonAccount{
		ID:       uuid.New(),
		Username: "tester",
		PersonID: personID,
		Role:     "manager",
	}

	tok, err := svc.IssueToken(acc, orgID)
	if err != nil {
		t.Fatalf("IssueToken failed: %v", err)
	}

	claims, err := svc.ParseToken(tok)
	if err != nil {
		t.Fatalf("ParseToken failed: %v", err)
	}
	if claims.Username != "tester" || claims.PersonID != personID || claims.OrgID != orgID || claims.Role != "manager" {
		t.Fatalf("claims mismatch: %+v", claims)
	}
}
```

**Verify**: Run `go test ./internal/auth -v` in `api/` → PASS.

### Step 2: Create `api/internal/middleware/auth_test.go`
Create HTTP middleware unit tests using Fiber's in-memory `app.Test(...)`:
1. `RequireAuth` returns 401 when Authorization header is missing or lacks "Bearer".
2. `RequireAuth` returns 401 when token is invalid.
3. `RequireAuth` sets `claimsKey` in locals on valid token.
4. `RequireRole` returns 403 when role does not match allowed roles.
5. `RequireRole` returns 200 when role matches.

**Verify**: Run `go test ./internal/middleware -v` in `api/` → PASS.

### Step 3: Create `api/internal/middleware/error_test.go`
Create unit tests for `ErrorHandler`:
1. Fiber errors (e.g. `fiber.NewError(404, "not found")`) produce status 404 with RFC 7807 envelope.
2. Internal errors (status 500) do not expose sensitive internal error text.

**Verify**:
Run `go test ./...` in `api/` → all unit tests execute and PASS.

## Test plan

- Execute `go test ./... -v` without `TEST_DATABASE_URL` environment variable set.
- All unit test files must pass without network or database dependencies.

## Done criteria

- [ ] `api/internal/auth/auth_test.go` exists and passes.
- [ ] `api/internal/middleware/auth_test.go` exists and passes.
- [ ] `api/internal/middleware/error_test.go` exists and passes.
- [ ] `go test ./...` in `api/` exits 0.
- [ ] `plans/README.md` status row updated.

## STOP conditions

- If tests require database connection, stop: these must remain pure unit tests utilizing mocks or Fiber's in-memory `app.Test(req)`.

## Maintenance notes

- When new middleware (e.g. rate limiting or audit logging) is added, add corresponding unit tests in `internal/middleware/`.

