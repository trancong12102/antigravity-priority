package management

// helpMarkup contains the HTML owned by the Help feature asset.
const helpMarkup = `
        <section id="panelHelp" hidden>
            <div class="editor">
                <header class="editor-head">
                    <div class="editor-title-wrap">
                        <p class="editor-kicker">Help</p>
                        <h1>How priorities are computed</h1>
                        <p class="editor-subtitle">Built for Google Antigravity in CLIProxyAPI: two quota windows drive adaptive priority scheduling and write-back to CPA.</p>
                    </div>
                </header>
                <section class="card">
                    <div class="card-head"><div><h2>Mechanics</h2><p>What the scheduler watches and how it reacts.</p></div></div>
                    <div class="info-callout">
                        <span class="callout-icon">` + iconBolt + `</span>
                        <div><strong>5-hour window: soft demotion</strong><p>Monitors burst quota and its reset countdown. When the short window runs out the credential is only demoted (priority=-1, disabled=false) and recovers automatically when the window resets, which avoids 429 rate limits.</p></div>
                    </div>
                    <div class="info-callout">
                        <span class="callout-icon">` + iconCalendar + `</span>
                        <div><strong>7-day window: global control</strong><p>Tracks the remaining weekly balance and reset progress. Weekly depletion triggers a hard disable (disabled=true) so a credential is not exhausted early in the cycle.</p></div>
                    </div>
                    <div class="info-callout">
                        <span class="callout-icon">` + iconTrend + `</span>
                        <div><strong>Adaptive burn-rate learning (C_cycle)</strong><p>Learns each credential's real consumption from probe deltas with EMA smoothing, so the safe reserve is computed automatically instead of configured by hand.</p></div>
                    </div>
                    <div class="info-callout">
                        <span class="callout-icon">` + iconScale + `</span>
                        <div><strong>Weekly urgency and load balancing</strong><p>Quantifies quota pressure per unit of time and groups credentials within the tolerance into equal priority tiers, so traffic rotates between them.</p></div>
                    </div>
                    <div class="info-callout">
                        <span class="callout-icon">` + iconArrowUp + `</span>
                        <div><strong>Dynamic boost horizon</strong><p>Credentials with plenty of quota left close to their weekly reset move into the top tier (900-999), so large balances are spent instead of wasted.</p></div>
                    </div>
                    <div class="info-callout">
                        <span class="callout-icon">` + iconShield + `</span>
                        <div><strong>429 cooldown and self-healing</strong><p>A 429 from Google demotes the credential to -1 for the cooldown period; the next scheduled run probes it again and restores it without manual work.</p></div>
                    </div>
                </section>
            </div>
        </section>
`

var helpPageAsset = managementPageAsset{
	name:   "Help",
	markup: helpMarkup,
}
