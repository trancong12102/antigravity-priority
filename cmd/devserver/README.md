# Local Runtime Dev Server

`cmd/devserver` starts a local CPA emulation that never touches the real Antigravity network. It runs the production Runtime, Planner, Host Transition, State Cache and management API, and only replaces the CPA Host and the Antigravity quota endpoint with local implementations.

## Start

```bash
go run ./cmd/devserver
```

The default address is `http://localhost:8080/status`. By default it generates 10 auth files:

```text
data/devserver/auth-files/
data/devserver/quota-state.json
data/devserver/refresh-cache.json
```

Existing `.json` auth files are read and kept. The account count flag is the minimum number of accounts in the directory; manually edited files are never overwritten.

## Flags

```text
-addr            HTTP listen address, default :8080
-accounts        minimum number of generated accounts, default 10
-auth-dir        CPA auth JSON directory
-quota-state     simulated quota state file
-state-cache     Runtime state cache file
-seed            random seed; 0 uses a time-based seed
-auto-apply      turn on Runtime auto scheduling at startup
-management-key  require this Management Key on /v0/management/ routes, as CPA does; empty disables the check
```

Each flag also has a matching `ANTIGRAVITY_DEVSERVER_*` environment variable (`ANTIGRAVITY_DEVSERVER_MANAGEMENT_KEY` for `-management-key`).

## Management session

The page never asks for the Management Key: it reuses the session CPAMC stores on the same origin (`isLoggedIn` and the obfuscated `cli-proxy-auth` entry in `localStorage`). The devserver has no CPAMC, so seed that session yourself in the browser console on `http://localhost:8080/status`, then press **Retry session**:

```js
localStorage.setItem("isLoggedIn", "true");
localStorage.setItem("cli-proxy-auth", JSON.stringify({ state: { managementKey: "dev-management-key" } }));
```

A plain JSON value is accepted the same way CPAMC's legacy storage is. Without `-management-key` any key is accepted; with `-management-key dev-management-key` a wrong key brings back the session notice, which is how the 401 path can be checked locally. Clearing `isLoggedIn` from another tab locks the page, and setting it again reconnects it automatically.

The theme follows CPAMC's saved theme; to preview one, set `localStorage.setItem("cli-proxy-theme", JSON.stringify({ state: { theme: "dark" } }))` and reload, or call `setTheme("light" | "white" | "dark")` in the console.

## Simulation rules

- Auth files use the CPA Antigravity credential format; the Runtime reads and writes them through `host.auth.list/get/get_runtime/save`.
- `host.http.do` answers locally with `retrieveUserQuotaSummary`-compatible JSON and makes no real network request.
- One credential probe returns the 5h and 7d windows for both the Gemini and Claude/GPT model groups.
- The first probe returns full quota; later probes consume a random amount from both model groups.
- Once any window of any model group is exhausted, the next probe of that credential restores full quota for both groups.
- Quota state, sample history, priority write-back and scheduler diagnostics are persisted by the simulated Host and the production Runtime respectively.

## 429 circuit breaker and recovery (manual testing)

The dev server exposes helper endpoints that inject a 429 failure or a successful request, for checking the 429 circuit breaker and self-healing by hand:

```bash
# Inject a 429 error for a credential (two in a row demote it to -1 and start the 429 cooldown):
curl -X POST "http://localhost:8080/dev/simulate-429?auth_index=auth-001"

# Report a successful request for a credential (resets the consecutive 429 count, or ends the cooldown early):
curl -X POST "http://localhost:8080/dev/simulate-success?auth_index=auth-001"
```

Without `auth_index`, the first available credential is used.
