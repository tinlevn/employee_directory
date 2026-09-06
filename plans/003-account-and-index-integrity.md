# Plan 003: Account Identity Integrity and Partial Index Optimization

> **Executor instructions**: Follow this plan step by step. Run every verification command and confirm the expected result before moving to the next step. If anything in the "STOP conditions" section occurs, stop and report — do not improvise. When done, update the status row for this plan in `plans/README.md`.
>
> **Drift check (run first)**: `git diff --stat e2eb40e..HEAD -- api/internal/handler/auth_handler.go api/internal/repository/auth_repo.go api/migrations/`
> If any in-scope file changed since this plan was written, compare the "Current state" excerpts against the live code before proceeding; on a mismatch, treat it as a STOP condition.

## Status

- **Priority**: P1
- **Effort**: S
- **Risk**: LOW
- **Depends on**: none
- **Category**: bug / performance
- **Planned at**: commit `e2eb40e`, 2026-09-06

## Why this matters

1. Currently, `person_accounts` has a unique constraint on `username`, but NOT on `person_id`. An employee can have multiple disjoint accounts created, or a malicious actor can register an account linked to an employee who already possesses an account.
2. In `employment_records`, the primary directory listing query repeatedly executes `LEFT JOIN employment_records e ON e.person_id = p.id AND e.is_current = true`. While `person_id` has an index, missing a partial index on `(person_id) WHERE is_current = true` will cause query performance degradation as SCD2 version history accumulates.

This plan adds migration `000003_account_and_employment_indexes` and enforces single-account-per-person validation in `AuthHandler.Register`.

## Current state

- In `api/migrations/000001_init.up.sql:270-286`:
  ```sql
  CREATE TABLE person_accounts (
    ...
    person_id UUID NOT NULL REFERENCES persons(id) ON DELETE CASCADE,
    username VARCHAR(200) UNIQUE NOT NULL,
    ...
  );
  CREATE INDEX idx_accounts_person ON person_accounts(person_id);
  ```
  `person_id` is indexed but not unique.
- In `api/internal/handler/auth_handler.go:56-60`, `Register` checks `GetAccountByUsername`, but never checks if an account already exists for `personID`.
- In `api/migrations/000001_init.up.sql:81`, `idx_employment_person` indexes `person_id` without filtering on `is_current = true`.

## Commands you will need

| Purpose | Command | Expected on success |
|---------|---------|---------------------|
| Build API | `go build ./...` (in `api/`) | exit 0, no errors |
| Vet API | `go vet ./...` (in `api/`) | exit 0 |
| Run Tests | `go test ./...` (in `api/`) | exit 0 |

## Scope

**In scope**:
- `api/migrations/000003_account_and_employment_indexes.up.sql` (new)
- `api/migrations/000003_account_and_employment_indexes.down.sql` (new)
- `api/internal/repository/auth_repo.go`
- `api/internal/handler/auth_handler.go`

**Out of scope**:
- Changing existing `000001` or `000002` migration files (migrations must remain immutable).

## Git workflow

- Branch: `advisor/003-account-and-index-integrity`
- Commit style: Conventional Commits, e.g. `feat(db): add unique person_account constraint and partial employment index`

## Steps

### Step 1: Create Migration `000003_account_and_employment_indexes.up.sql`
Create file `api/migrations/000003_account_and_employment_indexes.up.sql`:
```sql
-- Enforce 1:1 account per person
CREATE UNIQUE INDEX idx_unique_account_per_person ON person_accounts(person_id);

-- Partial index for high-traffic current employment lookup
CREATE INDEX idx_employment_current_person ON employment_records(person_id) WHERE is_current = TRUE;
```

### Step 2: Create Migration `000003_account_and_employment_indexes.down.sql`
Create file `api/migrations/000003_account_and_employment_indexes.down.sql`:
```sql
DROP INDEX IF EXISTS idx_employment_current_person;
DROP INDEX IF EXISTS idx_unique_account_per_person;
```

### Step 3: Add `GetAccountByPersonID` to `AuthRepository`
In `api/internal/repository/auth_repo.go`, add:
```go
func (r *AuthRepository) GetAccountByPersonID(ctx context.Context, personID uuid.UUID) (*domain.PersonAccount, error) {
	q := `SELECT id, person_id, username, password_hash, role, is_active FROM person_accounts WHERE person_id = $1 LIMIT 1`
	var acc domain.PersonAccount
	err := r.pool.QueryRow(ctx, q, personID).Scan(&acc.ID, &acc.PersonID, &acc.Username, &acc.PasswordHash, &acc.Role, &acc.IsActive)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &acc, nil
}
```

### Step 4: Validate in `AuthHandler.Register`
In `api/internal/handler/auth_handler.go`, after resolving `orgID` and before creating the account:
```go
if existingAcc, err := h.repo.GetAccountByPersonID(c.Context(), personID); err != nil {
	return middleware.RepositoryError(err)
} else if existingAcc != nil {
	return fiber.NewError(fiber.StatusConflict, "an account already exists for this person")
}
```

**Verify**:
Run `go build ./...` in `api/` → exit 0.
Run `go vet ./...` in `api/` → exit 0.

## Test plan

- Test attempting to register an account with a `person_id` that already has an account returns HTTP 409 Conflict.
- Test that migration files apply cleanly with `psql` or `golang-migrate`.

## Done criteria

- [ ] `000003_account_and_employment_indexes.up.sql` and `.down.sql` exist in `api/migrations/`.
- [ ] `AuthHandler.Register` rejects duplicate accounts for the same `person_id` with 409 Conflict.
- [ ] `go build ./...` in `api/` exits 0.
- [ ] `plans/README.md` status row updated.

## STOP conditions

- If duplicate accounts already exist in a database, the unique index creation will fail until duplicates are reconciled. (In local dev or CI, fresh seed data only assigns 1 account to Ada Lovelace).

## Maintenance notes

- Any account recovery or password reset workflow should use `GetAccountByPersonID` or `GetAccountByUsername`.

