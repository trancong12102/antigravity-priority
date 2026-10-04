package management

const templateStyleOverviewCards = `        /* Credential cards */
        .credentials-grid {
            display: grid;
            grid-template-columns: repeat(auto-fill, minmax(18.5rem, 1fr));
            align-content: start;
            gap: 0.72rem;
        }
        .cred-card { display: flex; flex-direction: column; overflow: hidden; border: 1px solid var(--border); border-radius: var(--radius-md); background: var(--surface-solid); box-shadow: 0 2px 9px rgba(0, 0, 0, 0.025); }
        .cred-head { display: grid; gap: 0.55rem; padding: 0.85rem 0.95rem 0.75rem; border-bottom: 1px solid var(--border); background: color-mix(in srgb, var(--surface-subtle) 42%, transparent); }
        .cred-name { min-width: 0; display: flex; align-items: center; justify-content: space-between; gap: 0.6rem; }
        .cred-title { min-width: 0; display: flex; align-items: baseline; gap: 0.25rem; font-size: 0.78rem; }
        .cred-email { min-width: 0; overflow: hidden; font-weight: 650; text-overflow: ellipsis; white-space: nowrap; }
        .cred-meta { flex: 0 0 auto; color: var(--text-tertiary); font-size: 0.64rem; white-space: nowrap; }
        .cred-badge-row { display: flex; flex-wrap: wrap; align-items: center; gap: 0.35rem; }
        .cred-meters { display: grid; gap: 0.7rem; padding: 0.85rem 0.95rem; }
        .cred-priority { display: grid; gap: 0.3rem; padding: 0 0.95rem 0.85rem; }
        .priority-score-box { display: flex; flex-wrap: wrap; align-items: center; gap: 0.4rem; font-size: 0.68rem; }
        .priority-label { color: var(--text-tertiary); font-weight: 600; }
        .priority-actual { color: var(--text-secondary); font-family: var(--mono); font-weight: 600; }
        .priority-arrow { color: var(--text-tertiary); }
        .priority-score { color: var(--text); font-family: var(--mono); font-size: 0.95rem; font-weight: 700; letter-spacing: -0.02em; line-height: 1; }
        .priority-predicted { color: var(--boost); }
        .cred-reason { color: var(--text-tertiary); font-size: 0.66rem; line-height: 1.4; }
        .cred-actions { display: flex; justify-content: flex-end; gap: 0.35rem; margin-top: auto; padding: 0.6rem 0.95rem; border-top: 1px solid var(--border); }

`
