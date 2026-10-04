package management

const templateStyleControls = `        /* Controls */
        .button, .icon-button, .health-badge {
            -webkit-tap-highlight-color: transparent;
            transition: transform 130ms var(--ease-out), color 150ms ease, background-color 150ms ease, border-color 150ms ease, opacity 150ms ease, box-shadow 150ms ease;
        }
        .button {
            min-height: 2.55rem;
            display: inline-flex;
            align-items: center;
            justify-content: center;
            gap: 0.48rem;
            padding: 0.56rem 0.9rem;
            border: 1px solid transparent;
            border-radius: var(--radius-sm);
            font-size: 0.8rem;
            font-weight: 670;
            white-space: nowrap;
            cursor: pointer;
        }
        .button svg { width: 1rem; height: 1rem; }
        .button.wide { width: 100%; }
        .button.compact { min-height: 2.3rem; padding: 0.43rem 0.76rem; }
        .button.mini { min-height: 1.8rem; gap: 0.32rem; padding: 0.22rem 0.55rem; font-size: 0.68rem; font-weight: 640; }
        .button.mini svg { width: 0.8rem; height: 0.8rem; }
        .button.primary { color: #fff; background: var(--accent); box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.18), 0 4px 12px var(--accent-soft); }
        .button.secondary { border-color: var(--border); background: var(--surface-solid); box-shadow: 0 1px 1px rgba(0, 0, 0, 0.025); }
        .button.danger { border-color: color-mix(in srgb, var(--danger) 24%, var(--border)); color: var(--danger); background: var(--danger-soft); }
        .button.danger-solid { color: #fff; background: var(--danger); box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.18), 0 4px 12px var(--danger-soft); }
        .button:disabled, .icon-button:disabled { cursor: not-allowed; opacity: 0.44; }
        .button:not(:disabled):active, .icon-button:not(:disabled):active { transform: scale(0.97); }
        .icon-button {
            width: 2.3rem;
            height: 2.3rem;
            display: inline-grid;
            flex: 0 0 auto;
            place-items: center;
            padding: 0;
            border: 1px solid transparent;
            border-radius: var(--radius-sm);
            color: var(--text-secondary);
            background: transparent;
            cursor: pointer;
        }
        button:focus-visible, input:focus-visible, select:focus-visible { outline: 2px solid var(--accent); outline-offset: 2px; }

        input[type="text"], input[type="number"], select {
            width: 100%;
            min-height: 2.3rem;
            padding: 0.43rem 0.7rem;
            color: var(--text);
            background: var(--surface-raised);
            border: 1px solid var(--border-strong);
            border-radius: var(--radius-sm);
            font-size: 0.78rem;
            outline: none;
            transition: border-color 150ms ease, box-shadow 150ms ease, background-color 150ms ease, opacity 150ms ease;
        }
        input::placeholder { color: var(--text-tertiary); }
        input:focus, select:focus { border-color: var(--accent); box-shadow: 0 0 0 3px var(--accent-soft); }
        input:disabled { opacity: 0.5; cursor: not-allowed; }
        select {
            width: auto;
            padding-right: 2rem;
            appearance: none;
            cursor: pointer;
            background-image: linear-gradient(45deg, transparent 50%, currentColor 50%), linear-gradient(135deg, currentColor 50%, transparent 50%);
            background-position: calc(100% - 0.92rem) 55%, calc(100% - 0.68rem) 55%;
            background-size: 0.24rem 0.24rem, 0.24rem 0.24rem;
            background-repeat: no-repeat;
        }
        .select-compact { min-height: 2.3rem; font-weight: 600; }

        .switch { position: relative; display: inline-flex; align-items: center; gap: 0.6rem; cursor: pointer; user-select: none; font-size: 0.75rem; font-weight: 600; color: var(--text-secondary); }
        .switch input { position: absolute; width: 1px; height: 1px; opacity: 0; }
        .switch-track { position: relative; width: 2.35rem; height: 1.35rem; flex: 0 0 auto; border-radius: 999px; background: var(--meter-track-border); transition: background-color 180ms var(--ease-out); }
        .switch-track::after { content: ""; position: absolute; top: 0.17rem; left: 0.17rem; width: 1.01rem; height: 1.01rem; border-radius: 50%; background: #fff; box-shadow: 0 1px 2px rgba(0, 0, 0, 0.2); transition: transform 180ms var(--ease-out); }
        .switch input:checked + .switch-track { background: var(--success); }
        .switch input:checked + .switch-track::after { transform: translateX(1rem); }
        .switch input:focus-visible + .switch-track { outline: 2px solid var(--accent); outline-offset: 2px; }

        /* Status */
        .health-badge {
            min-height: 1.8rem;
            display: inline-flex;
            align-items: center;
            gap: 0.44rem;
            padding: 0.28rem 0.62rem;
            border: 1px solid var(--border);
            border-radius: 999px;
            color: var(--text-secondary);
            background: var(--surface-solid);
            font-size: 0.67rem;
            white-space: nowrap;
        }
        button.health-badge { cursor: pointer; }
        .status-dot { width: 0.44rem; height: 0.44rem; flex: 0 0 auto; border-radius: 50%; background: var(--success); box-shadow: 0 0 0 3px var(--success-soft); }
        .status-dot.warning { background: var(--caution); box-shadow: 0 0 0 3px var(--caution-soft); }
        .status-dot.boost { background: var(--boost); box-shadow: 0 0 0 3px var(--boost-soft); }
        .status-dot.neutral { background: var(--text-tertiary); box-shadow: none; }

        .badge {
            display: inline-flex;
            align-items: center;
            gap: 0.3rem;
            min-height: 1.35rem;
            padding: 0.12rem 0.5rem;
            border: 1px solid var(--border);
            border-radius: 999px;
            color: var(--text-secondary);
            background: var(--surface-raised);
            font-size: 0.62rem;
            font-weight: 640;
            line-height: 1.2;
            white-space: nowrap;
        }
        .badge-success { border-color: color-mix(in srgb, var(--success) 22%, var(--border)); color: var(--success); background: var(--success-soft); }
        .badge-warning, .badge-pending { border-color: color-mix(in srgb, var(--caution) 26%, var(--border)); color: var(--caution); background: var(--caution-soft); }
        .badge-danger { border-color: color-mix(in srgb, var(--danger) 24%, var(--border)); color: var(--danger); background: var(--danger-soft); }
        .badge-boost, .badge-predicted { border-color: color-mix(in srgb, var(--boost) 24%, var(--border)); color: var(--boost); background: var(--boost-soft); }
        .badge-subtle { color: var(--text-tertiary); background: var(--surface-subtle); border-color: transparent; }

        .chip { display: inline-flex; align-items: center; gap: 0.3rem; min-height: 1.6rem; padding: 0.18rem 0.55rem; border: 1px solid var(--border); border-radius: 999px; color: var(--text-secondary); background: var(--surface-raised); font-size: 0.66rem; white-space: nowrap; }
        .chip strong { color: var(--text); font-family: var(--mono); font-weight: 540; }

        /* Surfaces */
        .stat-grid { display: grid; grid-template-columns: repeat(4, minmax(0, 1fr)); gap: 0.72rem; margin-bottom: 1rem; }
        .stat-grid.three { grid-template-columns: repeat(3, minmax(0, 1fr)); }
        .stat-card { min-width: 0; padding: 1.05rem; border: 1px solid var(--border); border-radius: var(--radius-md); background: var(--surface-solid); box-shadow: 0 2px 8px rgba(0, 0, 0, 0.025); }
        .stat-card > span, .stat-card small { display: block; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
        .stat-card > span { color: var(--text-secondary); font-size: 0.68rem; }
        .stat-card strong { display: block; margin: 0.3rem 0 0.2rem; overflow: hidden; font-size: 1.72rem; line-height: 1.05; letter-spacing: -0.05em; text-overflow: ellipsis; white-space: nowrap; }
        .stat-card strong.text-value { padding: 0.32rem 0 0.12rem; font-size: 0.8rem; line-height: 1.35; letter-spacing: -0.01em; font-family: var(--mono); font-weight: 520; }
        .stat-card small { color: var(--text-tertiary); font-size: 0.61rem; }
        .stat-head { display: flex; align-items: center; justify-content: space-between; gap: 0.5rem; }
        .stat-head > span:first-child { min-width: 0; overflow: hidden; color: var(--text-secondary); font-size: 0.68rem; text-overflow: ellipsis; white-space: nowrap; }

        .card { margin-bottom: 1rem; overflow: hidden; border: 1px solid var(--border); border-radius: var(--radius-md); background: var(--surface-solid); box-shadow: 0 2px 9px rgba(0, 0, 0, 0.025); }
        .card-head { display: flex; align-items: flex-start; justify-content: space-between; gap: 1rem; padding: 1.05rem 1.15rem 0.88rem; border-bottom: 1px solid var(--border); background: color-mix(in srgb, var(--surface-subtle) 42%, transparent); }
        .card-head h2 { margin: 0; font-size: 0.86rem; letter-spacing: -0.012em; }
        .card-head p { margin: 0.3rem 0 0; color: var(--text-tertiary); font-size: 0.68rem; }
        .card-body { padding: 1.15rem; }
        .setting-row { min-height: 4.15rem; display: flex; align-items: center; justify-content: space-between; gap: 1.5rem; padding: 0.88rem 1.15rem; border-bottom: 1px solid var(--border); }
        .setting-row:last-child { border-bottom: 0; }
        .setting-copy { min-width: 0; }
        .setting-copy strong { display: block; font-size: 0.78rem; font-weight: 650; }
        .setting-copy span { display: block; margin-top: 0.18rem; color: var(--text-tertiary); font-size: 0.68rem; line-height: 1.4; }
        .setting-control { width: min(20rem, 48%); flex: 0 0 auto; display: flex; align-items: center; justify-content: flex-end; flex-wrap: wrap; gap: 0.45rem; color: var(--text-secondary); font-size: 0.72rem; text-align: right; overflow-wrap: anywhere; }
        .setting-control.wide { width: min(27rem, 60%); flex-wrap: nowrap; }
        .setting-control input[type="number"] { width: 8rem; text-align: right; font-family: var(--mono); }
        .setting-control input.time-input { width: 5.4rem; text-align: center; font-family: var(--mono); }
        .setting-control select { min-width: 12rem; }

        .default-notice { display: flex; align-items: flex-start; gap: 0.7rem; margin-bottom: 1rem; padding: 0.88rem 0.98rem; border: 1px solid color-mix(in srgb, var(--success) 20%, transparent); border-radius: var(--radius-sm); color: var(--success); background: var(--success-soft); font-size: 0.71rem; line-height: 1.52; }
        .default-notice svg { width: 1rem; height: 1rem; flex: 0 0 auto; margin-top: 0.08rem; }
        .card .default-notice { margin: 0; }
        .info-callout { display: flex; align-items: flex-start; gap: 0.9rem; padding: 1.15rem; border-bottom: 1px solid var(--border); }
        .info-callout:last-child { border-bottom: 0; }
        .callout-icon { width: 2.3rem; height: 2.3rem; flex: 0 0 auto; display: grid; place-items: center; border-radius: 0.72rem; color: var(--accent); background: var(--accent-soft); }
        .info-callout strong { display: block; font-size: 0.8rem; }
        .info-callout p { margin: 0.32rem 0 0; color: var(--text-secondary); font-size: 0.71rem; line-height: 1.58; }
        .mono { font-family: var(--mono); letter-spacing: -0.015em; }
        .muted { color: var(--text-tertiary); }
        .empty-state { grid-column: 1 / -1; padding: 2.6rem 1.25rem; border: 1px dashed var(--border-strong); border-radius: var(--radius-md); color: var(--text-tertiary); font-size: 0.75rem; text-align: center; }
        .empty-state.error { border-color: color-mix(in srgb, var(--danger) 30%, var(--border)); color: var(--danger); }

        /* Quota meters */
        .meter { display: grid; gap: 0.32rem; }
        .meter-label-row { display: flex; align-items: center; justify-content: space-between; gap: 0.6rem; color: var(--text-secondary); font-size: 0.66rem; font-weight: 600; }
        .meter-label-row strong { color: var(--text); font-family: var(--mono); font-weight: 540; }
        .meter-track {
            height: 0.42rem;
            overflow: hidden;
            border: 1px solid var(--meter-track-border);
            border-radius: 999px;
            background: var(--meter-track-bg);
        }
        .meter-fill { height: 100%; border-radius: inherit; transition: width 400ms var(--ease-out), background-color 300ms ease; }
        .meter-fill-healthy { background: var(--success); }
        .meter-fill-warning { background: var(--caution); }
        .meter-fill-danger { background: var(--danger); }
        .meter-countdown { color: var(--text-tertiary); font-family: var(--mono); font-size: 0.64rem; font-weight: 500; font-variant-numeric: tabular-nums; }
        .mini-meter { display: flex; align-items: center; gap: 0.5rem; min-width: 8.5rem; max-width: 11rem; }
        .mini-meter .meter-track { flex: 1; min-width: 5rem; }
        .mini-meter strong { min-width: 2.3rem; font-family: var(--mono); font-size: 0.68rem; text-align: right; }

`
