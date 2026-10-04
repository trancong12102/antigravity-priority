package management

const templateScriptControls = `        function onModelGroupChange() {
            userSelectedModelGroup = true;
            renderDashboard();
        }

        // Follows the server's active model group until the user picks one.
        function syncModelGroupSelect(group) {
            if (userSelectedModelGroup || !group) return;
            const select = document.getElementById("modelGroupSelect");
            if (select && select.value !== group) select.value = group;
        }

        function switchTab(tabId) {
            activeTab = tabId;
            document.querySelectorAll(".nav-item[data-tab]").forEach(tab => {
                tab.setAttribute("aria-current", tab.dataset.tab === tabId ? "page" : "false");
            });
            document.getElementById("panelOverview").hidden = tabId !== "overview";
            document.getElementById("panelHistory").hidden = tabId !== "history";
            document.getElementById("panelDiagnostics").hidden = tabId !== "diagnostics";
            document.getElementById("panelConfig").hidden = tabId !== "config";
            document.getElementById("panelHelp").hidden = tabId !== "help";
            window.scrollTo(0, 0);
            if (isAuthBlocked) return;
            if (tabId === "history") fetchDiagnostics();
            if (tabId === "diagnostics") fetchDiagnostics();
            if (tabId === "config") fetchDynamicConfig();
            if (tabId === "overview") refreshDashboard({ withSync: true, silent: true });
        }

        function stopBackgroundRefresh() {
            if (dashboardRefreshInterval) {
                clearInterval(dashboardRefreshInterval);
                dashboardRefreshInterval = null;
            }
            if (countdownInterval) {
                clearInterval(countdownInterval);
                countdownInterval = null;
            }
        }

        // Locks the page behind the session gate. The page never asks for the
        // Management Key; it waits for CPAMC to hold a valid session again.
        function lockSession(title, message) {
            sessionActive = false;
            stopBackgroundRefresh();
            ["diffModal", "samplesModal", "confirmModal"].forEach(function(id) {
                const modal = document.getElementById(id);
                if (modal) modal.hidden = true;
            });
            document.getElementById("app").hidden = true;
            document.getElementById("authGate").hidden = false;
            setSessionState("error", title, message);
        }

        function getAuthHeader() {
            const headers = { "Content-Type": "application/json" };
            const key = readCPAMCManagementKey();
            if (key) {
                headers["Authorization"] = "Bearer " + key;
                headers["X-Management-Key"] = key;
            }
            return headers;
        }

        function sessionError(message) {
            const error = new Error(message);
            error.status = 401;
            return error;
        }

        async function apiFetch(path, options) {
            if (isAuthBlocked) {
                throw sessionError("CPAMC management session required (401 Unauthorized)");
            }
            const resp = await fetch(path, {
                ...(options || {}),
                cache: "no-store",
                headers: { ...getAuthHeader(), ...((options && options.headers) || {}) }
            });
            if (resp.status === 401) {
                isAuthBlocked = true;
                lockSession("Management session unavailable", "The Management Key saved by CPAMC was rejected. Log in to CPAMC again with \"Remember password\" on; this page reconnects automatically.");
                throw sessionError("CPAMC management session required (401 Unauthorized)");
            }
            if (resp.status === 403) {
                const text = await resp.text();
                if ((text.indexOf("banned") >= 0 || text.indexOf("IP") >= 0) && dashboardRefreshInterval) {
                    clearInterval(dashboardRefreshInterval);
                    dashboardRefreshInterval = null;
                }
                throw new Error(text || resp.statusText);
            }
            const text = await resp.text();
            if (!resp.ok) {
                let errMessage = text || resp.statusText;
                try {
                    const parsed = JSON.parse(text);
                    if (parsed.error) errMessage = parsed.error;
                } catch (_) {}
                throw new Error(errMessage);
            }
            return text ? JSON.parse(text) : {};
        }

`
