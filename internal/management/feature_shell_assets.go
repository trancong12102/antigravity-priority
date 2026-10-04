package management

// managementPageShellMarkup contains the session gate and the document-level
// chrome (top bar and page navigation) shared by every feature panel.
const managementPageShellMarkup = `
</head>
<body>
    <main id="authGate" class="auth-shell">
        <section class="auth-card" aria-labelledby="authTitle">
            <div id="sessionIcon" class="brand-mark" aria-hidden="true">` + iconLock + `</div>
            <p class="eyebrow">CPAMC SESSION</p>
            <h1 id="authTitle">Connecting to management session</h1>
            <p id="authMessage" class="auth-copy" role="status" aria-live="polite">Reading the current CPAMC connection (read-only). No need to enter the Management Key again.</p>
            <button id="retrySessionButton" class="button secondary wide" type="button" onclick="startSession()" hidden>Retry session</button>
            <div class="security-note">
                <span class="status-dot neutral" aria-hidden="true"></span>
                <span>This page only reads the same-origin session CPAMC already saved; it never copies or persists the Management Key.</span>
            </div>
        </section>
    </main>

    <div id="app" class="app" hidden>
        <header class="topbar">
            <div class="brand">
                <div class="brand-mark small" aria-hidden="true">` + iconBrand + `</div>
                <div class="brand-copy">
                    <strong>Antigravity Priority</strong>
                    <small>CPA plugin · <span class="version-badge">v1.2.18-en.1</span></small>
                </div>
            </div>
            <p class="topbar-context">Quota-aware priorities for Antigravity credentials</p>
        </header>

        <div class="workspace">
            <aside class="sidebar" aria-label="Page navigation">
                <div class="sidebar-head">
                    <div>
                        <p class="section-label">WORKSPACE</p>
                        <p class="sidebar-count">Antigravity credentials</p>
                    </div>
                </div>
                <nav class="page-nav" aria-label="Dashboard sections">
                    <button type="button" class="nav-item" data-tab="overview" aria-current="page" onclick="switchTab('overview')">
                        <span class="nav-icon">` + iconOverview + `</span>
                        <span class="nav-copy"><strong>Overview</strong><span>Quota meters and priorities</span></span>
                    </button>
                    <button type="button" class="nav-item" data-tab="history" aria-current="false" onclick="switchTab('history')">
                        <span class="nav-icon">` + iconHistory + `</span>
                        <span class="nav-copy"><strong>Run history</strong><span>Last 10 probes and applies</span></span>
                    </button>
                    <button type="button" class="nav-item" data-tab="diagnostics" aria-current="false" onclick="switchTab('diagnostics')">
                        <span class="nav-icon">` + iconPulse + `</span>
                        <span class="nav-copy"><strong>Diagnostics</strong><span>Scheduler, 429 cooldowns, audit</span></span>
                    </button>
                    <button type="button" class="nav-item" data-tab="config" aria-current="false" onclick="switchTab('config')">
                        <span class="nav-icon">` + iconSliders + `</span>
                        <span class="nav-copy"><strong>Configuration</strong><span>Runtime settings</span></span>
                    </button>
                    <button type="button" class="nav-item" data-tab="help" aria-current="false" onclick="switchTab('help')">
                        <span class="nav-icon">` + iconHelp + `</span>
                        <span class="nav-copy"><strong>Help</strong><span>How priorities are computed</span></span>
                    </button>
                </nav>
                <div class="sidebar-note">Signed in through the CPAMC session on this origin. Log out of CPAMC and this page locks again.</div>
            </aside>

            <main class="main-stage">`

const managementPageShellMarkupEnd = `
            </main>
        </div>
    </div>
`

const managementPageShellStyles = templateStyleTokens +
	templateStyleShell +
	templateStyleControls

const managementPageShellStylesTail = templateStyleOverlays + templateStyleResponsive

const managementPageShellScripts = templateScriptPrelude +
	templateScriptIcons +
	templateScriptShellCore +
	templateScriptControls +
	templateScriptModals +
	templateScriptLiveUI
