# Plan 007: OpenAPI 3.0 Contract Specification Update for Auth and Security

> **Executor instructions**: Follow this plan step by step. Run every verification command and confirm the expected result before moving to the next step. If anything in the "STOP conditions" section occurs, stop and report — do not improvise. When done, update the status row for this plan in `plans/README.md`.
>
> **Drift check (run first)**: `git diff --stat e2eb40e..HEAD -- api/openapi.yaml`
> If any in-scope file changed since this plan was written, compare the "Current state" excerpts against the live code before proceeding; on a mismatch, treat it as a STOP condition.

## Status

- **Priority**: P2
- **Effort**: S
- **Risk**: LOW
- **Depends on**: none
- **Category**: docs / api-contract
- **Planned at**: commit `e2eb40e`, 2026-09-06

## Why this matters

The repository maintains a checked-in OpenAPI 3.0 specification at `api/openapi.yaml`, which is served by the API at `/openapi.yaml` and `/swagger`. However, this contract was written before authentication was implemented: it does not document `/api/v1/auth/register` or `/api/v1/auth/login`, and lacks the `securitySchemes` definition for `BearerAuth`. Consequently, API consumers and contract linters see all endpoints as anonymous public endpoints. This plan synchronizes `openapi.yaml` with the live API implementation.

## Current state

- In `api/openapi.yaml:8-58`, paths start directly with `/persons` and `/analytics/headcount`.
- No `/auth/register` or `/auth/login` paths are documented.
- No `components.securitySchemes` block exists.

## Commands you will need

| Purpose | Command | Expected on success |
|---------|---------|---------------------|
| Validate YAML syntax | `npx -y yaml-lint api/openapi.yaml` or python yaml check | valid YAML |
| Build API | `go build ./...` (in `api/`) | exit 0 |

## Scope

**In scope**:
- `api/openapi.yaml`

**Out of scope**:
- Modifying Go routes or business logic.

## Git workflow

- Branch: `advisor/007-openapi-contract-update`
- Commit style: Conventional Commits, e.g. `docs(api): document auth endpoints and bearer security scheme in openapi.yaml`

## Steps

### Step 1: Add `/auth/register` and `/auth/login` paths to `api/openapi.yaml`
Under `paths:`, add:
```yaml
  /auth/register:
    post:
      summary: Register employee account
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              required: [username, password, person_id]
              properties:
                username: {type: string}
                password: {type: string}
                person_id: {type: string, format: uuid}
                role: {type: string, enum: [staff, manager, admin, read-only], default: staff}
      responses:
        '201': {description: Account created, token returned}
        '400': {description: Invalid input}
        '409': {description: Username or Person ID conflict}
  /auth/login:
    post:
      summary: Authenticate employee
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              required: [username, password]
              properties:
                username: {type: string}
                password: {type: string}
      responses:
        '200': {description: Authentication successful}
        '401': {description: Invalid credentials}
```

### Step 2: Add `BearerAuth` to `components.securitySchemes`
In `api/openapi.yaml`, under `components:`, add:
```yaml
  securitySchemes:
    BearerAuth:
      type: http
      scheme: bearer
      bearerFormat: JWT
```
And add global security requirement at root level:
```yaml
security:
  - BearerAuth: []
```
(Exempt `/auth/login` and `/auth/register` with `security: []`).

**Verify**:
Run `go build ./...` in `api/` → exit 0.
Inspect `/openapi.yaml` syntax.

## Test plan

- Validate that `api/openapi.yaml` parses without syntax errors.
- Verify GET `/openapi.yaml` serves the updated YAML specification.

## Done criteria

- [ ] `api/openapi.yaml` documents `/auth/register` and `/auth/login`.
- [ ] `api/openapi.yaml` declares `BearerAuth` security scheme.
- [ ] `plans/README.md` status row updated.

## STOP conditions

- Maintain OpenAPI version `3.0.3` for compatibility with Swagger UI. Do not upgrade to 3.1.0 unless required.

## Maintenance notes

- Any future endpoints must be documented in `openapi.yaml` in lockstep with handler creation.

