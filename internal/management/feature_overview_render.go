package management

const templateScriptOverviewRender = `        function quotaFillClass(percent) {
            if (percent <= 10) return "meter-fill-danger";
            if (percent <= 30) return "meter-fill-warning";
            return "meter-fill-healthy";
        }

        function renderQuotaMeter(label, percent, resetAt) {
            return "<div class=\"meter\">" +
                "<div class=\"meter-label-row\">" +
                    "<span>" + label + " <strong>" + percent + "%</strong></span>" +
                    "<span class=\"meter-countdown\" data-reset-time=\"" + escapeHTML(resetAt) + "\">" + escapeHTML(formatCountdown(resetAt)) + "</span>" +
                "</div>" +
                "<div class=\"meter-track\"><div class=\"meter-fill " + quotaFillClass(percent) + "\" style=\"width: " + percent + "%\"></div></div>" +
            "</div>";
        }

        function renderDashboard() {
            if (latestDiagnostics) {
                var auditEl = document.getElementById("valLastAudit");
                if (auditEl) {
                    auditEl.textContent = latestDiagnostics.latest_audit || "-";
                    auditEl.title = latestDiagnostics.latest_audit || "";
                }

                var sched = latestDiagnostics.scheduler || {};
                var isPaused = scheduleConfig ? Boolean(scheduleConfig.paused) : Boolean(sched.paused);
                var isSleeping = false;
                if (scheduleConfig && scheduleConfig.window_enabled && scheduleConfig.window_start && scheduleConfig.window_end) {
                    isSleeping = !isCurrentTimeInScheduleWindow(scheduleConfig.window_start, scheduleConfig.window_end);
                } else if (sched.window_enabled && sched.window_start && sched.window_end) {
                    isSleeping = !isCurrentTimeInScheduleWindow(sched.window_start, sched.window_end);
                }

                var el = document.getElementById("valNextProbe");
                if (el) {
                    if (!isAutoApplyEnabled()) {
                        el.textContent = "Auto schedule disabled";
                    } else if (isPaused) {
                        el.textContent = "Auto schedule paused";
                    } else if (isSleeping) {
                        el.textContent = "Sleeping outside the schedule window";
                    } else {
                        var nextRunAt = sched.next_run_at;
                        var nextStr = nextRunAt ? formatCountdown(nextRunAt) : (sched.next_wait || "-");
                        el.innerHTML = "Next run in <span class=\"meter-countdown\" data-scheduler-countdown=\"" + escapeHTML(nextRunAt) + "\">" + escapeHTML(nextStr) + "</span>";
                    }
                }
            }

            const container = document.getElementById("credentialsContainer");
            if (!container || !latestSnapshot) return;

            const groupSelect = document.getElementById("modelGroupSelect");
            const selectedGroup = groupSelect ? groupSelect.value : "gemini";
            const activeGroup = (latestSnapshot && latestSnapshot.active_model_group) || "gemini";
            const isPredictedView = selectedGroup !== activeGroup;
            const groupData = (latestSnapshot && latestSnapshot.groups && latestSnapshot.groups[selectedGroup]) || {};
            const items = groupData.items || [];
            const pendingAuthIndexes = new Set((groupData.changes || []).map(change => change.auth_index));
            let boostedCount = 0;
            let depletedCount = 0;
            let activeCount = 0;

            items.forEach(item => {
                if (item.is_boosted) boostedCount++;
                const isManuallyDisabled = item.reason && item.reason.indexOf("disabled on host") >= 0;
                if (item.target && item.target.disabled && !isManuallyDisabled) depletedCount++;
                if (item.target && !item.target.disabled) activeCount++;
            });

            document.getElementById("valTotalCreds").textContent = items.length;
            document.getElementById("valTotalDesc").textContent = activeCount + " active";
            document.getElementById("valBoosted").textContent = boostedCount;
            document.getElementById("valDepleted").textContent = depletedCount;

            container.innerHTML = "";
            if (items.length === 0) {
                container.innerHTML = "<div class=\"empty-state\">No Antigravity credentials found</div>";
                return;
            }

            items.forEach(item => {
                const card = document.createElement("article");
                card.className = "cred-card";

                const isBoosted = item.is_boosted;
                const r5hPercent = Math.max(0, Math.min(100, Math.round((item.r5h || 0) * 100)));
                const r7dPercent = Math.max(0, Math.min(100, Math.round((item.r7d || 0) * 100)));
                const urgency = (item.urgency || 0).toFixed(2);
                const burnRate = (item.cycle_burn_rate || 0).toFixed(2);

                let actualP = "-";
                if (item.current) {
                    if (item.current.disabled) actualP = "[Disabled]";
                    else if (item.current.priority_missing || item.current.priority === undefined || item.current.priority === null) actualP = "[Unset]";
                    else actualP = item.current.priority;
                } else if (item.priority_missing || item.priority === undefined || item.priority === null) {
                    actualP = "[Unset]";
                } else if (item.priority !== undefined && item.priority !== null) {
                    actualP = item.disabled ? "[Disabled]" : item.priority;
                }

                let targetP = "-";
                if (item.target) {
                    if (item.target.disabled) targetP = "[Disabled]";
                    else if (item.target.priority_missing || item.target.priority === undefined || item.target.priority === null) targetP = "[Unset]";
                    else targetP = item.target.priority;
                } else {
                    targetP = actualP;
                }

                let tagBadge = "";
                if (isPredictedView) {
                    tagBadge = "<span class=\"badge badge-predicted\">Predicted</span>";
                } else if (pendingAuthIndexes.has(item.auth_index)) {
                    tagBadge = "<span class=\"badge badge-pending\">Pending</span>";
                } else if (String(actualP) !== String(targetP)) {
                    tagBadge = "<span class=\"badge badge-warning\">Target changed</span>";
                }

                const probeFailed = item.reason && (item.reason.indexOf("probe failed") >= 0 || item.reason.indexOf("probe invalid") >= 0);
                const isManuallyDisabled = item.reason && item.reason.indexOf("disabled on host") >= 0;
                let statusBadge = "<span class=\"badge badge-success\">Active</span>";
                if (probeFailed) {
                    statusBadge = "<span class=\"badge badge-warning\">Probe failed</span>";
                } else if (isManuallyDisabled) {
                    statusBadge = "<span class=\"badge badge-subtle\">Disabled</span>";
                } else if (item.target && item.target.disabled) {
                    statusBadge = "<span class=\"badge badge-danger\">Weekly depleted</span>";
                } else if (item.reason && item.reason.indexOf("429") >= 0) {
                    statusBadge = "<span class=\"badge badge-warning\">429 cooldown</span>";
                } else if (isBoosted) {
                    statusBadge = "<span class=\"badge badge-boost\">Boosted</span>";
                }
                const formattedReason = formatReason(item.reason, isBoosted, item.target && item.target.disabled);
                const authIdx = item.auth_index || "";
                const credDisplayName = item.email || "Credential";

                card.innerHTML =
                    "<div class=\"cred-head\">" +
                        "<div class=\"cred-name\">" +
                            "<div class=\"cred-title\">" +
                                "<span class=\"cred-email\" title=\"" + escapeHTML(credDisplayName) + "\">" + escapeHTML(credDisplayName) + "</span>" +
                                "<span class=\"cred-meta\"> · " + escapeHTML(item.plan_type || "Antigravity") + "</span>" +
                            "</div>" +
                        "</div>" +
                        "<div class=\"cred-badge-row\">" +
                            statusBadge +
                            "<span class=\"chip\">Urgency <strong>" + urgency + "</strong></span>" +
                            "<span class=\"chip\">Burn rate <strong>" + burnRate + "</strong></span>" +
                        "</div>" +
                    "</div>" +
                    "<div class=\"cred-meters\">" +
                        renderQuotaMeter("5h window", r5hPercent, item.short_window_reset_at) +
                        renderQuotaMeter("7d window", r7dPercent, item.long_window_reset_at) +
                    "</div>" +
                    "<div class=\"cred-priority\">" +
                        "<div class=\"priority-score-box\">" +
                            "<span class=\"priority-label\">Actual</span>" +
                            "<span class=\"priority-actual\">" + escapeHTML(String(actualP)) + "</span>" +
                            "<span class=\"priority-arrow\">&rarr;</span>" +
                            "<span class=\"priority-label\">" + (isPredictedView ? "Predicted" : "Target") + "</span>" +
                            "<span class=\"priority-score" + (isPredictedView ? " priority-predicted" : "") + "\">" + escapeHTML(String(targetP)) + "</span>" +
                            tagBadge +
                        "</div>" +
                        "<div class=\"cred-reason\">" + escapeHTML(formattedReason) + "</div>" +
                    "</div>" +
                    "<div class=\"cred-actions\">" +
                        "<button type=\"button\" class=\"button secondary mini\" data-credential-action=\"probe\" data-auth-index=\"" + escapeHTML(authIdx) + "\">" + ICONS.probe + "<span>Probe</span></button>" +
                        "<button type=\"button\" class=\"button secondary mini\" data-credential-action=\"samples\" data-auth-index=\"" + escapeHTML(authIdx) + "\" data-credential-name=\"" + escapeHTML(credDisplayName) + "\">" + ICONS.chart + "<span>Samples</span></button>" +
                    "</div>";

                container.appendChild(card);
            });
        }

`
