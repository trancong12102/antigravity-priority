# Architecture Deepening Proposals: Round 2

This document records the **module deepening opportunities** found in the second architecture review of `antigravity-priority`.

This round focuses on areas with frequent recent changes, mainly:

- the scheduling and execution chain across `internal/runtime`, `internal/state`, `internal/priority` and `internal/apply`;
- the recent asset split of the embedded management page in `internal/management`;
- cross-module semantics such as Fresh Evidence, the dual Model Group projection, Host transitions and execution records;
- whether the current interfaces form a stable test surface, and whether complexity leaks across seams.

The review uses `module`, `interface`, `implementation`, `depth`, `deep`, `shallow`, `seam`, `adapter`, `leverage` and `locality` as its shared architecture vocabulary, and applies a **deletion test** to each candidate.

> **Status note**: the status in each section is authoritative. Proposals that have finished design interrogation record their confirmed scope, and implementation details are carried by separate specs.

---

## Relationship to Round 1

Round 1 is recorded in [`architecture-deepening-proposals.md`](./architecture-deepening-proposals.md). These four items are complete:

1. Absorb probe failure policy into the Priority Planner.
2. Unify static and dynamic configuration into a deep config module.
3. Split physical document patching from the Host RPC adapter.
4. Deepen the state cache into a high-level state engine.

This round does not re-propose that work. Some issues sit next to Round 1 results, but they concern interface complexity, semantic leaks and caller orchestration knowledge that remain after those changes landed.

---

## Contents

1. [Unify Host transitions and execution records](#1-unify-host-transitions-and-execution-records)
2. [Consolidate ownership of Fresh Evidence](#2-consolidate-ownership-of-fresh-evidence)
3. [Deepen the dual Model Group planning projection](#3-deepen-the-dual-model-group-planning-projection)
4. [Organize management page assets by user feature](#4-organize-management-page-assets-by-user-feature)
5. [Recommended order and how the proposals combine](#5-recommended-order-and-how-the-proposals-combine)

---

## 1. Unify Host transitions and execution records

- **Status**: `Completed`
- **Landed in**: `1fe981c feat(apply): unify atomic host transition execution`
- **Recommendation**: `Strong`
- **Dependency category**: `local-substitutable`
- **Suggested priority**: `P1 / Top recommendation`
- **Files involved**:
  - `internal/apply/apply.go`
  - `internal/apply/audit.go`
  - `internal/runtime/production_runner.go`
  - `internal/runtime/runtime.go`
  - `internal/state/store.go`

### Current pain points

Normal Apply, 429 Reactive Cooldown and ResetAllPriorities are three kinds of Host state change, but they use three execution paths that are not fully consistent.

The main call chain for normal Apply is:

```text
priority.Plan
  → apply.Apply
  → Host PatchPriority / PatchDisabled
  → Runtime.snapshotRunEntry
  → json.Marshal
  → Store.SetRuntimeSnapshot
  → Store.SaveAtomic
```

429 Reactive Cooldown has `Runtime.triggerCooldown` build a `priority.Plan` and `priority.Change` by hand, call `apply.Apply`, and then generate and persist its own history record.

ResetAllPriorities bypasses `apply.Apply` entirely: it loops over `client.ResetPriority` directly, then builds its own `apply.Result`, snapshot, history record and persisted data.

Because of these forks, no single authority defines ordering, errors, results, redaction or persistence semantics for a Host transition.

### Concrete evidence

#### 1.1 The production adapter for the `Auditor` seam is an empty implementation

`apply.Request` requires an `Auditor`:

```go
type Auditor interface {
    SaveSnapshot(ctx context.Context, snapshot PlanSnapshot) error
    RecordEvent(ctx context.Context, event AuditEvent) error
}
```

`apply.Apply` calls `SaveSnapshot` and `RecordEvent` before writing to the Host. But the Runtime's production implementation only returns `ctx.Err()` and neither uses nor saves its arguments:

```go
func (r *Runtime) SaveSnapshot(ctx context.Context, snapshot apply.PlanSnapshot) error {
    return ctx.Err()
}

func (r *Runtime) RecordEvent(ctx context.Context, event apply.AuditEvent) error {
    return ctx.Err()
}
```

Real recording happens after the Host write, scattered across the Runtime's execution branches. So the write-ahead guarantee implied by the `Auditor` interface does not match production behavior.

There is only one production adapter, and it has no real implementation. Under the principle "one adapter means a hypothetical seam; two adapters prove real variation", this seam has not proven its value.

#### 1.2 Execution record persistence is duplicated across branches

The following steps are repeated in different forms across Probe, zero-change Apply, normal Apply, Reset and 429 Cooldown:

```text
snapshotRunEntry
→ json.Marshal(result)
→ json.Marshal(history)
→ store.SetRuntimeSnapshot
→ store.SaveAtomic
```

Some paths check the return value of `SaveAtomic`; others ignore the error with `_ = store.SaveAtomic(ctx)`. Whether an execution record reliably reaches disk depends on which path triggered it.

#### 1.3 Some Host transitions lack clear result semantics

`applyChange` can execute in this order:

```text
PatchPriority succeeds
→ PatchDisabled fails
```

The Host has now been partially changed, but the final result is mostly marked `failed`. The current result does not state which fields were written successfully or what partial state the Host may be in.

#### 1.4 Reset collapses attempts and failures

`ResetAllPriorities` only increments `resetCount` when `ResetPriority` succeeds, yet uses that value for both:

```go
Attempted: resetCount,
Succeeded: resetCount,
```

If 10 credentials are attempted and 7 succeed, the record may show 7 attempted and 7 succeeded; the 3 failed transitions never enter a trustworthy result model.

### Why deepening is needed

Host state changes are the core ownership of the Apply Layer. Callers and tests need one stable interface to answer:

- what the plan intends to change;
- what was actually attempted;
- which fields were written successfully;
- which fields failed;
- whether a partial success occurred;
- whether the redacted execution record was persisted;
- what state the Host may currently be in.

Today this knowledge is spread across `apply`, `runtime` and `state`, so the full behavior of one transition cannot be verified through a single seam.

### Deletion test

Deleting the current `Auditor` interface and the Runtime's empty implementation would barely change production behavior, which shows the module is shallow.

Conversely, if a deep module truly owned the Host transition and execution record lifecycle, deleting it would scatter gating, write ordering, partial success, result counting, redaction and persistence complexity back into three callers: normal Apply, Reset and 429 Cooldown. That shows this direction yields real depth and leverage.

### Deepening direction

Deepen the existing Apply module so that the following behavior of one Host transition has locality:

```text
Validate write gate
→ Capture pre-write state
→ Execute Host transition
→ Express full success, partial success, failure or skip
→ Generate redacted execution record
→ Persist execution result
```

Normal Apply, Reset and 429 Reactive Cooldown should reuse the same transition semantics. The Runtime still owns execution of the ManualApply, AutoApply, Reset and Cooldown use cases; this proposal does not move use-case ownership into the Apply module.

The concrete interface must be decided in a later design phase; no new port or adapter is assumed now.

### Goals

- Make the factual outcome of a Host transition the single authority.
- Remove the contract gap between nominal auditing and actual persistence.
- Unify gating, error and counting semantics across the three write paths.
- Keep execution records consistent with actual Host changes.
- Establish a stable test surface that covers partial success.

### Expected results

- `Attempted`, `Succeeded`, `Failed`, `Skipped` and partial-success counts are trustworthy.
- Reset, 429 Cooldown and normal Apply share one result model.
- Execution record persistence errors are no longer silently ignored.
- Tests can directly cover transitions such as "priority succeeded, disabled failed".
- Duplicated history, serialization and persistence orchestration in the Runtime shrinks noticeably.
- The recent write health shown on the diagnostics page corresponds to real Host transitions.
- The Apply Layer better matches the AGENTS.md ownership rule of "write-back gating and fully redacted audit persistence".

### Risks and constraints

- Priority computation must not move into the Apply Layer; the Planner must remain a side-effect-free pure function.
- Retry behavior after partial success must be defined explicitly, so repeated writes do not cause new state drift.
- Internal seams should not be exposed just for testing.
- If there is still only one production adapter in the end, abstractions with no behavioral difference should not be kept.
- All audit information must remain fully redacted.

### ADR relationship

This proposal does not conflict with existing ADRs. The decisions from design interrogation are recorded as [`ADR-0005`](../adr/0005-atomic-host-transition-and-truthful-outcomes.md), and they reinforce:

- the actual write semantics of Soft Depletion and Hard Depletion in ADR-0003;
- the non-destructive Host transition for 429 Reactive Cooldown in ADR-0004;
- the Apply Layer's ownership of gating, Host Patch and redacted auditing in AGENTS.md.

### Decisions confirmed in design interrogation (2026-08-22)

The following conclusions were confirmed with the user via `/grill-with-docs`; later design and implementation must not change them silently:

1. **A credential's target fields form one Host Transition**: when priority and disabled live in the same credential file, they are committed together through one file replacement, not kept as two separate business results.
2. **The result is defined by the Host's final state**: whether a request returned an error does not by itself decide whether the transition succeeded.
3. **The three operations share transition semantics**: normal Apply, 429 Reactive Cooldown and Reset Priority use consistent success, failure, verification and execution record semantics.
4. **Credentials execute independently**: one credential's failure does not block the others.
5. **Current Host state is the concurrency authority**: if an external decision-state change after the read is detected, the old Plan is abandoned rather than overwritten, and the next round of Fresh Evidence and the Planner decide again.
6. **Successful file replacement is the commit point**: cancellation before the commit point can safely return failure; context cancellation after the commit point cannot negate a Host Transition that already happened.
7. **Re-read verification after commit**: verify the target fields and JSON validity, without requiring the whole file to be byte-identical; report `uncertain` when the final state cannot be confirmed.
8. **Host outcome and Record outcome are separate**: Host state facts and execution record reliability are two orthogonal results.
9. **Concurrency comparison covers only decision-relevant state**: compare credential identity, priority and disabled; unrelated fields take the latest Host document and are left unchanged.
10. **Field targets use explicit operation semantics**: each target field expresses `unchanged`, `set` or `unset`, so all three use cases can share one transition implementation.
11. **No same-round automatic retry**: failures are recorded directly, and later scheduled rounds compensate based on new Fresh Evidence.
12. **429 Cooldown facts are persisted independently**: a failed Host demotion does not undo a cooldown that already happened.
13. **Host outcome has five terminal states**: `no_change`, `committed`, `failed`, `conflict` and `uncertain`, with the specific cause expressed through a reason.
14. **Round totals are derived from transition details**: totals must not be maintained by hand in each execution path.
15. **All three use cases produce a transition intent**: normal Apply converts from Planner changes; 429 uses `set priority=-1 / unchanged disabled`; Reset uses `unset priority / unchanged disabled`. The Planner keeps sole ownership of scheduling computation.

#### Speculative complexity explicitly rejected

The following designs were explicitly ruled out this round, unless a reproducible failure or a new real use case appears:

- a write-ahead sidecar journal;
- precise process crash recovery and a pending reconciliation protocol;
- journal schema, fingerprint identity, idempotent cleanup and a separate persistence path;
- same-round automatic retry, backoff or replay;
- extra ports or adapters for a multi-writer scenario that does not yet exist.

---

## 2. Consolidate ownership of Fresh Evidence

- **Status**: `Accepted`
- **Recommendation**: `Strong`
- **Dependency category**: `in-process`
- **Suggested priority**: `P1`
- **Files involved**:
  - `internal/priority/planner.go`
  - `internal/state/store.go`
  - `internal/runtime/production_runner.go`
  - `internal/runtime/runtime.go`

### Current pain points

`CONTEXT.md` defines Fresh Evidence as validated quota observations obtained by a successful probe in the current scheduling round.

In the code, however, whether evidence counts as Fresh Evidence has to be inferred jointly by the Store, Runtime and Planner. `priority.ProbeEvidence` exposes all of:

```text
ObservedAt
Freshness
ProbeStatus
Status
EvidenceFresh
```

Successful, failed, cached, current-round and historical prediction evidence are not expressed through one interface; they rely on combinations of several fields and on call order.

### Concrete evidence

#### 2.1 The Planner uses different freshness conditions for success and failure evidence

For success evidence the Planner checks:

```text
EvidenceFresh == true
Freshness == Fresh
ProbeStatus == Ready
Status == Ready
```

But `freshEvidenceByAuthIndex` also accepts evidence with `Status == ProbeFailed` without applying the same full current-round check to failure evidence. `initialItems` then marks that failure evidence `EvidenceFresh = true` and produces a disable decision.

#### 2.2 The Store turns cached state into evidence that looks fresh

When `state.GetCachedEvidence` builds a `priority.ProbeEvidence` from a successful cache entry, it sets `FreshnessFresh`, `ProbeStatusReady` and `EvidenceFresh = true`.

But a cache entry existing does not automatically mean "Fresh Evidence for the current scheduling round". Whether the cache may be used for read-only display, prediction or Host write-back is left for callers to decide.

#### 2.3 A private Runtime filter acts as a hidden gate

Production scheduling uses `currentRoundEvidence` to check:

```go
item.ObservedAt.Equal(observedAt)
```

to select current-round evidence. This filter protects production Apply, but it is private Runtime knowledge that is not reflected in the Planner or Store interface.

#### 2.4 Planner callers do not apply the filter consistently

`runProductionTask` calls `currentRoundEvidence` before planning.

`SyncHost` and `ResetAllPriorities`, when generating snapshots, go straight through:

```text
store.BuildGroupEvidence
→ priority.PlanFreshOnly
```

These paths currently produce read-only snapshots and do not apply the Plan, but an old failure cache entry may be shown as a current failure or a target disable. More importantly, any future caller must know about and copy this hidden gate, or it may promote historical cache into writable Fresh Evidence.

### Why deepening is needed

Fresh Evidence decides whether a Plan may produce a Host write-back; it is the scheduling system's safety gate. Callers should not have to reconstruct this domain concept by combining five status fields and comparing timestamps.

If every caller must know the cache source, observation time, probe result and intended use, the current evidence interface is too complex and easy to misuse.

### Deletion test

Deleting `currentRoundEvidence` would not remove complexity; it would immediately scatter back into every Planner caller, each of which would have to re-implement filtering on observation time, success status and failure status.

This shows that choosing between "current-round writable evidence" and "historical read-only evidence" is worth centralizing in the existing evidence module. This direction needs no new adapter; it gains depth by deepening the in-process implementation.

### Deepening direction

Make the evidence generation and selection module the single authority for Fresh Evidence, absorbing internally:

- how `ObservedAt` relates to the current scheduling round;
- the provenance of successful and failed probes;
- the difference between cached observations and Fresh Evidence;
- evidence that may be used for Host write-back;
- evidence used only for historical prediction or diagnostic display.

Callers only state their intended use and no longer interpret multiple freshness fields themselves. The concrete interface will be decided in a later design phase.

### Goals

- Give the Fresh Evidence definition exactly one implementation.
- Prevent cached evidence from being accidentally promoted to writable evidence.
- Remove repeated timestamp and status-combination checks from Runtime callers.
- Give the Planner semantically clear input with explicit provenance.
- Make the evidence rules fully testable through one stable interface.

### Expected results

- Apply, Sync, Reset and dual Model Group prediction share consistent evidence semantics.
- A probe failure in the current round only produces a diagnostic record and does not change that credential's Host state.
- An old `ProbeFailed` is not mistaken for a current-round failure.
- New Planner callers do not need to copy `currentRoundEvidence`.
- The Planner no longer depends on callers assembling five freshness fields correctly.
- Tests can directly cover current-round success, current-round failure, stale success, stale failure and historical prediction.
- The Fresh Evidence domain definition in `CONTEXT.md` has a single authority in code.

### Risks and constraints

- Cached predictions and Fresh Evidence must not be given the same write eligibility.
- The read-only prediction capability for the alternate Model Group must be kept.
- Neither current nor historical probe failures may trigger quota-driven Host write-back.
- No port with only one adapter should be added.
- The Planner's pure-function property must not be broken.

### ADR relationship

This proposal does not conflict with existing ADRs. The failure semantics from design interrogation are recorded as [`ADR-0006`](../adr/0006-failed-quota-probes-preserve-host-state.md), and they reinforce:

- the definition of Fresh Evidence in `CONTEXT.md`;
- ADR-0003's handling of real Soft Depletion and Hard Depletion, without misclassifying probe failure as depletion;
- the Round 1 state engine deepening proposal, without duplicating its caching and sampling responsibilities.

### Decisions confirmed in design interrogation (2026-08-22)

The following conclusions were confirmed with the user via `/grill-with-docs`; later design and implementation must not change them silently:

1. **Fresh Evidence comes only from a successful probe in the current round**: the data must belong to the current scheduling round, come from a successful probe, and be validated.
2. **A probe failure means unknown**: when a Manual Apply probe fails, or an Auto Apply probe still fails after retry, that credential keeps its Host state; the failure is only recorded and displayed.
3. **Credentials are planned independently**: one credential's probe failure does not stop other credentials that have Fresh Evidence from being scheduled.
4. **Historical cache is read-only**: the last successful data may be shown or used for prediction when clearly marked as historical, but it cannot authorize Host write-back; if there was never a success, show unknown.
5. **Runtime events and manual commands are authorized separately**: 429 Reactive Cooldown and Reset Priority are authorized by a real 429 event and an explicit manual action respectively, not disguised as Fresh Evidence.
6. **Later rounds compensate naturally**: this round does not modify the Host because of a probe failure and adds no new recovery flow; later scheduled rounds re-probe and decide again.

#### Speculative complexity explicitly rejected

The following are not core requirements this round, unless a verifiable new business scenario appears:

- extending probe failures into a permanent, temporary or configurable error taxonomy;
- adding Fresh Evidence TTL, failure policy or retry policy configuration;
- building a generic evidence port, adapter or cross-provider evidence platform;
- introducing a separate recovery queue, event sourcing or replay mechanism for probe failures.

---

## 3. Deepen the dual Model Group planning projection

- **Status**: `Accepted`
- **Recommendation**: `Strong`
- **Dependency category**: `in-process`
- **Suggested priority**: `P2`
- **Files involved**:
  - `internal/runtime/production_runner.go`
  - `internal/runtime/runtime.go`
  - `internal/apply/audit.go`
  - `internal/priority/planner.go`

### Current pain points

Antigravity's `gemini` and `claude_gpt` are two independent Model Groups. The system uses the configured Control Model Group to produce actual write-back, and generates a read-only prediction for the alternate Model Group.

This full projection flow is currently repeated across several Runtime use cases:

```text
Read control group evidence
→ PlanFreshOnly
→ Snapshot

Compute alternate Model Group
→ Read alternate evidence
→ PlanFreshOnly
→ SnapshotPredicted

→ NewDualGroupSnapshot
```

### Concrete evidence

#### 3.1 Production scheduling assembles the dual-group projection itself

`runProductionTask` handles evidence, Plan and Snapshot for the control group and the alternate group separately, then calls `NewDualGroupSnapshot`.

#### 3.2 SyncHost repeats the same knowledge

`SyncHost` again runs, for both the control and alternate groups:

```text
BuildGroupEvidence
→ PlanFreshOnly
→ Snapshot / SnapshotPredicted
→ NewDualGroupSnapshot
```

#### 3.3 ResetAllPriorities repeats it a third time

After Reset finishes its Host state changes, it runs the same dual-group projection again to update the management page snapshot.

#### 3.4 The LatestSnapshot fallback maintains another set of result-building rules

When there is no latest dual-group snapshot, `LatestSnapshot` builds its own fallback containing the legacy single-group result and an empty predicted Snapshot.

So every caller must know:

- where the Control Model Group comes from;
- how the alternate Model Group is computed;
- which Plan is writable;
- which Snapshot must be marked predicted;
- how cooldown options and observation time are passed;
- the shape of the fallback result.

### Why deepening is needed

This is not just duplicated code; it is a set of domain rules that must stay consistent across all callers.

Any change to prediction semantics, cooldown propagation or Snapshot structure has to be made in Production, Sync, Reset and the fallback at once. Missing any path can make the management page show different semantics after different operations.

Today the projection can only be tested through large Runtime tests that set up Host inventory, Store, Clock and two groups of evidence; the pure projection logic has no test surface of its own.

### Deletion test

Deleting any one duplicated projection block does not remove complexity; it just forces that caller to re-implement the control group, alternate group, Planner and Snapshot rules.

If the dual Model Group projection were deepened into a Runtime-internal module and then deleted, all four call paths would have to rebuild the full algorithm. This proves the module provides real leverage.

### Deepening direction

Centralize the following knowledge in a deep module inside the Runtime:

- selection of the control and alternate Model Groups;
- mapping two sets of evidence to two Plans;
- the roles of the control result and the predicted result;
- cooldown options and observation time;
- uniform construction of `DualGroupSnapshot`, including the empty state.

The module only does in-process projection; it performs no Host write-back or external network requests. The Planner stays a side-effect-free pure function, and the Runtime keeps ownership of scheduling use cases.

The concrete interface will be decided in a later design phase.

### Goals

- Give dual Model Group projection semantics exactly one implementation.
- Prevent behavioral drift between Production, Sync, Reset and the fallback.
- Extract pure computation from large Runtime use cases into a stable test surface.
- Reduce the blast radius of future changes to predicted rules.
- Keep the explicit contract that the control group is writable and the alternate group is read-only.

### Expected results

- All Runtime use cases produce a consistent `DualGroupSnapshot`.
- Callers no longer compute the alternate Model Group themselves.
- No path can forget the `SnapshotPredicted` marker.
- The control and predicted groups can be tested with pure inputs, without a full Host and disk setup.
- Runtime use cases focus on execution order rather than repeated data projection.
- If new read-only display rules are added later, the change stays in one place.

### Risks and constraints

- The predicted Plan must never reach Host write-back.
- Network probing and Store persistence must not be absorbed into the pure projection implementation.
- The Runtime's use-case ownership of ManualApply, AutoApply, Probe, SyncHost and Reset must not move.
- The Planner must not gain side effects.
- This must be coordinated with the Fresh Evidence deepening so the new module does not reinterpret freshness fields.

### ADR relationship

This proposal does not conflict with existing ADRs and is consistent with:

- ADR-0001's specialization on Antigravity as the single provider;
- the domain definitions of the `gemini` and `claude_gpt` Model Groups in `CONTEXT.md`;
- the Runtime's sole ownership of scheduling use cases in AGENTS.md.

This round produced no new, hard-to-reverse business trade-off, so no new ADR is added. The implementation spec only consolidates the dual-group contract already in effect.

### Decisions confirmed in design interrogation (2026-08-22)

The following rules are already backed by requirements, tests or earlier decisions, and were confirmed unchanged via `/grill-with-docs` this round:

1. **Config is the only control authority**: `antigravity_model_group` in Dynamic Config decides the Control Model Group; API parameters and the page selector cannot change write-back authority.
2. **Switching on the page only changes the view**: viewing the other Model Group triggers no network request, does not change config, and does not swap the Target and Predicted roles.
3. **One probe refreshes both groups**: a single Antigravity quota response is parsed and saved for both `gemini` and `claude_gpt`; which group the page is viewing does not change what is collected.
4. **Predicted is always read-only**: the Control Model Group produces the Target; the Predicted Model Group only shows the result "if it became the control group" and must never reach Host write-back.
5. **The two groups' data are independent**: when one group fails or has no data it shows unknown, and must not borrow, copy or relabel the other group's data.
6. **All use cases output the same shape**: Production, Probe, SyncHost, Reset and the startup fallback all return a complete dual-group Snapshot with explicit roles.
7. **Saving config does not write the Host immediately**: after switching the Control Model Group, the new write-back authority takes effect at the next Manual Apply or Auto Apply that obtains Fresh Evidence.
8. **Fresh Evidence authority stays separate**: the projection only consumes evidence already classified as current or historical, and does not reinterpret freshness fields.

#### Speculative complexity explicitly rejected

The following are not core requirements this round:

- a generic framework for a third or arbitrary number of Model Groups;
- new projection strategies, switching strategies or automatic takeover config;
- absorbing network probing, State Store persistence or Host write-back into the projection module;
- new databases, cache layers, background refresh jobs or management APIs;
- redesigning history records, the Management page or the public Snapshot schema.

---

## 4. Organize management page assets by user feature

- **Status**: `Accepted`
- **Recommendation**: `Worth exploring`
- **Dependency category**: `in-process`
- **Suggested priority**: `P3`
- **Files involved**:
  - `internal/management/templates.go`
  - `internal/management/template_scripts.go`
  - `internal/management/template_styles.go`
  - `internal/management/template_*.go`
  - `internal/management/templates_test.go`

### Current pain points

Recent commits split what used to be thousands of lines of embedded CSS and JavaScript strings into about 27 files. That fixed single-file length, but the split mostly follows technology type: HTML, CSS, JavaScript and i18n each live separately.

So one user feature still spans several modules. For example, knowledge about the config center is spread across:

```text
TemplateConfig
templateScriptConfig
templateStyleConfig
templateScriptI18N
shared controls
shared modal
boot initialization
```

The files are shorter, but making one feature change still requires rebuilding the full context from several places.

### Concrete evidence

#### 4.1 JavaScript fragments depend on manual concatenation order

`template_scripts.go` explicitly concatenates several private fragments in dependency order. The fragments form an implicit interface through global functions, global variables and DOM IDs.

Inline calls in the HTML, such as:

```html
onclick="openKeyModal()"
onclick="toggleLanguage()"
onclick="toggleTheme()"
onclick="switchTab('overview')"
```

require the related functions to exist in the final global scope and to be initialized in the right order. This interface is not expressed explicitly by types or an assembly structure.

#### 4.2 CSS fragments depend on manual cascade order

`template_styles.go` concatenates the token, shell, controls, overview, history, diagnostics, config, overlays and responsive fragments in a fixed order.

Changing the order can change which selectors win, so file order itself becomes implicit knowledge callers must maintain.

#### 4.3 One user feature spans several technology files

Features such as config, diagnostics and overview each span markup, style, behaviour and language resources. To understand one button, a maintainer or AI has to open several files and rebuild the global dependencies.

#### 4.4 Tests mostly check source strings

`templates_test.go` currently relies mainly on:

- `strings.Contains` to check for elements or source fragments;
- regexes to check external URLs;
- regexes to check CSS tokens;
- specific JavaScript source substrings to check behavior requirements.

These tests catch missing resources, but cannot reliably verify:

- whether the JavaScript parses and runs;
- whether global functions exist when called;
- whether DOM IDs and selectors match;
- whether i18n keys are complete;
- whether user actions produce the expected DOM behavior;
- whether changing fragment order breaks runtime behavior.

The tests bypass the final resource interface and depend directly on implementation text. They can fail when the wording changes but behavior does not, and keep passing when the source is still there but the behavior is broken.

### Why deepening is needed

The real unit of maintenance for the management page is the user feature, not "all CSS" or "all JavaScript".

Ideal locality puts a feature's markup, style, behaviour and language resources next to each other, while hiding assembly order and shared runtime details inside the implementation.

The current split improved file length, but most of the 27 fragments have no independent behavior or stable interface. They only work through final string concatenation and global naming, which makes them shallow modules.

### Deletion test

Deleting any single fragment file and moving its string back into the aggregate file neither concentrates nor removes behavioral complexity, which shows that individual fragments have no significant depth.

If deep modules were formed per user feature, deleting one feature module would scatter that feature's markup, style, behaviour, copy, assembly and verification knowledge back across several places. This shows that grouping by user feature gives stronger locality.

### Deepening direction

Keep the small final `StatusHTML` interface, but organize the management page implementation by user feature, conceptually for example:

```text
Overview feature
History feature
Diagnostics feature
Config feature
Modal feature
```

Each feature gathers its own markup, style, behaviour and i18n. Final assembly still produces a single native, self-contained `StatusHTML`.

The test surface should gradually move from source substrings to executable behavior of the final resource, including script parsing, DOM dependencies, i18n completeness and core user actions.

The concrete file layout and interface will be decided in a later design phase. This proposal does not advocate merging all assets back into one giant file.

### Goals

- Keep a single user feature change local.
- Hide assembly knowledge such as global functions, DOM IDs, script order and CSS order.
- Reduce the number of files a maintainer or AI must cross to understand one feature.
- Keep the small `StatusHTML` interface and self-contained delivery.
- Make final browser behavior the main test surface.

### Expected results

- Changes to the config, diagnostics or overview features have their related knowledge in one module.
- Less leakage of JavaScript global names and assembly order.
- CSS cascade order is more centralized and easier to verify.
- Deleting a feature deletes its markup, style, behaviour and copy together.
- Tests catch script parse errors, DOM dependency errors and real interaction failures.
- `StatusHTML` remains the management page's only external interface.
- The management page still meets the CPA Host's embedded loading requirements.

### Risks and constraints

- Must stay pure native, self-contained, with zero external CDN at runtime.
- Must remain compatible with the CPA Host's dark and light themes.
- Must not introduce a frontend build chain as a CPA runtime dependency.
- Should not fall back to a single giant string thousands of lines long.
- Code shared between feature modules must pass the deletion test, to avoid creating new shallow modules.
- Browser behavior tests should keep their complexity in check, so the test toolchain does not outgrow the management page implementation itself.

### ADR relationship

This proposal does not conflict with existing ADRs. Implementation must strictly keep the Web UI constraints in AGENTS.md:

- pure native frontend implementation;
- self-contained;
- zero external CDN;
- full adaptation to the CPA Host's dark and light themes.

This round produced no new, hard-to-reverse business trade-off, so no new ADR is added. Organizing assets by user feature is an internal architectural consolidation under the existing management page contract and introduces no new domain terms.

### Decisions confirmed in design interrogation (2026-08-22)

The following scope was confirmed via `/grill-with-docs`; later design and implementation should stay focused on it:

1. **Behavior-equivalent reorganization**: this item only improves locality, assembly reliability and the test surface of management page assets; it does not change user-visible layout, interaction flows, the meaning of copy, or Management API behavior.
2. **The single external interface stays the same**: `StatusHTML` remains the management page's only external resource, and the `/status` route, auth and request semantics stay compatible.
3. **Assets are owned by user feature**: Overview, History, Diagnostics, Config and Help each own their markup, style, behaviour and i18n; the Shell only owns truly cross-feature navigation, theme, language, notifications, common controls and modal foundations.
4. **Runtime constraints stay the same**: page resources remain pure native, self-contained, with zero remote URLs/CDN, compliant with a strict CSP, and add no frontend build chain to the CPA runtime.
5. **Product capabilities stay the same**: keep the five feature areas, the dual Model Group view, Probe, Apply, Reset, Diff confirmation, Chinese/English switching, CPA Host theme adaptation, independent light/dark selection, responsive layout, and both embedded and direct browser access.
6. **The final resource is the main test surface**: verify that the final script parses, that inline handlers and DOM dependencies are closed, and that i18n references are complete; add behavior tests only for a few high-value flows such as saving config and confirming Apply.
7. **Shared code must prove its leverage**: only capabilities genuinely reused by several user features go into the Shell; a shared module must not become a shallow pass-through layer that re-centralizes all code.

#### Explicit scope boundaries

This round introduces no new UI design, frontend framework, bundler, npm app, CSS-in-JS, plugin system or generic feature registration framework, and adds no browser version matrix. Concrete file layout, internal interfaces, selector management and test tooling are implementation decisions.

---

## 5. Recommended order and how the proposals combine

### Recommended implementation order

| Order | Proposal | Strength | Main reason |
|---|---|---|---|
| 1 | Unify Host transitions and execution records | Strong | Directly affects Host state, failure reporting and audit trustworthiness |
| 2 | Consolidate ownership of Fresh Evidence | Strong | Decides whether evidence is eligible for Host write-back; it is the domain safety gate |
| 3 | Deepen the dual Model Group planning projection | Strong | Removes duplicated domain orchestration and creates a pure-computation test surface |
| 4 | Organize management page assets by user feature | Worth exploring | Mainly improves maintenance locality and behavior testing; current production risk is low |

### How the three scheduling candidates combine

The first three items form an execution chain with clear responsibilities:

```text
Evidence module
  │  outputs only current-round or historical evidence with explicit purpose
  ▼
Dual Model Group projection module
  │  uniform projection of control group and predicted group
  ▼
priority.Planner
  │  side-effect-free pure function, outputs only an immutable Plan
  ▼
Host transition / execution record module
  │  executes gating, Host transition and redacted recording
  ▼
CPA Host + persisted execution history
```

This combination keeps the existing architectural ownership:

- the Domain Core owns evidence semantics and pure planning;
- the Apply Layer owns Host transitions and redacted execution records;
- the Runtime owns scheduling use cases, the Ticker and `runMu` single-flight;
- Management and the CGO adapter stay thin pass-throughs.

### Questions to settle before implementation

Before any candidate moves into concrete design or implementation, settle for each:

1. Which behaviors belong to the target module's external interface, and which should only be internal seams.
2. Whether there are at least two real adapters, or only an in-process implementation is needed.
3. Which existing tests should move to the new interface test surface, and which shallow-module tests should be deleted.
4. The exact semantics of partial Host transitions, failure evidence and predicted Plans.
5. Whether `CONTEXT.md` domain terms need to be added or updated.
6. Whether any load-bearing decision needs to be recorded as a new ADR.

Until these questions are discussed, this document defines no concrete interface and does not authorize implementation code changes.
