package management

// historyMarkup contains the HTML owned by the History feature asset.
const historyMarkup = `
        <section id="panelHistory" hidden>
            <div class="command-wrap">
                <div class="commandbar" aria-label="History actions">
                    <div class="command-state">
                        <span class="command-label">Run history</span>
                        <span id="historyCountBadge" class="health-badge"><span class="status-dot neutral"></span><span>No runs yet</span></span>
                    </div>
                    <div class="command-actions">
                        <button type="button" class="button secondary compact" onclick="fetchDiagnostics()" title="Reload run history">` + iconRefresh + `<span class="label-long">Refresh</span></button>
                    </div>
                </div>
            </div>
            <div class="editor">
                <header class="editor-head">
                    <div class="editor-title-wrap">
                        <p class="editor-kicker">Execution history</p>
                        <h1>Recent runs</h1>
                        <p class="editor-subtitle">The last 10 probes, applies, resets and automatic schedule runs, with the details recorded for each.</p>
                    </div>
                </header>
                <section class="card">
                    <div id="historyList" class="history-list">
                        <div class="empty-state">No execution history yet</div>
                    </div>
                </section>
            </div>
        </section>
`

var historyPageAsset = managementPageAsset{
	name:    "History",
	markup:  historyMarkup,
	styles:  templateStyleHistory,
	scripts: templateScriptHistory + templateScriptHistoryDetails,
}
