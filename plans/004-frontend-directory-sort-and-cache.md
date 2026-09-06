# Plan 004: Frontend Directory Server-Side Sort Alignment and Hover Cache Eviction

> **Executor instructions**: Follow this plan step by step. Run every verification command and confirm the expected result before moving to the next step. If anything in the "STOP conditions" section occurs, stop and report — do not improvise. When done, update the status row for this plan in `plans/README.md`.
>
> **Drift check (run first)**: `git diff --stat e2eb40e..HEAD -- web/src/components/Directory.tsx web/src/components/PersonHoverCard.tsx`
> If any in-scope file changed since this plan was written, compare the "Current state" excerpts against the live code before proceeding; on a mismatch, treat it as a STOP condition.

## Status

- **Priority**: P1
- **Effort**: S
- **Risk**: LOW
- **Depends on**: none
- **Category**: correctness / performance
- **Planned at**: commit `e2eb40e`, 2026-09-06

## Why this matters

1. In `Directory.tsx:112-120`, the frontend attempts to sort `res.data` in-memory after fetching from the API. Because the API returns a paginated slice (e.g. 20 of 1,003 items), sorting only within the client's current slice produces confusing, out-of-order pagination (e.g. page 1 sorted, page 2 sorted independently). Furthermore, the backend SQL query (`person_repo.go:170-195`) ALREADY sorts all rows by `first_name` or `city`. Removing client-side re-sorting allows the server-sorted page results to render faithfully.
2. In `PersonHoverCard.tsx:9`, `cache` is a raw unbounded JavaScript `Map`. Hovering over many employees in a long browsing session fills memory with stale entries that never expire. Capping this cache with a bounded size and TTL prevents memory creep.

## Current state

- In `web/src/components/Directory.tsx:111-122`:
  ```ts
  let list = [...res.data];
  if (sortVal === "first_name") {
    list.sort((a, b) => `${a.first_name} ${a.last_name}`.localeCompare(`${b.first_name} ${b.last_name}`));
  } else if (sortVal === "-first_name") {
    list.sort((a, b) => `${b.first_name} ${b.last_name}`.localeCompare(`${a.first_name} ${a.last_name}`));
  } else if (sortVal === "city") {
    list.sort((a, b) => (a.city || "").localeCompare(b.city || ""));
  } else if (sortVal === "-city") {
    list.sort((a, b) => (b.city || "").localeCompare(a.city || ""));
  }
  setPersons(list);
  ```
- In `web/src/components/PersonHoverCard.tsx:9`:
  ```ts
  const cache = new Map<string, Details>();
  ```
  Unbounded growth with no eviction policy.

## Commands you will need

| Purpose | Command | Expected on success |
|---------|---------|---------------------|
| Check Web Types | `npm run check` (in `web/`) | 0 errors, 0 warnings |
| Build Web | `npm run build` (in `web/`) | exit 0, dist/ generated |

## Scope

**In scope**:
- `web/src/components/Directory.tsx`
- `web/src/components/PersonHoverCard.tsx`

**Out of scope**:
- Backend Go code or API contracts.

## Git workflow

- Branch: `advisor/004-frontend-directory-sort-and-cache`
- Commit style: Conventional Commits, e.g. `fix(web): rely on server-side paginated sorting and bound hover card cache`

## Steps

### Step 1: Remove client-side sorting from `Directory.tsx`
In `web/src/components/Directory.tsx`, replace lines 111-122:
```ts
setPersons(res.data);
```
Remove the local in-memory `.sort()` block. The server's SQL query (`ORDER BY ... LIMIT ... OFFSET ...`) handles full-table sorting across all pages.

**Verify**:
Run `npm run check` in `web/` → 0 errors.

### Step 2: Add bounded cache with simple eviction in `PersonHoverCard.tsx`
In `web/src/components/PersonHoverCard.tsx`, replace the raw `Map` with a bounded cache helper:
```ts
const MAX_CACHE_SIZE = 50;
const cache = new Map<string, { details: Details; timestamp: number }>();
const CACHE_TTL_MS = 5 * 60 * 1000; // 5 minutes

function getCached(id: string): Details | null {
  const entry = cache.get(id);
  if (!entry) return null;
  if (Date.now() - entry.timestamp > CACHE_TTL_MS) {
    cache.delete(id);
    return null;
  }
  return entry.details;
}

function setCached(id: string, details: Details) {
  if (cache.size >= MAX_CACHE_SIZE) {
    const oldestKey = cache.keys().next().value;
    if (oldestKey) cache.delete(oldestKey);
  }
  cache.set(id, { details, timestamp: Date.now() });
}
```
Update `useEffect` in `PersonHoverCard.tsx` to use `getCached` and `setCached`.

**Verify**:
Run `npm run check` in `web/` → 0 errors.
Run `npm run build` in `web/` → exit 0.

## Test plan

- Test sorting by first name ascending: page 1 displays employees with first names beginning with "A", and page 2 displays the next alphabetical slice without jumping back to "A".
- Verify hovering over employees continues to display details card with cached performance.

## Done criteria

- [ ] `npm run check` in `web/` exits 0 with 0 errors.
- [ ] `npm run build` in `web/` exits 0.
- [ ] No client-side `.sort()` executes on paginated API response slices.
- [ ] `plans/README.md` status row updated.

## STOP conditions

- If `api.listPersons` fails to sort on a specific column, verify the sort key against the backend sort whitelist in `person_repo.go:170`. Do not re-add client-side array sorting.

## Maintenance notes

- If client-side search filtering is ever desired, it must query the API with parameter `q` rather than filtering the local 20-item array.

