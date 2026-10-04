package management

const templateStyleDiagnostics = `        /* Diagnostics */
        .diag-audit { display: grid; gap: 0.55rem; }
        .diag-audit .command-label { margin: 0; }
        .diag-audit-text { padding: 0.75rem 0.85rem; border: 1px solid var(--border); border-radius: var(--radius-sm); background: var(--surface-raised); color: var(--text); font-family: var(--mono); font-size: 0.7rem; line-height: 1.5; overflow-wrap: anywhere; }
        .diag-pills-row { display: flex; flex-wrap: wrap; align-items: center; gap: 0.35rem; }
        .diag-cooldown-list { display: grid; gap: 0.5rem; }
        .diag-cooldown-item { display: flex; align-items: center; justify-content: space-between; flex-wrap: wrap; gap: 0.75rem; padding: 0.75rem 0.85rem; border: 1px solid color-mix(in srgb, var(--caution) 24%, var(--border)); border-radius: var(--radius-sm); background: var(--caution-soft); }
        .diag-cooldown-who { display: grid; gap: 0.25rem; min-width: 0; }
        .diag-cooldown-name { display: flex; flex-wrap: wrap; align-items: center; gap: 0.4rem; font-size: 0.76rem; font-weight: 650; }
        .diag-cooldown-reason { color: var(--text-tertiary); font-family: var(--mono); font-size: 0.66rem; }
        .diag-cooldown-timer { display: grid; gap: 0.1rem; text-align: right; color: var(--text-tertiary); font-size: 0.64rem; }
        .diag-cooldown-timer .meter-countdown { color: var(--caution); font-size: 0.8rem; font-weight: 650; }

`
