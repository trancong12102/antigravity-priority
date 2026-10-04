# Antigravity Priority Release & Audit Checklist

This document defines the **standard review flow and change-check matrix** for every version iteration, code commit, and formal release of the `antigravity-priority` plugin. After finishing a feature or bug fix, and before committing code or publishing a Release, **go through this checklist item by item** so every related file and piece of metadata is updated.

---

## Contents
1. [Version Consistency Matrix](#1-version-consistency-matrix)
2. [Documentation Review Checklist](#2-documentation-review-checklist)
3. [CPA Host and Plugin Architecture Contracts](#3-cpa-host-and-plugin-architecture-contracts)
4. [Code Quality and Test Matrix](#4-code-quality-and-test-matrix)
5. [Standard Release Workflow](#5-standard-release-workflow)

---

## 1. Version Consistency Matrix

When bumping the version (such as `v1.0.1` $\to$ `v1.1.0`), **the version in these 4 places must match exactly**:

| Check | Target file | Key location / field | Rule |
| :--- | :--- | :--- | :--- |
| **1. Plugin registry** | `registry.json` | `plugins[0].version` | Plain three-part version, such as `"1.1.0"` (no `v` prefix) |
| **2. Runtime metadata** | `internal/runtime/runtime.go` | `buildMetadata().Version` | Must match `registry.json` exactly, such as `"1.1.0"` |
| **3. Embedded web dashboard** | `internal/management/feature_shell_assets.go` | `<span class="version-badge">v1.1.0</span>` | Must carry the `v` prefix, such as `v1.1.0` |
| **4. Release note** | `.github/release-notes/vX.Y.Z.md` | File name and top-level heading `# antigravity-priority vX.Y.Z` | File name must be `vX.Y.Z.md`; content is English-only |

---

## 2. Documentation Review Checklist

Public, user-facing docs should stay clear and professional and **focus on core feature value**. In this fork, the README, release notes, and UI are English-only, with no Chinese text. Follow these rules:

1. **Focus on core features**:
   - The overview lists only high-level core capabilities (such as Antigravity dual-window scheduling, Dynamic Boost, Weekly Urgency rotation, adaptive learning rate, self-healing soft depletion, the UI config center, and the dual-theme dashboard).
   - **Do not expand on technical details** (no separate bullets for development-level optimizations such as "reuses the CPA host path", "Fresh Evidence gate", or "two-phase in-memory computation").
2. **No internal codenames**:
   - **Never expose internal `(REQ-xx)` tags in `README.md` or release notes**.

| Document | Review focus | Checks |
| :--- | :--- | :--- |
| **`README.md`** | Feature overview, minimal config, management API | - Focused on core features, no pile of small implementation details;<br>- English-only, no Chinese text;<br>- **No internal `(REQ-xx)` tags**;<br>- Config docs reflect the "minimal YAML + web config center" model;<br>- API list stays in sync with route changes. |
| **`.github/release-notes/vX.Y.Z.md`** | English release note | - **Must cover every commit in `git log <last-tag>..HEAD`; never drop intermediate commits**;<br>- English-only, a single English body with no Chinese section (follow `TEMPLATE.md`);<br>- Organize core feature upgrades, performance improvements, and UX improvements by module;<br>- **No internal `(REQ-xx)` tags**. |
| **`docs/requirements/`** | Requirements and technical spec roadmap | - Update status from `Ready` to `Completed`;<br>- Covers requirement specs, domain model, and architecture design (internal docs may keep REQ numbers). |

---

## 3. CPA Host and Plugin Architecture Contracts

| Dimension | Principle and check | Pass criteria |
| :--- | :--- | :--- |
| **Minimal host config** | The CPA host `config.yaml` is only the basic load switch and carries no day-to-day business parameters. | User YAML needs only `enabled: true`; the `state_cache_path` file path defaults to `data/antigravity-priority-cache.json`. |
| **Clean metadata (`buildMetadata`)** | The official CPA management UI renders a drawer form from `Metadata.ConfigFields`, and saving that drawer writes back to `config.yaml`, which forces CPA to reload. So that the config center alone manages business parameters, **`ConfigFields` stays `nil`** (option A, minimal and non-intrusive). | `buildMetadata()` contains no `ConfigFields`, avoiding form conflicts and unexpected reloads. |
| **Layered config baseline (`config.Default()`)** | `config.Default()` is the **single source of truth** for the cold-start baseline when no cache file exists and for the UI's restore-defaults action. | The defaults returned by `Default()` (15m, gemini, 6, 1, 999, 100) must match the recommended config exactly. |
| **Dynamic config persistence and hot reload** | Business and scheduling parameters are saved to the `app_config` node of `data/antigravity-priority-cache.json`; a change immediately resets the ticker worker. | Dynamic config is persisted (`SaveAtomic`), merged automatically at startup and on Reconfigure, and takes effect instantly without a restart. |

---

## 4. Code Quality and Test Matrix

After each feature or bug fix, run lint locally first, then the unit tests. Before committing or releasing, run the same lint again and complete the full quality matrix. Every check must pass **100% with zero warnings**.

### 4.1 Local verification after finishing a feature

After writing the code, run these in this order. `go vet` does not replace `golangci-lint`:

```bash
# 1. Lint (must use golangci-lint v2.12.2, same as the release workflow)
golangci-lint run --timeout=5m

# 2. Unit tests
go test ./...
```

### 4.2 Full verification before commit and release

Before committing, creating a version tag, or publishing a Release, run lint again, then the full verification commands below:

```bash
# 1. Lint (must use golangci-lint v2.12.2, same as the release workflow)
golangci-lint run --timeout=5m

# 2. Build check (all packages and submodules free of syntax/type errors)
go build ./...

# 3. Basic static analysis
go vet ./...

# 4. Full unit tests (including documentation, registry compliance, and template compliance tests)
go test -v ./...

# 5. Race detection (zero race conditions in concurrent scheduling, the probe worker pool, and cache reads/writes)
go test -race ./...
```

### Key tests
- **`TestRegistryJSON_SchemaValidation`**: Validates that `registry.json` meets the CPA plugin store standard.
- **`TestDocumentation_EnglishReadmeCompleteness`**: Checks that `README.md` exists, is English-only, and covers the core sections and CPA paths.
- **`TestStatusHTML_NoRemoteURLs`**: Ensures `StatusHTML` has no external CDN and no remote fonts/scripts, is fully self-contained offline, and meets a strict CSP.
- **`TestStatusHTML_CSSTokenCompleteness`**: Ensures the light (white) and dark theme variable sets are complete.

---

## 5. Standard Release Workflow

```text
[Feature development / requirement iteration]
         │
         ▼
[1. Update related code and docs]
   ├─ Update business code and unit tests (TDD)
   ├─ Check version consistency (registry.json, runtime.go, feature_shell_assets.go)
   ├─ Review every commit since the last release: git log <last-tag>..HEAD --oneline
   ├─ Write .github/release-notes/vX.Y.Z.md (from the full git log range, English-only, no REQ tags)
   └─ Update README.md (concise, no pile of technical details)
         │
         ▼
[2. Local verification after finishing the feature]
   ├─ golangci-lint run --timeout=5m
   └─ go test ./...
         │
         ▼
[3. Full quality check before commit and release]
   ├─ golangci-lint run --timeout=5m again
   ├─ go build ./...
   ├─ go vet ./...
   ├─ go test -v ./...
   └─ go test -race ./...
         │
         ▼
[4. User confirmation (safety baseline)]
   └─ "Never commit without explicit user confirmation during the conversation"
         │
         ▼
[5. Git commit and tag]
   ├─ Commit messages follow Conventional Commits (such as chore: bump version to v1.1.0)
   └─ Create the Git tag: git tag vX.Y.Z
         │
         ▼
[6. Automated GitHub Release]
   └─ CI cross-compiles the dynamic library artifacts for each platform (.so / .dylib / .dll)
```
