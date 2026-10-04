package management

const templateStyleShell = `        /* Session gate */
        .auth-shell {
            min-height: 100vh;
            min-height: 100dvh;
            display: grid;
            place-items: center;
            padding: 2rem 1.25rem;
        }
        .auth-card {
            width: min(100%, 27rem);
            padding: 2.2rem;
            border: 1px solid var(--border);
            border-radius: 1.4rem;
            background: var(--surface);
            backdrop-filter: blur(26px) saturate(150%);
            box-shadow: var(--shadow-card);
        }
        .brand-mark {
            width: 3.15rem;
            height: 3.15rem;
            display: grid;
            place-items: center;
            border: 1px solid color-mix(in srgb, var(--border-strong) 78%, transparent);
            border-radius: 0.8rem;
            color: var(--text);
            background: var(--surface-solid);
            box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.04), var(--shadow-card);
        }
        .brand-mark svg { width: 1.5rem; height: 1.5rem; }
        .brand-mark.error { color: var(--danger); background: var(--danger-soft); border-color: color-mix(in srgb, var(--danger) 35%, var(--border)); }
        .brand-mark.small { width: 2.25rem; height: 2.25rem; border-radius: 0.65rem; box-shadow: none; }
        .brand-mark.small svg { width: 1.08rem; height: 1.08rem; }
        .eyebrow, .section-label {
            margin: 1.65rem 0 0.36rem;
            color: var(--text-tertiary);
            font-size: 0.68rem;
            line-height: 1.2;
            font-weight: 740;
            letter-spacing: 0.1em;
        }
        .auth-card h1 {
            margin: 0;
            font-size: clamp(1.8rem, 5vw, 2.12rem);
            line-height: 1.08;
            letter-spacing: -0.04em;
        }
        .auth-copy { margin: 0.9rem 0 1.7rem; color: var(--text-secondary); font-size: 0.91rem; line-height: 1.68; }
        .security-note {
            display: flex;
            gap: 0.62rem;
            margin-top: 1.4rem;
            padding-top: 1.2rem;
            border-top: 1px solid var(--border);
            color: var(--text-tertiary);
            font-size: 0.71rem;
            line-height: 1.58;
        }
        .security-note .status-dot { margin-top: 0.34rem; flex: 0 0 auto; }

        /* App chrome */
        .app { min-height: 100vh; min-height: 100dvh; }
        .topbar {
            position: sticky;
            top: 0;
            z-index: 30;
            height: var(--topbar-height);
            display: flex;
            align-items: center;
            gap: 1.25rem;
            padding: 0.8rem max(1.25rem, env(safe-area-inset-right)) 0.8rem max(1.25rem, env(safe-area-inset-left));
            background: color-mix(in srgb, var(--surface) 88%, transparent);
            backdrop-filter: blur(16px);
            border-bottom: 1px solid color-mix(in srgb, var(--border) 78%, transparent);
        }
        /* CPA renders its own floating controls in the upper-right corner. Keep that host-owned area empty. */
        :root.is-embedded .topbar { padding-right: max(12.5rem, env(safe-area-inset-right)); }
        .brand { display: flex; align-items: center; gap: 0.76rem; min-width: max-content; }
        .brand-copy strong, .brand-copy small { display: block; }
        .brand-copy strong { font-size: 0.88rem; line-height: 1.22; letter-spacing: -0.012em; }
        .brand-copy small { margin-top: 0.1rem; color: var(--text-tertiary); font-size: 0.66rem; }
        .version-badge { font-family: var(--mono); letter-spacing: -0.01em; }
        .topbar-context {
            margin: 0;
            padding-left: 1.25rem;
            border-left: 1px solid var(--border);
            color: var(--text-tertiary);
            font-size: 0.72rem;
            white-space: nowrap;
        }

        .workspace {
            min-height: calc(100vh - var(--topbar-height));
            min-height: calc(100dvh - var(--topbar-height));
            display: grid;
            grid-template-columns: var(--sidebar-width) minmax(0, 1fr);
        }
        .sidebar {
            position: sticky;
            top: var(--topbar-height);
            height: calc(100vh - var(--topbar-height));
            height: calc(100dvh - var(--topbar-height));
            display: flex;
            flex-direction: column;
            padding: 1.25rem 0.95rem 1rem;
            border-right: 1px solid var(--border);
            background: var(--surface-solid);
            overflow: hidden;
        }
        .sidebar-head { display: flex; align-items: center; justify-content: space-between; padding: 0 0.32rem 0.9rem; }
        .sidebar-head .section-label { margin: 0 0 0.2rem; }
        .sidebar-count { margin: 0; color: var(--text-secondary); font-size: 0.75rem; }
        .page-nav { min-height: 0; flex: 1; display: grid; align-content: start; gap: 0.18rem; overflow: auto; scrollbar-width: thin; }
        .nav-item {
            position: relative;
            width: 100%;
            min-height: 3rem;
            display: flex;
            align-items: center;
            gap: 0.68rem;
            padding: 0.55rem 0.62rem;
            border: 1px solid transparent;
            border-radius: 0.72rem;
            color: var(--text-secondary);
            background: transparent;
            text-align: left;
            cursor: pointer;
            -webkit-tap-highlight-color: transparent;
            transition: color 150ms ease, background-color 150ms ease, border-color 150ms ease, box-shadow 150ms ease;
        }
        .nav-item[aria-current="page"] { color: var(--text); border-color: var(--border); background: var(--surface-raised); box-shadow: 0 2px 7px rgba(0, 0, 0, 0.045); }
        .nav-item[aria-current="page"]::before { content: ""; position: absolute; left: -0.96rem; width: 0.2rem; height: 1.35rem; border-radius: 0 999px 999px 0; background: var(--accent); }
        .nav-icon { width: 1.8rem; height: 1.8rem; flex: 0 0 auto; display: grid; place-items: center; border-radius: 0.52rem; color: var(--text-secondary); background: var(--surface-subtle); }
        .nav-item[aria-current="page"] .nav-icon { color: var(--accent); background: var(--accent-soft); }
        .nav-icon svg { width: 0.88rem; height: 0.88rem; }
        .nav-copy { min-width: 0; flex: 1; }
        .nav-copy strong, .nav-copy span { display: block; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
        .nav-copy strong { color: inherit; font-size: 0.75rem; font-weight: 640; }
        .nav-copy span { margin-top: 0.1rem; color: var(--text-tertiary); font-size: 0.64rem; line-height: 1.25; }
        .sidebar-note { margin-top: 0.8rem; padding: 0.72rem; border: 1px solid var(--border); border-radius: 0.72rem; color: var(--text-secondary); background: var(--surface); font-size: 0.65rem; line-height: 1.5; }

        /* The command dock sits below the host toolbar instead of competing for the upper-right corner. */
        .main-stage { min-width: 0; }
        .command-wrap {
            position: sticky;
            top: var(--topbar-height);
            z-index: 20;
            padding: 0.8rem clamp(1.2rem, 4vw, 3rem) 0.25rem;
            pointer-events: none;
            background: linear-gradient(to bottom, var(--bg) 45%, transparent);
        }
        .commandbar {
            width: 100%;
            max-width: var(--content-width);
            min-height: 3.35rem;
            display: flex;
            align-items: center;
            justify-content: space-between;
            gap: 1rem;
            margin: 0 auto;
            padding: 0.48rem 0.56rem 0.48rem 0.78rem;
            border: 1px solid var(--border);
            border-radius: 1rem;
            background: color-mix(in srgb, var(--surface-solid) 78%, transparent);
            backdrop-filter: blur(16px);
            box-shadow: 0 18px 44px rgba(0, 0, 0, 0.12), inset 0 1px 0 rgba(255, 255, 255, 0.04);
            pointer-events: auto;
        }
        .command-state { min-width: 0; display: flex; align-items: center; gap: 0.6rem; }
        .command-label { color: var(--text-tertiary); font-size: 0.66rem; font-weight: 680; letter-spacing: 0.04em; white-space: nowrap; }
        .command-actions { min-width: 0; display: flex; align-items: center; gap: 0.38rem; }
        .separator { width: 1px; height: 1.35rem; margin: 0 0.15rem; flex: 0 0 auto; background: var(--border); }

        .editor { width: 100%; max-width: var(--content-width); margin: 0 auto; padding: 1.6rem clamp(1.2rem, 4vw, 3rem) 5rem; }
        .editor-head { display: flex; align-items: flex-start; justify-content: space-between; gap: 1rem; margin-bottom: 1.65rem; }
        .editor-title-wrap { min-width: 0; }
        .editor-kicker { margin: 0 0 0.42rem; color: var(--accent); font-size: 0.65rem; font-weight: 750; letter-spacing: 0.095em; text-transform: uppercase; }
        .editor h1 { margin: 0; font-size: clamp(1.85rem, 4vw, 2.5rem); line-height: 1.06; letter-spacing: -0.048em; }
        .editor-subtitle { max-width: 47rem; margin: 0.68rem 0 0; color: var(--text-secondary); font-size: 0.82rem; line-height: 1.62; }

`
