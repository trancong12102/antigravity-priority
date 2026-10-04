package management

// overviewMarkup contains the HTML owned by the Overview feature asset.
const overviewMarkup = `
        <section id="panelOverview">
            <div class="command-wrap">
                <div class="commandbar" aria-label="Overview actions">
                    <div class="command-state">
                        <span class="command-label">Auto schedule</span>
                        <button type="button" id="scheduleStatusBadge" class="health-badge" onclick="toggleSchedulePause()" title="Click to toggle pause">
                            <span class="status-dot neutral"></span><span id="scheduleStatusText">Loading</span>
                        </button>
                    </div>
                    <div class="command-actions">
                        <select id="modelGroupSelect" class="select-compact" aria-label="Model group" onchange="onModelGroupChange()">
                            <option value="gemini">Gemini models</option>
                            <option value="claude_gpt">Claude &amp; GPT models</option>
                        </select>
                        <button type="button" class="icon-button" id="btnRefresh" onclick="refreshDashboard({ withSync: true, notifySync: true })" aria-label="Refresh from CPA" title="Refresh from CPA">` + iconRefresh + `</button>
                        <span class="separator" aria-hidden="true"></span>
                        <button type="button" class="button secondary compact" id="btnProbe" onclick="triggerProbe()" title="Probe quota for every credential">` + iconProbe + `<span class="label-long">Fetch quota</span></button>
                        <button type="button" class="button danger compact" id="btnReset" onclick="triggerReset()" title="Reset every priority to unset">` + iconReset + `<span class="label-long">Reset priority</span></button>
                        <button type="button" class="button primary compact" id="btnApply" onclick="triggerApplyWithConfirm()" title="Write target priorities to CPA">` + iconBolt + `<span class="label-long">Apply</span></button>
                    </div>
                </div>
            </div>

            <div class="editor">
                <header class="editor-head">
                    <div class="editor-title-wrap">
                        <p class="editor-kicker">Priority overview</p>
                        <h1>Antigravity credentials</h1>
                        <p class="editor-subtitle">Live 5h and 7d quota for every credential, with the priority the scheduler would write back to CPA.</p>
                    </div>
                </header>

                <section class="stat-grid" aria-label="Credential statistics">
                    <article class="stat-card">
                        <span>Total credentials</span>
                        <strong id="valTotalCreds">0</strong>
                        <small id="valTotalDesc">0 active</small>
                    </article>
                    <article class="stat-card">
                        <span>Boosted</span>
                        <strong id="valBoosted">0</strong>
                        <small>High quota abundance</small>
                    </article>
                    <article class="stat-card">
                        <span>Depleted / down</span>
                        <strong id="valDepleted">0</strong>
                        <small>Awaiting window reset</small>
                    </article>
                    <article class="stat-card">
                        <span>Latest audit</span>
                        <strong id="valLastAudit" class="text-value">-</strong>
                        <small id="valNextProbe">-</small>
                    </article>
                </section>

                <div id="credentialsContainer" class="credentials-grid">
                    <div class="empty-state">Loading credentials and quota…</div>
                </div>
            </div>
        </section>
`

var overviewPageAsset = managementPageAsset{
	name:    "Overview",
	markup:  overviewMarkup + overviewSamplesMarkup,
	styles:  templateStyleOverviewCards,
	scripts: templateScriptOverviewActionsCore + templateScriptOverviewReset + templateScriptOverviewRender + templateScriptProbeSchedule + templateScriptOverviewSamples,
}
