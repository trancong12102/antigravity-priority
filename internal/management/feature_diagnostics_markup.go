package management

// diagnosticsMarkup contains the HTML owned by the Diagnostics feature asset.
const diagnosticsMarkup = `
        <section id="panelDiagnostics" hidden>
            <div class="command-wrap">
                <div class="commandbar" aria-label="Diagnostics actions">
                    <div class="command-state">
                        <span class="command-label">Diagnostics</span>
                        <span class="health-badge"><span class="status-dot neutral"></span><span>Secrets redacted</span></span>
                    </div>
                    <div class="command-actions">
                        <button type="button" class="icon-button" onclick="fetchDiagnostics()" aria-label="Reload diagnostics" title="Reload diagnostics">` + iconRefresh + `</button>
                        <button type="button" class="button secondary compact" id="btnCopyDiagnostics" onclick="copyDiagnosticsJSON()" title="Copy diagnostics JSON">` + iconCopy + `<span class="label-long">Copy diagnostics JSON</span></button>
                    </div>
                </div>
            </div>
            <div class="editor">
                <header class="editor-head">
                    <div class="editor-title-wrap">
                        <p class="editor-kicker">System diagnostics</p>
                        <h1>Scheduler and health</h1>
                        <p class="editor-subtitle">How the scheduler is running, which credentials sit in 429 cooldown, and what the last write to CPA did.</p>
                    </div>
                </header>

                <section class="stat-grid three diag-kpi-grid" aria-label="Diagnostics summary">
                    <article class="stat-card">
                        <div class="stat-head"><span>Scheduler engine</span><span id="diagSchedBadge" class="badge badge-subtle">-</span></div>
                        <strong id="diagSchedInterval">-</strong>
                        <small id="diagSchedCountdown">-</small>
                    </article>
                    <article class="stat-card">
                        <div class="stat-head"><span>429 circuit breakers</span><span id="diagCooldownBadge" class="badge badge-success">Healthy</span></div>
                        <strong id="diagCooldownCount">0</strong>
                        <small id="diagCooldownSub">Auto demotion to -1 ready</small>
                    </article>
                    <article class="stat-card">
                        <div class="stat-head"><span>Last apply health</span><span id="diagApplyBadge" class="badge badge-subtle">-</span></div>
                        <strong id="diagApplyStats">-</strong>
                        <small id="diagApplyTime">-</small>
                    </article>
                </section>

                <section class="card" id="diagSectionScheduler">
                    <div class="card-head"><div><h2>Scheduling details</h2><p>Interval, background worker and the active time window.</p></div></div>
                    <div class="setting-row">
                        <div class="setting-copy"><strong>Base interval</strong></div>
                        <div id="diagSchedDetailInterval" class="setting-control mono">-</div>
                    </div>
                    <div class="setting-row">
                        <div class="setting-copy"><strong>Background worker</strong></div>
                        <div id="diagSchedDetailWorker" class="setting-control">-</div>
                    </div>
                    <div class="setting-row">
                        <div class="setting-copy"><strong>Last auto run</strong></div>
                        <div id="diagSchedDetailLastRun" class="setting-control mono">-</div>
                    </div>
                    <div class="setting-row">
                        <div class="setting-copy"><strong>Next run (estimated)</strong></div>
                        <div id="diagSchedDetailNextRun" class="setting-control mono">-</div>
                    </div>
                    <div class="setting-row">
                        <div class="setting-copy"><strong>Window policy</strong><span id="diagSchedDetailWindowText" class="mono">-</span></div>
                        <div class="setting-control"><span id="diagSchedDetailWindowBadge" class="badge badge-subtle">-</span></div>
                    </div>
                </section>

                <section class="card" id="diagSectionCooldown">
                    <div class="card-head"><div><h2>429 cooldowns</h2><p>Credentials temporarily demoted to -1 after rate limiting.</p></div></div>
                    <div id="diagCooldownContent" class="card-body">
                        <div class="default-notice">` + iconCheck + `<span>No credentials in 429 cooldown; all accounts are operating normally.</span></div>
                    </div>
                </section>

                <section class="card" id="diagSectionAudit">
                    <div class="card-head"><div><h2>Latest write and redacted audit</h2><p>The most recent apply recorded by the runtime.</p></div></div>
                    <div class="card-body diag-audit">
                        <p class="command-label">AUDIT STREAM</p>
                        <div id="diagAuditText" class="diag-audit-text">-</div>
                        <div id="diagAuditPills" class="diag-pills-row"></div>
                    </div>
                </section>
            </div>
        </section>
`

var diagnosticsPageAsset = managementPageAsset{
	name:    "Diagnostics",
	markup:  diagnosticsMarkup,
	styles:  templateStyleDiagnostics,
	scripts: templateScriptDiagnostics,
}
