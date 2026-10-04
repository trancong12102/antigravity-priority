package management

const templateScriptDiagnostics = `        let latestDiagnostics = null;

        function copyDiagnosticsJSON() {
            if (!latestDiagnostics) {
                showToast("No diagnostics data to copy yet", "info");
                return;
            }
            var text = JSON.stringify(latestDiagnostics, null, 2);
            if (navigator.clipboard && navigator.clipboard.writeText) {
                navigator.clipboard.writeText(text).then(function() {
                    showToast("Diagnostics JSON copied to clipboard", "success");
                }).catch(function() {
                    fallbackCopyText(text);
                });
            } else {
                fallbackCopyText(text);
            }
        }

        function fallbackCopyText(text) {
            var textarea = document.createElement("textarea");
            textarea.value = text;
            textarea.style.position = "fixed";
            textarea.style.left = "-9999px";
            textarea.style.top = "0";
            document.body.appendChild(textarea);
            textarea.focus();
            textarea.select();
            try {
                var successful = document.execCommand("copy");
                if (successful) {
                    showToast("Diagnostics JSON copied to clipboard", "success");
                } else {
                    showToast("Copy failed", "error");
                }
            } catch (_) {
                showToast("Copy failed", "error");
            }
            document.body.removeChild(textarea);
        }

        function renderDiagnostics() {
            if (!latestDiagnostics) return;

            var sched = latestDiagnostics.scheduler || {};
            var isPaused = scheduleConfig ? Boolean(scheduleConfig.paused) : Boolean(sched.paused);
            var isSleeping = false;
            var windowStart = scheduleConfig && scheduleConfig.window_start ? scheduleConfig.window_start : sched.window_start;
            var windowEnd = scheduleConfig && scheduleConfig.window_end ? scheduleConfig.window_end : sched.window_end;
            var windowEnabled = scheduleConfig && scheduleConfig.window_enabled !== undefined ? scheduleConfig.window_enabled : sched.window_enabled;

            if (windowEnabled && windowStart && windowEnd) {
                isSleeping = !isCurrentTimeInScheduleWindow(windowStart, windowEnd);
            }

            // --- 1. Scheduler engine ---
            var schedBadge = document.getElementById("diagSchedBadge");
            var schedInterval = document.getElementById("diagSchedInterval");
            var schedCountdown = document.getElementById("diagSchedCountdown");
            if (schedBadge && schedInterval && schedCountdown) {
                var nextRunAt = sched.next_run_at;
                var nextStr = nextRunAt ? formatCountdown(nextRunAt) : (sched.next_wait || "-");

                if (!isAutoApplyEnabled()) {
                    schedBadge.className = "badge badge-subtle";
                    schedBadge.textContent = "Disabled";
                    schedCountdown.textContent = "Auto schedule disabled";
                } else if (isPaused) {
                    schedBadge.className = "badge badge-warning";
                    schedBadge.textContent = "Paused";
                    schedCountdown.textContent = "Auto schedule paused";
                } else if (isSleeping) {
                    schedBadge.className = "badge badge-predicted";
                    schedBadge.textContent = "Sleeping";
                    schedCountdown.textContent = "Outside the schedule window";
                } else {
                    schedBadge.className = "badge badge-success";
                    schedBadge.textContent = "Running";
                    schedCountdown.innerHTML = "Next run in <span class=\"meter-countdown\" data-scheduler-countdown=\"" + (nextRunAt || "") + "\">" + nextStr + "</span>";
                }
                schedInterval.textContent = sched.interval || "-";
            }

            // --- 2. 429 cooldowns ---
            var cooldowns = (latestDiagnostics.active_cooldowns || []).slice();
            var cdBadge = document.getElementById("diagCooldownBadge");
            var cdCount = document.getElementById("diagCooldownCount");
            var cdSub = document.getElementById("diagCooldownSub");
            if (cdBadge && cdCount && cdSub) {
                cdCount.textContent = String(cooldowns.length);
                if (cooldowns.length === 0) {
                    cdBadge.className = "badge badge-success";
                    cdBadge.textContent = "Healthy";
                    cdSub.textContent = "Auto demotion to -1 ready";
                } else {
                    cdBadge.className = "badge badge-danger";
                    cdBadge.textContent = "Cooling down";
                    cdSub.textContent = cooldowns.length + (cooldowns.length === 1 ? " credential" : " credentials") + " demoted to -1 for protection";
                }
            }

            // --- 3. Latest apply health ---
            var latestApply = latestDiagnostics.latest_apply || null;
            var succ = latestApply ? (latestApply.succeeded || 0) : 0;
            var fail = latestApply ? (latestApply.failed || 0) : 0;
            var skip = latestApply ? (latestApply.skipped || 0) : 0;
            var attempted = latestApply ? (latestApply.attempted || 0) : 0;
            var lastRunTime = (latestApply && latestApply.at) || null;

            var applyBadge = document.getElementById("diagApplyBadge");
            var applyStats = document.getElementById("diagApplyStats");
            var applyTime = document.getElementById("diagApplyTime");
            if (applyBadge && applyStats && applyTime) {
                applyTime.textContent = lastRunTime ? formatDateTime(lastRunTime) : "No execution records yet";
                if (!lastRunTime && (attempted === 0 && succ === 0 && fail === 0 && skip === 0)) {
                    applyBadge.className = "badge badge-subtle";
                    applyBadge.textContent = "No runs";
                    applyStats.textContent = "-";
                } else if (fail > 0) {
                    applyBadge.className = "badge badge-danger";
                    applyBadge.textContent = "Has failures";
                    applyStats.textContent = succ + " ok · " + fail + " failed";
                } else if (succ > 0) {
                    applyBadge.className = "badge badge-success";
                    applyBadge.textContent = "All succeeded";
                    applyStats.textContent = succ + " ok · " + skip + " skipped";
                } else {
                    applyBadge.className = "badge badge-subtle";
                    applyBadge.textContent = "No changes";
                    applyStats.textContent = skip + " skipped";
                }
            }

            // --- 4. Scheduling details & window ---
            var detInterval = document.getElementById("diagSchedDetailInterval");
            var detWorker = document.getElementById("diagSchedDetailWorker");
            var detLastRun = document.getElementById("diagSchedDetailLastRun");
            var detNextRun = document.getElementById("diagSchedDetailNextRun");
            var detWinText = document.getElementById("diagSchedDetailWindowText");
            var detWinBadge = document.getElementById("diagSchedDetailWindowBadge");

            if (detInterval) detInterval.textContent = sched.interval || "-";
            if (detWorker) {
                detWorker.innerHTML = sched.worker_active
                    ? "<span class=\"badge badge-success\">Active</span>"
                    : "<span class=\"badge badge-subtle\">Inactive</span>";
            }
            if (detLastRun) {
                var lastRun = formatDateTime(sched.last_auto_apply_at);
                detLastRun.textContent = lastRun === "-" ? "No execution records yet" : lastRun;
            }
            if (detNextRun) detNextRun.textContent = formatDateTime(sched.next_run_at);
            if (detWinText && detWinBadge) {
                if (!windowEnabled || !windowStart || !windowEnd) {
                    detWinText.textContent = "24/7 continuous scheduling";
                    detWinBadge.className = "badge badge-success";
                    detWinBadge.textContent = "24/7";
                } else {
                    detWinText.textContent = windowStart + " - " + windowEnd;
                    if (isSleeping) {
                        detWinBadge.className = "badge badge-predicted";
                        detWinBadge.textContent = "Sleeping, schedule skipped";
                    } else {
                        detWinBadge.className = "badge badge-success";
                        detWinBadge.textContent = "In active window";
                    }
                }
            }

            // --- 5. 429 cooldown list ---
            var cdContent = document.getElementById("diagCooldownContent");
            if (cdContent) {
                if (cooldowns.length === 0) {
                    cdContent.innerHTML = "<div class=\"default-notice\">" + ICONS.check + "<span>No credentials in 429 cooldown; all accounts are operating normally.</span></div>";
                } else {
                    var html = "<div class=\"diag-cooldown-list\">";
                    cooldowns.forEach(function(c) {
                        var cdUntil = c.cooldown_until;
                        var cdCountdownStr = cdUntil ? formatCountdown(cdUntil) : "-";
                        var groupTag = c.model_group ? "<span class=\"badge badge-subtle\">" + escapeHTML(c.model_group) + "</span>" : "";
                        html += "<div class=\"diag-cooldown-item\">" +
                            "<div class=\"diag-cooldown-who\">" +
                                "<div class=\"diag-cooldown-name\">" +
                                    "<span>" + escapeHTML(c.email || "Credential") + "</span>" +
                                    groupTag +
                                    "<span class=\"badge badge-danger\">429 cooldown</span>" +
                                "</div>" +
                                "<div class=\"diag-cooldown-reason\">" + escapeHTML(c.reason || "429 rate limit") + "</div>" +
                            "</div>" +
                            "<div class=\"diag-cooldown-timer\">" +
                                "<span>Recovers in</span>" +
                                "<span class=\"meter-countdown\" data-cooldown-countdown=\"" + (cdUntil || "") + "\">" + cdCountdownStr + "</span>" +
                            "</div>" +
                        "</div>";
                    });
                    html += "</div>";
                    cdContent.innerHTML = html;
                }
            }

            // --- 6. Latest audit ---
            var auditText = document.getElementById("diagAuditText");
            var auditPills = document.getElementById("diagAuditPills");
            if (auditText) {
                auditText.textContent = latestApply ? (latestApply.message || "-") : "No Apply write recorded yet";
            }
            if (auditPills) {
                auditPills.innerHTML =
                    "<span class=\"badge badge-success\">" + succ + " succeeded</span>" +
                    "<span class=\"badge" + (fail ? " badge-danger" : " badge-subtle") + "\">" + fail + " failed</span>" +
                    "<span class=\"badge badge-subtle\">" + skip + " skipped</span>" +
                    "<span class=\"badge\">" + attempted + " attempted</span>";
            }
        }

`
