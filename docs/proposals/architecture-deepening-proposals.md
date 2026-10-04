# Architecture Deepening Proposals

This document records the **deepening opportunities** found in the `antigravity-priority` codebase during an architecture review. Following the `/codebase-design` principles, the goals are to increase module **depth**, keep related logic together for maintainers (**locality**), give callers more **leverage**, and make sure **the interface is the test surface**.

---

## Contents

1. [Proposal 1: Absorb Probe Failure Policy into the Priority Planner [Completed]](#1-absorb-probe-failure-policy-into-the-priority-planner-completed)
2. [Proposal 2: Unify Static and Dynamic Configuration into a Deep Config Module [Completed]](#2-unify-static-and-dynamic-configuration-into-a-deep-config-module-completed)
3. [Proposal 3: Split Physical Document Patching from the Host RPC Client [Completed]](#3-split-physical-document-patching-from-the-host-rpc-client-completed)
4. [Proposal 4: Deepen the State Cache into a High-Level State Engine [Completed]](#4-deepen-the-state-cache-into-a-high-level-state-engine-completed)

---

## 1. Absorb Probe Failure Policy into the Priority Planner [Completed]

- **Status**: `Completed`
- **Landed in**: `a89510d refactor(priority): absorb probe failure policy into planner and archive deepening proposals`
- **Recommendation**: `Strong`
- **Dependency category**: `in-process`
- **Files involved**:
  - `internal/priority/planner.go`
  - `internal/runtime/production_runner.go`

### Problem
`priority.Planner` is designed as a pure function, but `PlanFreshOnly` only handled valid quota evidence with `EvidenceStatusReady`. When an upstream probe failed (`EvidenceStatusProbeFailed`), the planner did not classify that domain state itself. Instead, `withProbeFailureTemporaryDisables` in the runtime orchestration layer (`production_runner.go`) imperatively looped over `plan.Items` and `plan.Changes` and patched fields afterwards (`Disabled = true, Reason = "failedQuotaFetch"`). Domain disable policy leaked into the runtime orchestration layer, and the plan output's immutability contract was broken.

### Solution
Deepen `priority.PlanFreshOnly` so it consumes all fresh evidence, including `EvidenceStatusProbeFailed`. When an account is not disabled and its probe failed, the planner itself produces a `PlanItem` with `Disabled = true` and `Reason = "failedQuotaFetch"`, and the `changes()` mechanism emits the write change automatically. Remove `withProbeFailureTemporaryDisables` from `production_runner.go` entirely.

### Benefits
- **Locality**: every decision on account promotion, normal ordering, quota exhaustion and failure disabling lives in `internal/priority`.
- **Interface & Testability**: `priority.PlanFreshOnly` is the only test surface; unit tests cover all failure-disable scenarios without starting or faking runtime orchestration.
- **Deletion Test**: removes 50+ lines of fragile after-the-fact slice traversal and mutation.

---

## 2. Unify Static and Dynamic Configuration into a Deep Config Module [Completed]

- **Status**: `Completed`
- **Landed in**: `33cb31b refactor(config): unify dynamic and static configuration into deep config module`
- **Recommendation**: `Strong`
- **Dependency category**: `in-process`
- **Files involved**:
  - `internal/config/config.go`
  - `internal/config/schedule.go`
  - `internal/state/store.go`
  - `internal/state/schedule.go`
  - `internal/runtime/runtime.go`
  - `internal/runtime/production_runner.go`
  - `internal/management/handler.go`

### Problem
System configuration was split across three shallow data structures:
1. `config.Config`: static YAML/JSON parsing at startup.
2. `state.DynamicConfig`: dynamic hot updates from the Web UI config center.
3. `state.ScheduleConfig`: schedule intervals and quiet windows.

As a result, validation, defaulting, bounds checks (such as `quota_sample_capacity`, `urgency_tolerance`, `min_change`) and hot-reload merge logic were scattered across the `config`, `runtime` and `management` packages, and `runtime.go` contained 4 shallow conversion functions.

### Solution
Deepen `internal/config` into the single Configuration Authority:
- `config.Config` owns static loading, the dynamic hot-update structure, field bounds validation, default inheritance and serialization in both directions.
- `state.Store` only stores and reads the persisted raw config snapshot; it no longer reaches into concrete config structs or business validation.
- `Runtime` and `Management` talk only to the deepened `config` module, which removes manual field copying and duplicate validation.

### Benefits
- **Locality**: all config definitions, validation constraints and serialization rules live in `internal/config`.
- **Leverage**: callers use one interface for loading, hot updates and validation, with no duplicated hand-written mapping.
- **Deletion Test**: removes shallow glue code in `runtime.go` such as the redundant `mergePersistedDynamicConfig` and `applyDynamicConfigToRuntime`.

---

## 3. Split Physical Document Patching from the Host RPC Client [Completed]

- **Status**: `Completed`
- **Recommendation**: `Worth exploring`
- **Dependency category**: `ports & adapters`
- **Files involved**:
  - `internal/host/client.go`
  - `internal/host/patcher.go`
  - `internal/host/redact.go`
  - `internal/host/types.go`

### Problem
`host.Client` used to play two adapter roles with very different physical responsibilities:
1. **CGO host RPC client**: calls the CPA host's `ListAuthFiles`, `GetAuth`, `GetRuntime` and `HTTPDo` through the CGO bridge.
2. **Physical file AST patcher**: because CPA has no native priority-update API, `Client` patched local credential JSON files in place at the AST/byte level via `os.ReadFile`, `os.CreateTemp` and atomic rename (`PatchPriority`, `ResetPriority`, `PatchDisabled`).

So `host.Client` was both a network/RPC adapter and a disk I/O storage engine. Unit tests for RPC logic often had to depend on or create real temporary files on disk.

### Solution
Put a clear **seam** between host communication and local file patching:
- Extract a dedicated `DocumentPatcher` adapter that wraps atomic file writes, AST field location and write-back.
- `host.Client` becomes purely a CGO/RPC adapter that only talks to the CPA host context.
- Move shared redaction helpers (`RedactBytes` etc.) into a separate `redact.go` module.

### Benefits
- **Locality**: disk file operations are decoupled from network/RPC callbacks, each with a clear responsibility.
- **Testability**: host RPC tests can be 100% in-memory mocks with zero disk side effects.
- **Seam Discipline**: two concrete adapters, each serving a different external environment.

---

## 4. Deepen the State Cache into a High-Level State Engine [Completed]

- **Status**: `Completed`
- **Recommendation**: `Worth exploring`
- **Dependency category**: `in-process`
- **Files involved**:
  - `internal/state/store.go`
  - `internal/state/estimator.go`
  - `internal/state/schedule.go`
  - `internal/runtime/production_runner.go`
  - `internal/runtime/probe_evidence.go`

### Problem
`state.Store` used to be a passive thread-safe JSON dictionary exposing 15 fine-grained getters and setters.
Callers (such as `production_runner.go`) had to orchestrate complex internal state transitions explicitly:
- call `store.ClearExpiredCooldowns(now)` by hand;
- call `store.NeedsProbe(...)` by hand;
- convert `state.Entry` to `priority.ProbeEvidence` by hand;
- assemble fallback model group prediction evidence by hand.

In addition, `estimator.go` still contained the old `CalculateCycleBurnRate` estimator, which no production code referenced.

### Solution
Deepen `state.Store` into a high-leverage **Execution State Engine**:
- Provide semantic high-level interfaces `GetCachedEvidence` and `BuildGroupEvidence` that internally handle sample-queue sliding window updates, EMA smoothing, persistence and evidence generation.
- Provide `GetActiveCooldowns(now)`, which clears expired cooldowns internally.
- Delete the obsolete single-span estimator dead code in `estimator.go`.

### Benefits
- **Leverage**: callers make 1 semantic call for a state transition that used to need multi-step orchestration.
- **Locality**: quota sample history, burn rate estimation and circuit breaker state transitions are fully encapsulated in the state engine.
- **Deletion Test**: removes the unused redundant estimator and 80+ lines of assembly code in the runtime.
