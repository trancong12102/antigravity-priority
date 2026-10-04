package management

// templateStyleTokens mirrors the key-model-access design tokens so both CPA
// plugin pages read as one product. The light palette is the default; CPAMC's
// "white" and "dark" themes are selected through data-theme on <html>.
const templateStyleTokens = `
        :root {
            color-scheme: light;
            font: 100%/1.5 ui-sans-serif, -apple-system, BlinkMacSystemFont, "SF Pro Text", "Segoe UI", sans-serif;
            font-optical-sizing: auto;
            --bg: #faf9f5;
            --surface: rgba(240, 238, 232, 0.82);
            --surface-solid: #f0eee8;
            --surface-raised: #fffdf9;
            --surface-subtle: #e9e6df;
            --surface-selected: rgba(45, 42, 38, 0.08);
            --text: #2d2a26;
            --text-secondary: #6d6760;
            --text-tertiary: #a29c95;
            --border: #e3e1db;
            --border-strong: #d5d2cb;
            --accent: #8b8680;
            --accent-hover: #7f7a74;
            --accent-soft: rgba(139, 134, 128, 0.15);
            --success: #10b981;
            --success-soft: rgba(16, 185, 129, 0.1);
            --caution: #b7791f;
            --caution-soft: rgba(183, 121, 31, 0.12);
            --danger: #c65746;
            --danger-soft: rgba(198, 87, 70, 0.1);
            --boost: #7c5cbf;
            --boost-soft: rgba(124, 92, 191, 0.11);
            --meter-track-bg: #e9e6df;
            --meter-track-border: #d5d2cb;
            --overlay: rgba(45, 42, 38, 0.32);
            --shadow-card: 0 1px 2px rgba(0, 0, 0, 0.08);
            --shadow-float: 0 12px 26px rgba(0, 0, 0, 0.14);
            --radius-sm: 0.5rem;
            --radius-md: 0.65rem;
            --radius-lg: 1rem;
            --topbar-height: 4.5rem;
            --sidebar-width: 17rem;
            --content-width: 68rem;
            --mono: ui-monospace, "SFMono-Regular", Consolas, monospace;
            --ease-out: cubic-bezier(0.23, 1, 0.32, 1);
        }

        :root[data-theme="white"] {
            --bg: #ffffff;
            --surface: rgba(255, 255, 255, 0.84);
            --surface-solid: #ffffff;
            --surface-raised: #ffffff;
            --surface-subtle: #f6f6f6;
            --surface-selected: rgba(45, 42, 38, 0.07);
            --border: #e5e5e5;
            --border-strong: #d9d9d9;
            --meter-track-bg: #f0f0f0;
            --meter-track-border: #d9d9d9;
        }

        :root[data-theme="dark"] {
            color-scheme: dark;
            --bg: #151412;
            --surface: rgba(29, 27, 24, 0.82);
            --surface-solid: #1d1b18;
            --surface-raised: #2a2723;
            --surface-subtle: #262320;
            --surface-selected: rgba(246, 244, 241, 0.1);
            --text: #f6f4f1;
            --text-secondary: #c9c3bb;
            --text-tertiary: #9c958d;
            --border: #3a3530;
            --border-strong: #4a453f;
            --accent: #8b8680;
            --accent-hover: #9a948e;
            --accent-soft: rgba(139, 134, 128, 0.18);
            --success: #10b981;
            --success-soft: rgba(16, 185, 129, 0.12);
            --caution: #d69e2e;
            --caution-soft: rgba(214, 158, 46, 0.16);
            --danger: #c65746;
            --danger-soft: rgba(198, 87, 70, 0.2);
            --boost: #a78bda;
            --boost-soft: rgba(167, 139, 218, 0.16);
            --meter-track-bg: #2f2b27;
            --meter-track-border: #4a453f;
            --overlay: rgba(0, 0, 0, 0.55);
            --shadow-card: 0 1px 3px rgba(0, 0, 0, 0.3);
            --shadow-float: 0 14px 30px rgba(0, 0, 0, 0.4);
        }

        * { box-sizing: border-box; }
        html, body { min-height: 100%; }
        html { background: var(--bg); }
        body {
            margin: 0;
            overflow-x: hidden;
            background: var(--bg);
            color: var(--text);
            transition: background-color 200ms ease, color 200ms ease;
            -webkit-font-smoothing: antialiased;
            text-rendering: optimizeLegibility;
        }
        button, input, select { font: inherit; }
        button { color: inherit; }
        [hidden] { display: none !important; }

        svg {
            display: block;
            width: 1.1rem;
            height: 1.1rem;
            flex: 0 0 auto;
            fill: none;
            stroke: currentColor;
            stroke-width: 1.75;
            stroke-linecap: round;
            stroke-linejoin: round;
        }

`
