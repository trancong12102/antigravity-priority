package management

// configMarkup contains the HTML owned by the Config feature asset.
const configMarkup = `
        <section id="panelConfig" hidden>
            <div class="command-wrap">
                <div class="commandbar" aria-label="Configuration actions">
                    <div class="command-state">
                        <span class="command-label">Runtime config</span>
                        <span id="configStateBadge" class="health-badge"><span class="status-dot neutral"></span><span>Loading</span></span>
                    </div>
                    <div class="command-actions">
                        <button type="button" class="button secondary compact" onclick="resetDynamicConfigToDefaults()" id="btnResetConfig" title="Reset every runtime setting to its default">` + iconReset + `<span class="label-long">Reset to defaults</span></button>
                        <button type="button" class="button primary compact" onclick="saveDynamicConfig()" id="btnSaveConfig" title="Save and apply immediately" disabled>` + iconSave + `<span class="label-long">Save &amp; apply</span></button>
                    </div>
                </div>
            </div>
            <div class="editor">
                <header class="editor-head">
                    <div class="editor-title-wrap">
                        <p class="editor-kicker">Configuration</p>
                        <h1>Runtime settings</h1>
                        <p class="editor-subtitle">Changes take effect immediately after saving; nothing here needs a CPA restart.</p>
                    </div>
                </header>

                <section class="card">
                    <div class="card-head"><div><h2>Auto scheduling and time window</h2><p>When the scheduler probes quota and writes priorities back.</p></div></div>
                    <div class="setting-row">
                        <div class="setting-copy"><strong>Auto periodic scheduling</strong><span>Periodically probe and update host credential priorities</span></div>
                        <div class="setting-control">
                            <label class="switch">
                                <input type="checkbox" id="cfgAutoApply">
                                <span class="switch-track" aria-hidden="true"></span>
                                <span id="cfgAutoApplyStatusText">Disabled</span>
                            </label>
                        </div>
                    </div>
                    <div class="setting-row">
                        <div class="setting-copy"><strong>Scheduling interval</strong><span>How often probing and scheduling run (e.g. 15m, 30m, 1h)</span></div>
                        <div class="setting-control">
                            <select id="cfgIntervalSelect" aria-label="Scheduling interval" onchange="onIntervalSelectChange()">
                                <option value="5m">5 minutes</option>
                                <option value="15m" selected>15 minutes (recommended)</option>
                                <option value="30m">30 minutes</option>
                                <option value="1h">1 hour</option>
                                <option value="custom">Custom</option>
                            </select>
                            <input type="text" id="cfgIntervalCustom" class="time-input" placeholder="15m" aria-label="Custom interval" hidden>
                        </div>
                    </div>
                    <div class="setting-row">
                        <div class="setting-copy"><strong>Primary model group</strong><span>Model group used as the main basis for write-back</span></div>
                        <div class="setting-control">
                            <select id="cfgModelGroup" aria-label="Primary model group">
                                <option value="gemini">Gemini models</option>
                                <option value="claude_gpt">Claude &amp; GPT models</option>
                            </select>
                        </div>
                    </div>
                    <div class="setting-row">
                        <div class="setting-copy"><strong>Active schedule window</strong><span>Only schedule inside a daily window; crossing midnight works (e.g. 22:00-06:00)</span></div>
                        <div class="setting-control wide">
                            <label class="switch">
                                <input type="checkbox" id="cfgWindowEnabled" onchange="onWindowEnabledChange()">
                                <span class="switch-track" aria-hidden="true"></span>
                                <span>Window only</span>
                            </label>
                            <div id="cfgWindowInputs" class="window-inputs">
                                <input type="text" id="cfgWindowStart" class="time-input" placeholder="09:00" aria-label="Window start">
                                <span class="muted">to</span>
                                <input type="text" id="cfgWindowEnd" class="time-input" placeholder="23:00" aria-label="Window end">
                            </div>
                        </div>
                    </div>
                </section>

                <section class="card">
                    <div class="card-head"><div><h2>Probing and performance</h2><p>Concurrency, write thresholds and the adaptive sample window.</p></div></div>
                    <div class="setting-row">
                        <div class="setting-copy"><strong>Max probe concurrency</strong><span>Maximum concurrent goroutines for Google API probes (1-32)</span></div>
                        <div class="setting-control"><input type="number" id="cfgMaxConcurrency" min="1" max="32" value="6"></div>
                    </div>
                    <div class="setting-row">
                        <div class="setting-copy"><strong>Priority change threshold</strong><span>Minimum priority delta before writing back to the host (0-100)</span></div>
                        <div class="setting-control"><input type="number" id="cfgMinChange" min="0" max="100" value="1"></div>
                    </div>
                    <div class="setting-row">
                        <div class="setting-copy"><strong>Urgency bucket tolerance</strong><span>Credentials within this tolerance share the same priority (0.00-0.50)</span></div>
                        <div class="setting-control"><input type="number" id="cfgUrgencyTolerance" min="0.00" max="0.50" step="0.01" value="0.05"></div>
                    </div>
                    <div class="setting-row">
                        <div class="setting-copy"><strong>Adaptive sample capacity</strong><span>Samples kept in the sliding window for burn-rate estimation (2-30)</span></div>
                        <div class="setting-control"><input type="number" id="cfgSampleCapacity" min="2" max="30" value="6"></div>
                    </div>
                    <div class="setting-row">
                        <div class="setting-copy"><strong>429 cooldown (minutes)</strong><span>How long a credential stays demoted to -1 after a 429 (1-1440)</span></div>
                        <div class="setting-control"><input type="number" id="cfgCooldownMinutes" min="1" max="1440" value="5"></div>
                    </div>
                </section>

                <section class="card">
                    <div class="card-head"><div><h2>Priority scoring rules</h2><p>Starting priorities for each credential tier.</p></div></div>
                    <div class="setting-row">
                        <div class="setting-copy"><strong>Dynamic boost start priority</strong><span>Base priority for top-tier credentials with abundant quota (1-999)</span></div>
                        <div class="setting-control"><input type="number" id="cfgBoostStartPriority" min="1" max="999" value="999"></div>
                    </div>
                    <div class="setting-row">
                        <div class="setting-copy"><strong>Normal start priority</strong><span>Base priority for regular healthy credentials (1-999)</span></div>
                        <div class="setting-control"><input type="number" id="cfgNormalStartPriority" min="1" max="999" value="100"></div>
                    </div>
                </section>
            </div>
        </section>
`

var configPageAsset = managementPageAsset{
	name:    "Config",
	markup:  configMarkup,
	styles:  templateStyleConfig,
	scripts: templateScriptConfig,
}
