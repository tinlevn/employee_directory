# Plan 005: Frontend Auth Navigation, Session Indicator, and Dynamic Config

> **Executor instructions**: Follow this plan step by step. Run every verification command and confirm the expected result before moving to the next step. If anything in the "STOP conditions" section occurs, stop and report — do not improvise. When done, update the status row for this plan in `plans/README.md`.
>
> **Drift check (run first)**: `git diff --stat e2eb40e..HEAD -- web/src/layouts/Layout.astro web/src/lib/api.ts`
> If any in-scope file changed since this plan was written, compare the "Current state" excerpts against the live code before proceeding; on a mismatch, treat it as a STOP condition.

## Status

- **Priority**: P1
- **Effort**: S
- **Risk**: LOW
- **Depends on**: none
- **Category**: dx / ux
- **Planned at**: commit `e2eb40e`, 2026-09-06

## Why this matters

`LoginForm.tsx` and `RegisterForm.tsx` exist under `/login` and `/register`, but `Layout.astro` contains zero links to either page, no sign-out capability, and no indication of who is currently authenticated. If an employee logs in, they cannot sign out without opening browser devtools and manually deleting `localStorage.getItem("token")`. Additionally, `Layout.astro:32` hardcodes `http://localhost:8080/health`, which breaks whenever running on a remote server or Docker port mapping. This plan adds an interactive session navigation widget and dynamic environment-based URLs.

## Current state

- In `web/src/layouts/Layout.astro:29-35`:
  ```astro
  <div class="flex items-center gap-5 text-sm font-medium">
    <a href="/" ...>Directory</a>
    <a href="/analytics" ...>Analytics</a>
    <a href="http://localhost:8080/health" target="_blank" ...>API Health</a>
    <button id="theme-toggle" ...></button>
  </div>
  ```
  No link to `/login` or `/register`, no "Logout" button, and health URL is hardcoded.

## Commands you will need

| Purpose | Command | Expected on success |
|---------|---------|---------------------|
| Check Web Types | `npm run check` (in `web/`) | 0 errors, 0 warnings |
| Build Web | `npm run build` (in `web/`) | exit 0 |

## Scope

**In scope**:
- `web/src/components/AuthNav.tsx` (new React island)
- `web/src/layouts/Layout.astro`
- `web/src/lib/api.ts`

**Out of scope**:
- Backend API handlers.

## Git workflow

- Branch: `advisor/005-frontend-auth-layout-and-nav`
- Commit style: Conventional Commits, e.g. `feat(web): add AuthNav session widget and dynamic health link`

## Steps

### Step 1: Create `AuthNav.tsx` Component
Create `web/src/components/AuthNav.tsx`:
```tsx
import { useEffect, useState } from "react";

export default function AuthNav() {
  const [hasToken, setHasToken] = useState(false);

  useEffect(() => {
    setHasToken(!!localStorage.getItem("token"));
  }, []);

  function handleLogout() {
    localStorage.removeItem("token");
    setHasToken(false);
    window.location.href = "/login";
  }

  if (hasToken) {
    return (
      <button
        onClick={handleLogout}
        className="text-xs font-medium px-2.5 py-1.5 rounded-md border border-[#E6DBC5] dark:border-[#2b303c] hover:bg-red-500/10 hover:border-red-500/30 text-red-600 dark:text-red-400 transition-colors"
      >
        Sign out
      </button>
    );
  }

  return (
    <div className="flex items-center gap-2 text-xs">
      <a
        href="/login"
        className="px-2.5 py-1.5 rounded-md text-[#141E46] dark:text-slate-300 hover:text-[#41B06E] dark:hover:text-[#1DCD9F] transition-colors"
      >
        Sign in
      </a>
      <a
        href="/register"
        className="px-2.5 py-1.5 rounded-md bg-[#41B06E] dark:bg-[#1DCD9F] text-white dark:text-[#131519] font-semibold hover:opacity-90 transition-opacity"
      >
        Register
      </a>
    </div>
  );
}
```

### Step 2: Update `Layout.astro`
In `web/src/layouts/Layout.astro`:
1. Import `AuthNav`:
   ```astro
   import AuthNav from "../components/AuthNav";
   const apiBase = import.meta.env.PUBLIC_API_URL || "http://localhost:8080";
   ```
2. Replace hardcoded `http://localhost:8080/health` with `${apiBase}/health`.
3. Render `<AuthNav client:load />` in the top right navigation bar right before the theme toggle.

**Verify**:
Run `npm run check` in `web/` → 0 errors.
Run `npm run build` in `web/` → exit 0.

## Test plan

- Visit `/` with no token in localStorage: "Sign in" and "Register" buttons appear in the top navbar.
- Log in through `/login`: Navbar renders "Sign out".
- Click "Sign out": token is removed from localStorage and user is redirected to `/login`.

## Done criteria

- [ ] `npm run check` exits 0.
- [ ] `npm run build` exits 0.
- [ ] User can sign in, register, and sign out directly from the top navigation.
- [ ] `plans/README.md` status row updated.

## STOP conditions

- Astro renders static HTML; `AuthNav` MUST use `client:load` (or `client:idle`) so it can read `localStorage` in the browser client.

## Maintenance notes

- If user profile display (e.g. username or role badge) is needed, `AuthNav` can decode the JWT payload to display the authenticated username.

