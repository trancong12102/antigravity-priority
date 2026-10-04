# AGENTS.md

> **Note**: This document is kept identical in content to [CLAUDE.md](./CLAUDE.md). Together they define the project's core design principles, strict interaction protocol, Git commit acceptance gate, and release review rules.

---

`antigravity-priority` is an efficient, single-purpose **Google Antigravity credential quota scheduler and priority booster plugin** for CLIProxyAPI (CPA).

### Core design principles
1. **Single-provider focus**: All generic heuristics and non-Antigravity logic (such as Codex / xAI) are removed. The plugin optimizes only for Antigravity's two windows (5h short window + 7d weekly window).
2. **Adaptive burn-rate learning**: The per-cycle burn rate $C_{\text{cycle}}$ is inferred from consecutive probe deltas ($\Delta R_{\text{5h}} \ge 5\%, \Delta R_{\text{7d}} > 0$) and smoothed with an EMA, so users never configure it by hand.
3. **Globally optimal scheduling**: **Weekly Urgency** and the **Dynamic Boost Horizon** remove the "large-quota account overflows in the last 24 hours" problem.
4. **Self-healing soft depletion**: Exhausting the 5h short window only soft-demotes (`priority = -1, disabled = false`), which recovers on its own after the short window refreshes. Exhausting the weekly quota hard-disables (`priority = -1, disabled = true`).
5. **Deep CPA integration**: The embedded web management page is fully self-contained (no external CDN) and follows the CPA host's dark/light themes.

---

## Strict Interaction Protocol

Unless the user explicitly invokes a skill or command such as `/<command>` (for example `/impl`), when the user asks about a feature, asks you to investigate a problem, analyze a defect, or do something specific, **never change code or run operations without confirmation**. Agents must follow this confirmation flow:

1. **State your understanding first**: Tell the user how you understood their request and its core need.
2. **Explain your approach**: Describe how you will analyze the problem, likely root-cause directions, or the implementation strategy.
3. **List a concrete action plan**: Say what you will do (which files you will change, which commands you will run, which logic you will adjust).
4. **State expected results and impact**: Describe the expected outcome, changes in system behavior, and any potential side effects.
5. **Wait for explicit confirmation**: Give the user all of the above, and **start changing or running anything only after the user explicitly replies and confirms**.

---

## Strict Git Commit Protocol

**After finishing a feature or fix, agents must never run a Git commit on their own.** Even when the full quality gate passes, follow this reporting and confirmation flow:

1. **Report the work**: Tell the user in detail what was done, which files changed, and the quality gate results.
2. **State acceptance needs**: Say whether any part of the change needs the user to verify it on a real device or target environment (such as Web UI dark theme rendering, progress bar refresh, CPA host loading, or API interaction).
3. **Give clear acceptance steps**: Provide concrete, executable verification steps and the expected result of each.
4. **Wait for user confirmation**: Create a Git commit only after the user has personally confirmed the acceptance results and given an explicit instruction (such as "ok to commit" or "commit"). **Never commit without explicit user authorization.**

---

## Release & Pre-Commit Audit

Before each version iteration, after finishing a change, and before asking the user to commit or publishing a Release, **the agent must read and strictly follow [`docs/release-checklist.md`](./docs/release-checklist.md)**.

Check the following against that document:
1. **Four-way version consistency**: The version in `registry.json`, `internal/runtime/runtime.go`, `internal/management/feature_shell_assets.go`, and `.github/release-notes/vX.Y.Z.md` must match exactly.
2. **Public docs (this fork)**: `README.md`, release notes, and the UI are English-only, with no Chinese text. Keep docs focused on core value; **never put internal `(REQ-xx)` tags in user-facing docs or release notes**. Every release note draft starts from `git log <last-tag>..HEAD --oneline`, covering every commit since the previous release.
3. **Git history tracing (full coverage)**: never draft release notes from memory or from a single task. First run
   ```bash
   git log <last-release-tag>..HEAD --oneline
   ```
   and check every commit since the previous tag, so all features and fixes in that range (theme work, scheduling changes, style fixes and so on) make it into the notes.
4. **Architecture and metadata contracts**: Keep the CPA host `config.yaml` minimal, make sure `buildMetadata()` exposes no extra drawer fields, and keep `config.Default()` as the baseline fallback, consistent with the UI config center.
5. **Full quality verification**: After each feature or bug fix, run `golangci-lint v2.12.2` (same version as the release workflow) first, then `go test ./...`. Before committing or releasing, run `golangci-lint run --timeout=5m` again, plus `go build ./...`, `go vet ./...`, `go test -v ./...`, and `go test -race ./...`, with zero warnings and zero data races.

---

## Architecture & Ownership

- **Domain Core (`internal/core`, `internal/priority`, `internal/state`, `internal/provider/antigravity`)**:
  - `priority.Planner` must be a **side-effect-free pure function** that only outputs an immutable `Plan`.
  - The domain core must never write back to the CPA host or make external network requests.
- **Apply Layer (`internal/apply`, `internal/host`)**:
  - Only executes host `PatchPriority` / `PatchDisabled` according to the `Plan` gates (Fresh Evidence, `ForceWrite`, and `min_change`) and persists fully redacted audits.
- **Application Layer (`internal/runtime`)**:
  - Sole owner of scheduling use-case execution (`ManualApply` / `AutoApply` / `Probe` / `SyncHost` / `ResetAllPriorities`).
  - Coordinates the ticker and the single-flight mutex (`runMu`).
- **Adapters (`main.go`, `internal/management`)**:
  - The CGO ABI exports and the HTTP Management API are thin adapters that pass calls through to Runtime use cases; **they must never compute or rewrite priority logic themselves**.
  - The Web UI (`feature_shell_assets.go`) is plain native frontend code with no external CDN, styled entirely after the CPA theme.

---

## Domain Docs & Configuration

## Agent skills

### Issue tracker

Issues and specs are tracked as local markdown files under `.scratch/<feature>/`. See `docs/agents/issue-tracker.md`.

### Triage labels

Canonical five-role triage labels (`needs-triage`, `needs-info`, `ready-for-agent`, `ready-for-human`, `wontfix`). See `docs/agents/triage-labels.md`.

### Domain docs

Single-context repository. See `docs/agents/domain.md`, `CONTEXT.md`, and `docs/adr/`.
