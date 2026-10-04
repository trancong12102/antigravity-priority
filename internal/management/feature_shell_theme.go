package management

// templateScriptThemeBootstrap runs in <head>, before the stylesheet and body
// are parsed, so the first paint already uses CPAMC's theme instead of
// flashing the light palette for dark users. It owns theme resolution for the
// whole page; the shell scripts only re-sync on later changes.
const templateScriptThemeBootstrap = `
    <script>
        const CPAMC_THEME_KEY = "cli-proxy-theme";

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

        document.documentElement.classList.toggle("is-embedded", window.self !== window.top);
        syncTheme();
    </script>`
