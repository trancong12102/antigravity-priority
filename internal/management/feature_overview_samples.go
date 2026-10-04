package management

const overviewSamplesMarkup = `
    <!-- Quota Samples History Modal -->
    <div id="samplesModal" class="modal-backdrop" hidden>
        <div class="modal medium" role="dialog" aria-modal="true" aria-labelledby="samplesModalTitle">
            <div class="modal-head">
                <div>
                    <h2 id="samplesModalTitle">Adaptive quota samples</h2>
                    <p id="samplesModalSubtitle" class="mono"></p>
                </div>
                <button type="button" class="icon-button" onclick="closeSamplesModal()" aria-label="Close">` + iconClose + `</button>
            </div>
            <div id="samplesModalBody" class="modal-body"></div>
            <div class="modal-foot">
                <button type="button" class="button secondary compact" onclick="closeSamplesModal()">Close</button>
            </div>
        </div>
    </div>
`

const templateScriptOverviewSamples = `        async function openSamplesModal(authIndex, email) {
            const modal = document.getElementById("samplesModal");
            const title = document.getElementById("samplesModalTitle");
            const sub = document.getElementById("samplesModalSubtitle");
            const body = document.getElementById("samplesModalBody");
            if (!modal || !body) return;

            title.textContent = "Adaptive quota samples";
            sub.textContent = email || "";
            body.innerHTML = "<div class=\"empty-state\">Loading samples…</div>";
            modal.hidden = false;

            try {
                const path = SAMPLES_PATH + "?auth_index=" + encodeURIComponent(authIndex);
                const data = await apiFetch(path);
                const groups = (data && data.groups) || {};
                const groupSelect = document.getElementById("modelGroupSelect");
                const selectedGroup = (groupSelect && groupSelect.value) || "gemini";
                const selectedGroupData = groups[selectedGroup] || {};
                const samples = (Array.isArray(selectedGroupData.samples) ? selectedGroupData.samples : []).slice().sort(function(left, right) {
                    const leftTime = Date.parse(left && left.observed_at ? left.observed_at : "");
                    const rightTime = Date.parse(right && right.observed_at ? right.observed_at : "");
                    const leftValid = Number.isFinite(leftTime);
                    const rightValid = Number.isFinite(rightTime);
                    if (leftValid && rightValid && leftTime !== rightTime) return rightTime - leftTime;
                    if (leftValid !== rightValid) return rightValid ? 1 : -1;
                    return Number(right && right.sequence || 0) - Number(left && left.sequence || 0);
                });
                const groupLabel = selectedGroup === "claude_gpt" ? "Claude/GPT" : "Gemini";

                if (samples.length === 0) {
                    body.innerHTML = "<div class=\"empty-state\">No quota samples recorded yet</div>";
                    return;
                }

                function renderMiniMeter(val) {
                    if (val === undefined || val === null || val === "") return "<span class=\"muted\">-</span>";
                    const parsed = Number(val);
                    if (!Number.isFinite(parsed)) return "<span class=\"muted\">-</span>";
                    const num = Math.max(0, Math.min(100, Math.trunc(parsed)));
                    return "<div class=\"mini-meter\">" +
                        "<div class=\"meter-track\"><div class=\"meter-fill " + quotaFillClass(num) + "\" style=\"width:" + num + "%\"></div></div>" +
                        "<strong>" + num + "%</strong>" +
                    "</div>";
                }

                let html = "<table class=\"sample-table\">" +
                    "<thead><tr>" +
                        "<th>#</th>" +
                        "<th>Observed at</th>" +
                        "<th>5h window</th>" +
                        "<th>7d window</th>" +
                    "</tr></thead>";

                for (let i = 0; i < samples.length; i++) {
                    const sample = samples[i] || {};
                    html += "<tbody class=\"sample-group\">" +
                        "<tr>" +
                            "<td class=\"sample-index\">" + (i + 1) + "</td>" +
                            "<td class=\"sample-time\">" + escapeHTML(formatDateTime(sample.observed_at)) + "</td>" +
                            "<td>" + renderMiniMeter(sample.short_window_rem) + "</td>" +
                            "<td>" + renderMiniMeter(sample.long_window_rem) + "</td>" +
                        "</tr>" +
                    "</tbody>";
                }

                html += "</table>";
                body.innerHTML = html;
                sub.textContent = (email ? email + " · " : "") + groupLabel;
            } catch (err) {
                body.innerHTML = "<div class=\"empty-state error\">" + escapeHTML(err.message) + "</div>";
            }
        }

        function closeSamplesModal() {
            const modal = document.getElementById("samplesModal");
            if (modal) modal.hidden = true;
        }

`
