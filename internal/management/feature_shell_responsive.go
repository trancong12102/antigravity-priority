package management

const templateStyleResponsive = `        @media (hover: hover) and (pointer: fine) {
            .button.primary:not(:disabled):hover { background: var(--accent-hover); box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.18), 0 6px 16px var(--accent-soft); }
            .button.secondary:not(:disabled):hover, .icon-button:not(:disabled):hover, button.health-badge:hover { border-color: var(--border); background-color: var(--surface-subtle); }
            .button.danger:not(:disabled):hover { background: color-mix(in srgb, var(--danger) 16%, transparent); }
            .nav-item:hover { color: var(--text); background: var(--surface-subtle); }
            .nav-item[aria-current="page"]:hover { background: var(--surface-raised); }
            .sample-table tbody.sample-group:hover td { background: var(--surface-subtle); }
        }

        @media (max-width: 1180px) {
            .stat-grid { grid-template-columns: repeat(2, minmax(0, 1fr)); }
            .stat-grid.three { grid-template-columns: repeat(3, minmax(0, 1fr)); }
        }

        @media (max-width: 1040px) {
            :root { --sidebar-width: 16rem; }
            .command-actions .label-long { display: none; }
            .command-label { display: none; }
        }

        @media (max-width: 760px) {
            :root { --topbar-height: 4.25rem; }
            :root.is-embedded .topbar { padding-right: max(8.5rem, env(safe-area-inset-right)); }
            .topbar-context { display: none; }
            .workspace { display: block; }
            .sidebar { position: sticky; z-index: 25; width: 100%; height: auto; padding: 0.6rem 1rem; border-right: 0; border-bottom: 1px solid var(--border); overflow: visible; }
            .sidebar-head, .sidebar-note { display: none; }
            .page-nav { display: flex; gap: 0.42rem; overflow-x: auto; padding-bottom: 0.15rem; }
            .nav-item { width: auto; min-width: max-content; min-height: 2.5rem; }
            .nav-item[aria-current="page"]::before { display: none; }
            .nav-copy span { display: none; }
            .command-wrap { position: static; padding: 0.75rem 1rem 0; }
            .editor { padding-top: 1.4rem; }
            .stat-grid.three { grid-template-columns: 1fr; }
        }

        @media (max-width: 560px) {
            .auth-card { padding: 1.55rem; }
            .brand-copy small { display: none; }
            .brand-mark.small { width: 2.15rem; height: 2.15rem; }
            .icon-button { width: 2.65rem; height: 2.65rem; }
            .commandbar { align-items: stretch; flex-direction: column; padding: 0.65rem; }
            .command-state { justify-content: space-between; }
            .command-label { display: inline; }
            .command-actions { width: 100%; flex-wrap: wrap; }
            .command-actions .button.primary { margin-left: auto; }
            .command-actions .label-long { display: inline; }
            .command-actions select { flex: 1 1 100%; width: 100%; }
            .health-badge { max-width: 14rem; }
            .health-badge span:last-child { overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
            .editor-head, .setting-row { align-items: stretch; flex-direction: column; }
            .stat-grid { grid-template-columns: repeat(2, minmax(0, 1fr)); }
            .setting-control, .setting-control.wide { width: 100%; flex-wrap: wrap; justify-content: flex-start; text-align: left; }
            .setting-control input[type="number"], .setting-control select { width: 100%; }
            .credentials-grid { grid-template-columns: 1fr; }
            .modal-backdrop { padding: 0.75rem; place-items: end center; }
        }

        @media (prefers-reduced-motion: reduce) {
            *, *::before, *::after { scroll-behavior: auto !important; animation-duration: 1ms !important; animation-iteration-count: 1 !important; transition-duration: 1ms !important; }
        }
        @media (prefers-reduced-transparency: reduce) {
            .auth-card, .topbar, .commandbar { background: var(--surface-solid); backdrop-filter: none; }
        }
        @media (prefers-contrast: more) {
            :root { --border: rgba(0, 0, 0, 0.28); --border-strong: rgba(0, 0, 0, 0.48); }
            :root[data-theme="dark"] { --border: rgba(255, 255, 255, 0.3); --border-strong: rgba(255, 255, 255, 0.5); }
        }
`
