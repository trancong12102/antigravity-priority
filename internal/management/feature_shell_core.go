package management

const templateScriptShellCore = `        let activeTab = "overview";
        let countdownInterval = null;
        let dashboardRefreshInterval = null;
        let silentDashboardRefreshInFlight = false;
        let probeCooldownTimer = null;
        let scheduleConfig = null;
        let dynamicConfig = null;
        let isAuthBlocked = false;
        let sessionActive = false;
        let sessionConnecting = false;

        // Mirrors CPAMC secureStorage v1. This is reversible obfuscation for session
        // compatibility, not cryptography and not a new security boundary.
        function decodeStoredValue(raw) {
            if (!raw || raw.indexOf(OBFUSCATION_PREFIX) !== 0) return raw;
            const encoded = atob(raw.slice(OBFUSCATION_PREFIX.length));
            const encrypted = Uint8Array.from(encoded, function(character) { return character.charCodeAt(0); });
            const key = new TextEncoder().encode(OBFUSCATION_SALT + "|" + window.location.host + "|" + navigator.userAgent);
            const decoded = new Uint8Array(encrypted.length);
            for (let index = 0; index < encrypted.length; index += 1) decoded[index] = encrypted[index] ^ key[index % key.length];
            return new TextDecoder().decode(decoded);
        }

        function readStoredValue(name) {
            try {
                const raw = localStorage.getItem(name);
                if (raw === null) return null;
                const decoded = decodeStoredValue(raw);
                try { return JSON.parse(decoded); } catch (_) { return decoded; }
            } catch (_) {
                return null;
            }
        }

        // The Management Key is read from CPAMC's saved session on every request and
        // is never copied into this page's own storage.
        function readCPAMCManagementKey() {
            try {
                if (localStorage.getItem("isLoggedIn") !== "true") return "";
                const persisted = readStoredValue(CPAMC_AUTH_KEY);
                const current = persisted && typeof persisted === "object" && persisted.state ? persisted.state.managementKey : "";
                if (typeof current === "string" && current.trim()) return current.trim();
                const legacy = readStoredValue("managementKey");
                return typeof legacy === "string" ? legacy.trim() : "";
            } catch (_) {
                return "";
            }
        }

        function setSessionState(kind, title, message) {
            document.getElementById("authTitle").textContent = title;
            document.getElementById("authMessage").textContent = message;
            const icon = document.getElementById("sessionIcon");
            icon.classList.toggle("error", kind === "error");
            icon.innerHTML = kind === "loading" ? ICONS.spinner : kind === "error" ? ICONS.warning : ICONS.lock;
            document.getElementById("retrySessionButton").hidden = kind !== "error";
        }

        function resolveCPAMCTheme() {
            try {
                if (window.self !== window.top) {
                    const parentTheme = window.parent.document.documentElement.getAttribute("data-theme");
                    return parentTheme === "dark" || parentTheme === "white" ? parentTheme : "light";
                }
            } catch (_) {}
            try {
                const persisted = JSON.parse(localStorage.getItem(CPAMC_THEME_KEY) || "null");
                const theme = persisted && persisted.state ? persisted.state.theme : "";
                if (theme === "dark" || theme === "white" || theme === "light") return theme;
                if (theme === "auto") return window.matchMedia("(prefers-color-scheme: dark)").matches ? "dark" : "white";
            } catch (_) {}
            return window.matchMedia("(prefers-color-scheme: dark)").matches ? "dark" : "white";
        }

        function applyTheme(theme) {
            if (theme === "light") document.documentElement.removeAttribute("data-theme");
            else document.documentElement.setAttribute("data-theme", theme);
        }

        function syncTheme() {
            applyTheme(resolveCPAMCTheme());
        }

        // Console helper for the devserver and standalone debugging:
        // window.setTheme("light" | "white" | "dark")
        window.__setTheme = window.setTheme = function(theme) {
            const target = theme === "dark" || theme === "white" ? theme : "light";
            applyTheme(target);
            return "Theme switched to: " + target;
        };

        function initializeChrome() {
            document.documentElement.classList.toggle("is-embedded", window.self !== window.top);
            syncTheme();
            try {
                if (window.self !== window.top) {
                    new MutationObserver(syncTheme).observe(window.parent.document.documentElement, { attributes: true, attributeFilter: ["data-theme"] });
                }
            } catch (_) {}
        }

        function isAutoApplyEnabled() {
            if (dynamicConfig && dynamicConfig.auto_apply !== undefined) return Boolean(dynamicConfig.auto_apply);
            if (latestDiagnostics && latestDiagnostics.management_api && latestDiagnostics.management_api.auto_apply !== undefined) {
                return Boolean(latestDiagnostics.management_api.auto_apply);
            }
            return false;
        }

        function isCurrentTimeInScheduleWindow(startStr, endStr) {
            if (!startStr || !endStr) return true;
            var startParts = startStr.trim().split(":");
            var endParts = endStr.trim().split(":");
            if (startParts.length !== 2 || endParts.length !== 2) return true;
            var startMin = parseInt(startParts[0], 10) * 60 + parseInt(startParts[1], 10);
            var endMin = parseInt(endParts[0], 10) * 60 + parseInt(endParts[1], 10);
            if (isNaN(startMin) || isNaN(endMin) || startMin === endMin) return true;
            if (startMin === 0 && endMin === 24 * 60) return true;

            var now = new Date();
            var nowMin = now.getHours() * 60 + now.getMinutes();

            if (startMin < endMin) {
                return nowMin >= startMin && nowMin < endMin;
            }
            // Cross midnight (e.g. 22:00 to 06:00)
            return nowMin >= startMin || nowMin < endMin;
        }

        // Go encodes an unset time.Time as year 1; treat it like a missing value.
        function formatDateTime(value) {
            if (!value) return "-";
            var date = new Date(value);
            if (isNaN(date.getTime()) || date.getUTCFullYear() <= 1) return "-";
            return date.toLocaleString("en-US");
        }

`
