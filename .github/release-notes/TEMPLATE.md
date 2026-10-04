# antigravity-priority {version}

## Release Notes

- Write short, user-facing bullets: what changed and how it affects day-to-day use.
- Write in English only.
- Avoid internal symbol names, struct fields, or source paths.
- Never include internal `(REQ-xx)` tags.
- For config changes: name the setting, default behavior, and whether users must edit config.
- For behavior changes: who is affected and who is not.
- Do not add a separate "Upgrade notes" section; replace the plugin binary and restart the host as usual.

**This update ({version}) includes:**
- Exclusively tailored for Google Antigravity quota management with dual-window (5h burst + 7d weekly) scheduling.
- Implements Weekly Urgency Index and Dynamic Boost Horizon algorithms to eliminate end-of-cycle quota waste.
- Adaptive cycle burn-rate estimation ($C_{\text{cycle}}$) learns actual consumption patterns without manual tuning.
- Self-healing soft depletion for 5h window exhaustion and hard disabling for 7d weekly exhaustion.
- Native CPA light/dark theme embedded web management dashboard (`/status`).
