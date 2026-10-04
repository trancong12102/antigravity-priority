package management

const templateStyleHistory = `        /* Execution history */
        .history-list { display: flex; flex-direction: column; }
        .history-list > .empty-state { margin: 1.15rem; }
        .history-item { display: grid; gap: 0.35rem; padding: 0.85rem 1.15rem; border-bottom: 1px solid var(--border); }
        .history-item:last-child { border-bottom: 0; }
        .history-head { display: flex; align-items: center; justify-content: space-between; flex-wrap: wrap; gap: 0.6rem; }
        .history-when { display: flex; align-items: center; gap: 0.55rem; min-width: 0; }
        .history-when strong { font-size: 0.76rem; font-weight: 650; }
        .history-stats { display: flex; flex-wrap: wrap; align-items: center; gap: 0.35rem; }
        .history-message { color: var(--text-tertiary); font-family: var(--mono); font-size: 0.66rem; line-height: 1.45; overflow-wrap: anywhere; }

`
