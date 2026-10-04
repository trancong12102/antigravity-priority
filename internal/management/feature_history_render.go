package management

const templateScriptHistory = `        function formatHistoryKind(entry) {
			var kind = entry && entry.kind ? entry.kind : "";
			var trigger = entry && entry.trigger ? entry.trigger : "";
			var k = kind.toLowerCase();
            if (k === "auto_apply" || trigger.toLowerCase() === "auto_apply") return "Auto schedule";
            if (k === "apply" || k === "manual_apply") return "Apply";
            if (k === "probe") return "Probe";
            if (k === "reset") return "Reset";
            var label = (kind || "run").replace(/_/g, " ");
            return label.charAt(0).toUpperCase() + label.slice(1);
        }

        function renderHistory() {
            const list = document.getElementById("historyList");
            if (!list) return;

            const entries = (latestDiagnostics && latestDiagnostics.run_history) || [];
            const countBadge = document.getElementById("historyCountBadge");
            if (countBadge) {
                countBadge.innerHTML = "<span class=\"status-dot" + (entries.length ? "" : " neutral") + "\"></span><span>" +
                    (entries.length ? entries.length + (entries.length === 1 ? " run" : " runs") + " recorded" : "No runs yet") + "</span>";
            }
            list.innerHTML = "";
            if (entries.length === 0) {
                list.innerHTML = "<div class=\"empty-state\">No execution history yet</div>";
                return;
            }

            entries.forEach(function(entry, idx) {
                const item = document.createElement("div");
                item.className = "history-item";
                const dateStr = formatDateTime(entry.at);
				const kindText = formatHistoryKind(entry);
                const msg = entry.message || "";
				const isAutoSchedule = entry.kind === "auto_apply" || entry.trigger === "auto_apply";
				const hasSnapshot = entry.kind === "probe" ? Boolean(entry.probe_round_id) : (isAutoSchedule ? Boolean(entry.probe_round_id || (entry.snapshot && (entry.snapshot.items || entry.snapshot.changes))) : (entry.snapshot && (entry.snapshot.items || entry.snapshot.changes)));

                item.innerHTML =
                    "<div class=\"history-head\">" +
                        "<div class=\"history-when\">" +
                            "<span class=\"badge\">" + escapeHTML(kindText) + "</span>" +
                            "<strong>" + escapeHTML(dateStr) + "</strong>" +
                        "</div>" +
                        "<div class=\"history-stats\">" +
                            "<span class=\"badge badge-success\">Succeeded " + (Number(entry.succeeded) || 0) + "</span>" +
                            "<span class=\"badge" + (entry.failed ? " badge-danger" : " badge-subtle") + "\">Failed " + (Number(entry.failed) || 0) + "</span>" +
                            "<span class=\"badge badge-subtle\">Skipped " + (Number(entry.skipped) || 0) + "</span>" +
                            (hasSnapshot ? "<button type=\"button\" class=\"button secondary mini\" data-history-index=\"" + idx + "\">" + ICONS.search + "<span>Details</span></button>" : "") +
                        "</div>" +
                    "</div>" +
                    (msg ? "<div class=\"history-message\">" + escapeHTML(msg) + "</div>" : "");
                list.appendChild(item);
            });
        }

`
