package management

const templateScriptLiveUI = `        function updateAllCountdowns() {
            document.querySelectorAll(".meter-countdown[data-reset-time]").forEach(function(el) {
                var resetTime = el.getAttribute("data-reset-time");
                if (resetTime) {
                    el.textContent = formatCountdown(resetTime);
                }
            });

            document.querySelectorAll(".meter-countdown[data-cooldown-countdown]").forEach(function(el) {
                var cdTime = el.getAttribute("data-cooldown-countdown");
                if (cdTime) {
                    el.textContent = formatCountdown(cdTime);
                }
            });

            var isPaused = scheduleConfig ? Boolean(scheduleConfig.paused) : (latestDiagnostics && latestDiagnostics.scheduler && Boolean(latestDiagnostics.scheduler.paused));
            var isSleeping = false;
            if (scheduleConfig && scheduleConfig.window_enabled && scheduleConfig.window_start && scheduleConfig.window_end) {
                isSleeping = !isCurrentTimeInScheduleWindow(scheduleConfig.window_start, scheduleConfig.window_end);
            }

            if (isAutoApplyEnabled() && !isPaused && !isSleeping) {
                document.querySelectorAll("[data-scheduler-countdown]").forEach(function(el) {
                    var t = el.getAttribute("data-scheduler-countdown");
                    if (t) el.textContent = formatCountdown(t);
                });
            }
            // Periodically refresh scheduler status indicator to catch window transitions
            if (scheduleConfig) {
                renderScheduleStatus();
            }
        }

        async function refreshOverviewSilently() {
            if (isAuthBlocked || activeTab !== "overview" || document.hidden || silentDashboardRefreshInFlight) return;
            silentDashboardRefreshInFlight = true;
            try {
                await refreshDashboard({ silent: true });
            } finally {
                silentDashboardRefreshInFlight = false;
            }
        }

        function showToast(msg, type) {
            var root = document.getElementById("toastRoot");
            if (!root) return;
            var kind = type === "error" ? "error" : type === "info" ? "info" : "success";
            var toast = document.createElement("div");
            toast.className = "toast " + kind + " enter";
            toast.innerHTML = (kind === "error" ? ICONS.warning : kind === "info" ? ICONS.info : ICONS.check) + "<span>" + escapeHTML(msg) + "</span>";
            root.appendChild(toast);
            requestAnimationFrame(function() { toast.classList.remove("enter"); });
            setTimeout(function() {
                toast.classList.add("exit");
                setTimeout(function() { toast.remove(); }, 190);
            }, 3600);
        }

        function escapeHTML(str) {
            return String(str || "")
                .replace(/&/g, "&amp;")
                .replace(/</g, "&lt;")
                .replace(/>/g, "&gt;")
                .replace(/"/g, "&quot;")
                .replace(/'/g, "&#039;");
        }

        function showThemedConfirm(opts) {
            return new Promise(function(resolve) {
                var modal = document.getElementById("confirmModal");
                var titleEl = document.getElementById("confirmTitle");
                var msgEl = document.getElementById("confirmMessage");
                var okBtn = document.getElementById("confirmOkBtn");
                var cancelBtn = document.getElementById("confirmCancelBtn");

                titleEl.textContent = opts.title || "";
                msgEl.textContent = opts.message || "";
                okBtn.textContent = opts.confirmText || "OK";
                cancelBtn.textContent = opts.cancelText || "Cancel";

                okBtn.className = "button compact " + (opts.isDanger ? "danger-solid" : "primary");

                function cleanup(result) {
                    modal.hidden = true;
                    okBtn.removeEventListener("click", onOk);
                    cancelBtn.removeEventListener("click", onCancel);
                    document.removeEventListener("keydown", onKey);
                    resolve(result);
                }
                function onOk() { cleanup(true); }
                function onCancel() { cleanup(false); }
                function onKey(e) { if (e.key === "Escape") cleanup(false); }

                okBtn.addEventListener("click", onOk);
                cancelBtn.addEventListener("click", onCancel);
                document.addEventListener("keydown", onKey);
                modal.hidden = false;
            });
        }

`
