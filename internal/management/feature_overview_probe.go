package management

const templateScriptProbeSchedule = `        async function triggerProbe() {
            var btn = document.getElementById("btnProbe");
            if (!btn || btn.classList.contains("btn-cooldown")) return;

            btn.disabled = true;
            btn.classList.add("btn-cooldown");

            try {
                var groupSelect = document.getElementById("modelGroupSelect");
                var group = groupSelect ? groupSelect.value : "gemini";
                var path = RUN_PATH + "?mode=probe&antigravity_model_group=" + encodeURIComponent(group);
                await apiFetch(path, { method: "POST" });
                showToast("Quota probe complete", "success");
                await refreshDashboard();
            } catch (err) {
                showToast(err.message, "error");
            }

            var remaining = 10;
            var label = btn.querySelector(".label-long");
            var origText = label ? label.textContent : "";
            probeCooldownTimer = setInterval(function() {
                remaining--;
                if (label) label.textContent = "Cooling down (" + remaining + "s)";
                if (remaining <= 0) {
                    clearInterval(probeCooldownTimer);
                    probeCooldownTimer = null;
                    btn.disabled = false;
                    btn.classList.remove("btn-cooldown");
                    if (label) label.textContent = origText;
                }
            }, 1000);
        }

        async function fetchScheduleConfig() {
            if (isAuthBlocked) return;
            try {
                scheduleConfig = await apiFetch(SCHEDULE_CONFIG_PATH);
                renderScheduleStatus();
            } catch (_) {}
        }

        function setScheduleBadge(tone, text, title) {
            var badge = document.getElementById("scheduleStatusBadge");
            if (!badge) return;
            badge.title = title;
            badge.innerHTML = "<span class=\"status-dot " + tone + "\"></span><span id=\"scheduleStatusText\">" + escapeHTML(text) + "</span>";
        }

        function renderScheduleStatus() {
            if (!isAutoApplyEnabled()) {
                setScheduleBadge("neutral", "Disabled", "Auto scheduling is off. Click to open Configuration.");
                return;
            }

            if (!scheduleConfig) return;

            if (scheduleConfig.paused) {
                setScheduleBadge("warning", "Paused", "Click to resume auto scheduling");
                return;
            }

            var windowActive = true;
            if (scheduleConfig.window_enabled && scheduleConfig.window_start && scheduleConfig.window_end) {
                windowActive = isCurrentTimeInScheduleWindow(scheduleConfig.window_start, scheduleConfig.window_end);
            }

            if (scheduleConfig.window_enabled && !windowActive) {
                setScheduleBadge("boost", "Sleeping (outside window)", "The current time is outside the active schedule window");
                return;
            }

            var windowStr = "";
            if (scheduleConfig.window_enabled && scheduleConfig.window_start && scheduleConfig.window_end) {
                windowStr = " (" + scheduleConfig.window_start + "-" + scheduleConfig.window_end + ")";
            }
            setScheduleBadge("", "Running" + windowStr, "Click to pause auto scheduling");
        }

        async function toggleSchedulePause() {
            if (!isAutoApplyEnabled()) {
                switchTab("config");
                showToast("Turn on \"Auto periodic scheduling\" to start the scheduler", "info");
                return;
            }

            if (!scheduleConfig) {
                await fetchScheduleConfig();
                if (!scheduleConfig) return;
            }

            var newConfig = {
                paused: !scheduleConfig.paused,
                window_enabled: scheduleConfig.window_enabled || false,
                window_start: scheduleConfig.window_start || "00:00",
                window_end: scheduleConfig.window_end || "23:59"
            };

            try {
                scheduleConfig = await apiFetch(SCHEDULE_CONFIG_PATH, {
                    method: "POST",
                    body: JSON.stringify(newConfig)
                });
                renderScheduleStatus();
                renderDashboard();
                renderDiagnostics();
                showToast(scheduleConfig.paused ? "Auto schedule paused" : "Auto schedule running", "info");
            } catch (err) {
                showToast(err.message, "error");
            }
        }

`
