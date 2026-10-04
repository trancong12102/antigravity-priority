package management

const templateStyleOverlays = `        /* Modals */
        .modal-backdrop {
            position: fixed;
            inset: 0;
            z-index: 90;
            display: grid;
            place-items: center;
            padding: 1.25rem;
            background: var(--overlay);
            backdrop-filter: blur(3px);
            animation: backdropIn 180ms ease forwards;
        }
        .modal {
            width: min(720px, 100%);
            max-height: 85vh;
            display: flex;
            flex-direction: column;
            overflow: hidden;
            border: 1px solid var(--border);
            border-radius: var(--radius-lg);
            background: var(--surface-solid);
            box-shadow: var(--shadow-float);
            animation: modalIn 240ms var(--ease-out) forwards;
        }
        .modal.narrow { width: min(28rem, 100%); }
        .modal.medium { width: min(40rem, 100%); }
        .modal.modal-history-comparison { width: fit-content; max-width: 100%; }
        .modal-head { display: flex; align-items: flex-start; justify-content: space-between; gap: 1rem; padding: 1.05rem 1.15rem 0.9rem; border-bottom: 1px solid var(--border); background: color-mix(in srgb, var(--surface-subtle) 42%, transparent); }
        .modal-head h2 { margin: 0; font-size: 0.92rem; letter-spacing: -0.012em; }
        .modal-head p { margin: 0.3rem 0 0; color: var(--text-tertiary); font-size: 0.68rem; line-height: 1.45; }
        .modal-head .icon-button { margin: -0.35rem -0.4rem 0 0; }
        .modal-body {
            max-height: 60vh;
            display: flex;
            flex-direction: column;
            gap: 0.5rem;
            padding: 1rem 1.15rem;
            overflow-x: hidden;
            overflow-y: auto;
        }
        .modal-foot { display: flex; justify-content: flex-end; gap: 0.5rem; padding: 0.8rem 1.15rem; border-top: 1px solid var(--border); background: color-mix(in srgb, var(--surface-subtle) 30%, transparent); }
        .confirm-body { padding: 1.15rem; }
        .confirm-title { margin: 0; font-size: 0.95rem; letter-spacing: -0.012em; }
        .confirm-message { margin: 0.55rem 0 0; color: var(--text-secondary); font-size: 0.78rem; line-height: 1.6; }

        .diff-card { display: flex; align-items: center; justify-content: space-between; flex-wrap: wrap; gap: 0.7rem; padding: 0.75rem 0.85rem; border: 1px solid var(--border); border-radius: var(--radius-sm); background: var(--surface-raised); }
        .diff-name { display: flex; align-items: center; gap: 0.4rem; font-size: 0.76rem; font-weight: 650; }
        .diff-reason { margin-top: 0.15rem; color: var(--text-tertiary); font-size: 0.66rem; }
        .diff-value-box { display: flex; align-items: center; gap: 0.45rem; color: var(--text-tertiary); font-family: var(--mono); font-size: 0.74rem; font-weight: 600; }
        .diff-from, .diff-to { padding: 0.12rem 0.45rem; border-radius: 0.35rem; }
        .diff-from { color: var(--danger); background: var(--danger-soft); }
        .diff-to { color: var(--success); background: var(--success-soft); }

        .history-comparison-table { width: max-content; max-width: 100%; min-width: 0; overflow: hidden; border: 1px solid var(--border); border-radius: var(--radius-sm); background: var(--surface-raised); font-size: 0.68rem; }
        .history-comparison-row { display: grid; grid-template-columns: minmax(150px, max-content) 200px 200px minmax(140px, max-content); align-items: center; gap: 0.65rem; min-height: 3.2rem; padding: 0.5rem 0.7rem; border-top: 1px solid var(--border); }
        .history-comparison-quota-only .history-comparison-row { grid-template-columns: minmax(150px, max-content) 200px 200px; }
        .history-comparison-table, .history-comparison-row, .history-comparison-row > div { min-width: 0; }
        .history-comparison-row:first-child { border-top: 0; }
        .history-comparison-header { min-height: 2.2rem; color: var(--text-tertiary); background: var(--surface-subtle); font-size: 0.62rem; font-weight: 700; letter-spacing: 0.04em; text-transform: uppercase; }
        .history-account-cell { overflow: hidden; color: var(--text); font-weight: 650; text-overflow: ellipsis; white-space: nowrap; }
        .history-quota-values { display: flex; flex-direction: column; gap: 0.25rem; }
        .history-quota-values > div { display: flex; align-items: center; gap: 0.4rem; }
        .history-quota-values b { width: 1.3rem; color: var(--text-tertiary); font-size: 0.6rem; }
        .history-quota-meter { display: flex; align-items: center; gap: 0.4rem; width: 168px; }
        .history-quota-meter > span { min-width: 2rem; font-family: var(--mono); font-weight: 600; text-align: right; }
        .history-meter-track { flex: 1; min-width: 120px; height: 0.42rem; overflow: hidden; border: 1px solid var(--meter-track-border); border-radius: 999px; background: var(--meter-track-bg); }
        .history-meter-fill { height: 100%; border-radius: inherit; background: var(--success); }
        .history-priority-value { display: flex; align-items: center; gap: 0.4rem; font-family: var(--mono); font-weight: 650; }
        .history-priority-value b { color: var(--text-tertiary); font-weight: 400; }
        .history-priority-reason { margin-top: 0.25rem; color: var(--text-tertiary); line-height: 1.25; }

        /* The comparison needs ~760px as a table; below that each account
           becomes a card with labelled cells instead of being clipped. */
        @media (max-width: 820px) {
            .modal.modal-history-comparison { width: 100%; }
            .history-comparison-table { width: 100%; font-size: 0.7rem; }
            .history-comparison-header { display: none; }
            .history-comparison-row,
            .history-comparison-quota-only .history-comparison-row { grid-template-columns: repeat(2, minmax(0, 1fr)); align-items: start; gap: 0.55rem 0.9rem; padding: 0.75rem 0.85rem; }
            .history-comparison-header + .history-comparison-row { border-top: 0; }
            .history-account-cell { grid-column: 1 / -1; overflow-wrap: anywhere; white-space: normal; font-size: 0.76rem; }
            .history-priority-cell { grid-column: 1 / -1; }
            .history-cell::before { content: attr(data-label); display: block; margin-bottom: 0.3rem; color: var(--text-tertiary); font-size: 0.6rem; font-weight: 700; letter-spacing: 0.04em; text-transform: uppercase; }
            .history-quota-values > div, .history-quota-meter { width: 100%; min-width: 0; }
            .history-meter-track { min-width: 0; }
        }
        .history-empty-cell { display: block; min-height: 1.1rem; }

        .sample-table { width: 100%; border-collapse: separate; border-spacing: 0; overflow: hidden; border: 1px solid var(--border); border-radius: var(--radius-sm); background: var(--surface-raised); font-size: 0.7rem; text-align: left; }
        .sample-table th { padding: 0.55rem 0.75rem; border-bottom: 1px solid var(--border); color: var(--text-tertiary); background: var(--surface-subtle); font-size: 0.62rem; font-weight: 700; letter-spacing: 0.04em; text-transform: uppercase; white-space: nowrap; }
        .sample-table td { padding: 0.45rem 0.75rem; border-bottom: 1px solid var(--border); color: var(--text); white-space: nowrap; }
        .sample-table tbody.sample-group:last-child td { border-bottom: 0; }
        .sample-table td.sample-index { width: 2.6rem; color: var(--text-tertiary); font-family: var(--mono); text-align: center; }
        .sample-table td.sample-time { color: var(--text-secondary); font-family: var(--mono); }

        /* Feedback */
        .toast-region { position: fixed; z-index: 100; right: max(1rem, env(safe-area-inset-right)); bottom: max(1rem, env(safe-area-inset-bottom)); display: grid; gap: 0.5rem; pointer-events: none; }
        .toast { min-width: 16rem; max-width: min(24rem, calc(100vw - 2rem)); display: flex; align-items: center; gap: 0.65rem; padding: 0.78rem 0.84rem; border: 1px solid var(--border); border-radius: 0.78rem; background: var(--surface-raised); box-shadow: var(--shadow-float); color: var(--text); font-size: 0.73rem; opacity: 1; transform: translateY(0); transition: opacity 170ms ease, transform 190ms var(--ease-out); pointer-events: auto; }
        .toast.enter { opacity: 0; transform: translateY(0.5rem); }
        .toast.exit { opacity: 0; transform: translateY(0.35rem); }
        .toast svg { color: var(--success); }
        .toast.error svg { color: var(--danger); }
        .toast.info svg { color: var(--accent); }
        .toast span { flex: 1; }
        .spinner { animation: spin 700ms linear infinite; }
        @keyframes spin { to { transform: rotate(360deg); } }
        @keyframes backdropIn { from { opacity: 0; } to { opacity: 1; } }
        @keyframes modalIn { from { opacity: 0; transform: translateY(0.6rem) scale(0.98); } to { opacity: 1; transform: translateY(0) scale(1); } }

`
