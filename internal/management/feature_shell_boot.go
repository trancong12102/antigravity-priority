package management

const templateScriptBoot = `        // Connects through the CPAMC session; the first authenticated request doubles
        // as the session check before the dashboard is revealed.
        async function startSession() {
            if (sessionConnecting) return;
            if (!readCPAMCManagementKey()) {
                lockSession("No reusable CPAMC session found", "Automatic connection needs CPAMC on the same origin as CPA, logged in with \"Remember password\" on. Log in there and this page connects by itself; it never asks for the Management Key.");
                return;
            }
            sessionConnecting = true;
            isAuthBlocked = false;
            setSessionState("loading", "Connecting to management session", "Verifying the saved CPAMC connection and loading Antigravity credentials…");
            try {
                scheduleConfig = await apiFetch(SCHEDULE_CONFIG_PATH);
            } catch (err) {
                if (!isAuthBlocked) lockSession("Management session unavailable", "Cannot reach CPA: " + err.message);
                return;
            } finally {
                sessionConnecting = false;
            }

            sessionActive = true;
            document.getElementById("authGate").hidden = true;
            document.getElementById("app").hidden = false;
            await fetchDynamicConfig();
            if (isAuthBlocked) return;
            renderScheduleStatus();
            await refreshDashboard({ withSync: true, silent: true });
            if (isAuthBlocked) return;
            if (countdownInterval) clearInterval(countdownInterval);
            countdownInterval = setInterval(updateAllCountdowns, 1000);
            if (dashboardRefreshInterval) clearInterval(dashboardRefreshInterval);
            dashboardRefreshInterval = setInterval(refreshOverviewSilently, 15000);
        }

        // CPAMC writes its session and theme to this origin's localStorage, so a
        // login, logout or theme change in another document arrives here.
        window.addEventListener("storage", function(event) {
            if (event.key === CPAMC_THEME_KEY) syncTheme();
            if (event.key !== null && event.key !== CPAMC_AUTH_KEY && event.key !== "isLoggedIn") return;
            if (!sessionActive) {
                startSession();
                return;
            }
            if (!readCPAMCManagementKey()) {
                lockSession("CPAMC session ended", "Log in to CPAMC again with \"Remember password\" on; this page reconnects automatically.");
            }
        });

        // Credential cards are re-rendered from server data, so their buttons
        // carry identifiers in data attributes instead of inline handlers.
        document.getElementById("credentialsContainer").addEventListener("click", function(event) {
            var button = event.target.closest("button[data-credential-action]");
            if (!button) return;
            var authIndex = button.dataset.authIndex || "";
            if (button.dataset.credentialAction === "probe") probeSingleCredential(authIndex, button);
            else if (button.dataset.credentialAction === "samples") openSamplesModal(authIndex, button.dataset.credentialName || "");
        });

        document.getElementById("historyList").addEventListener("click", function(event) {
            var button = event.target.closest("button[data-history-index]");
            if (button) showHistoryDetails(Number(button.dataset.historyIndex));
        });

        initializeChrome();
        startSession();
    </script>
`
