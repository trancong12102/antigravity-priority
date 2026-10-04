package management

// Stroke icons shared by the static markup and the rendering scripts. They use
// the same 24px outline style as key-model-access; colour comes from
// currentColor so every icon follows the active theme.
const (
	iconBrand    = `<svg viewBox="0 0 24 24"><path d="M12 20V8M7.5 12.5 12 8l4.5 4.5"/><path d="M5 4h14"/></svg>`
	iconLock     = `<svg viewBox="0 0 24 24"><path d="M7.75 10.25V7.5a4.25 4.25 0 0 1 8.5 0v2.75M6.5 10.25h11a1.5 1.5 0 0 1 1.5 1.5v6.75a1.5 1.5 0 0 1-1.5 1.5h-11A1.5 1.5 0 0 1 5 18.5v-6.75a1.5 1.5 0 0 1 1.5-1.5Z"/><path d="M12 14v2.25"/></svg>`
	iconOverview = `<svg viewBox="0 0 24 24"><rect x="4" y="4" width="6" height="6" rx="1"/><rect x="14" y="4" width="6" height="6" rx="1"/><rect x="4" y="14" width="6" height="6" rx="1"/><rect x="14" y="14" width="6" height="6" rx="1"/></svg>`
	iconHistory  = `<svg viewBox="0 0 24 24"><path d="M3.5 12a8.5 8.5 0 1 0 2.6-6.1L3.5 8.5"/><path d="M3.5 4v4.5H8M12 7.5V12l3 2"/></svg>`
	iconPulse    = `<svg viewBox="0 0 24 24"><path d="M3 12h4l2.5-6.5 5 13L17 12h4"/></svg>`
	iconSliders  = `<svg viewBox="0 0 24 24"><path d="M4 7h9M17 7h3M4 17h3M11 17h9"/><circle cx="15" cy="7" r="2"/><circle cx="9" cy="17" r="2"/></svg>`
	iconHelp     = `<svg viewBox="0 0 24 24"><circle cx="12" cy="12" r="8.5"/><path d="M9.6 9.6a2.5 2.5 0 1 1 3.4 2.3c-.6.3-1 .8-1 1.5v.4M12 16.8h.01"/></svg>`
	iconRefresh  = `<svg viewBox="0 0 24 24"><path d="M20 6v5h-5M4 18v-5h5"/><path d="M6.1 8.2A7 7 0 0 1 18.8 9M17.9 15.8A7 7 0 0 1 5.2 15"/></svg>`
	iconProbe    = `<svg viewBox="0 0 24 24"><path d="M4.9 9.5a8 8 0 0 1 14.2 0M7.8 11.2a4.6 4.6 0 0 1 8.4 0"/><circle cx="12" cy="13.5" r="1.3"/><path d="M12 15v5"/></svg>`
	iconBolt     = `<svg viewBox="0 0 24 24"><path d="M13 3 5.5 13.5H11L10 21l7.5-10.5H12L13 3Z"/></svg>`
	iconReset    = `<svg viewBox="0 0 24 24"><path d="M4 4.5v5h5"/><path d="M5.4 14.5A7 7 0 1 0 6.6 7.2L4 9.5"/></svg>`
	iconCopy     = `<svg viewBox="0 0 24 24"><rect x="8.5" y="8.5" width="11" height="11" rx="2"/><path d="M15.5 8.5V6a1.5 1.5 0 0 0-1.5-1.5H6A1.5 1.5 0 0 0 4.5 6v8A1.5 1.5 0 0 0 6 15.5h2.5"/></svg>`
	iconSave     = `<svg viewBox="0 0 24 24"><path d="M5 4h12l2 2v14H5z"/><path d="M8 4v6h8V4M8 20v-6h8v6"/></svg>`
	iconClose    = `<svg viewBox="0 0 24 24"><path d="m7 7 10 10M17 7 7 17"/></svg>`
	iconWarning  = `<svg viewBox="0 0 24 24"><path d="M12 3 2.8 19h18.4L12 3Z"/><path d="M12 9v4M12 16.5h.01"/></svg>`
	iconCheck    = `<svg viewBox="0 0 24 24"><path d="m5 12 4 4L19 6"/></svg>`
	iconInfo     = `<svg viewBox="0 0 24 24"><circle cx="12" cy="12" r="8.5"/><path d="M12 11v5M12 8h.01"/></svg>`
	iconChart    = `<svg viewBox="0 0 24 24"><path d="M5 20V11M12 20V5M19 20v-6"/></svg>`
	iconCalendar = `<svg viewBox="0 0 24 24"><rect x="4" y="5" width="16" height="15" rx="2"/><path d="M4 10h16M9 3v4M15 3v4"/></svg>`
	iconShield   = `<svg viewBox="0 0 24 24"><path d="M12 3 5 6v5c0 4.4 3 8 7 10 4-2 7-5.6 7-10V6l-7-3Z"/></svg>`
	iconTrend    = `<svg viewBox="0 0 24 24"><path d="M4 17l5-5 4 4 7-7"/><path d="M15 9h5v5"/></svg>`
	iconScale    = `<svg viewBox="0 0 24 24"><path d="M12 4v16M7 20h10M5 8h14"/><path d="m5 8-2.5 6a3 3 0 0 0 5 0L5 8ZM19 8l-2.5 6a3 3 0 0 0 5 0L19 8Z"/></svg>`
	iconArrowUp  = `<svg viewBox="0 0 24 24"><path d="M12 19V5M6 11l6-6 6 6"/></svg>`
	iconSearch   = `<svg viewBox="0 0 24 24"><circle cx="11" cy="11" r="6.5"/><path d="m16 16 4 4"/></svg>`
	iconSpinner  = `<svg class="spinner" viewBox="0 0 24 24"><path d="M21 12a9 9 0 0 1-9 9"/><path d="M3 12a9 9 0 0 1 9-9" opacity=".35"/></svg>`
)

const templateScriptIcons = `        const ICONS = {
            probe: '` + iconProbe + `',
            chart: '` + iconChart + `',
            warning: '` + iconWarning + `',
            check: '` + iconCheck + `',
            info: '` + iconInfo + `',
            search: '` + iconSearch + `',
            spinner: '` + iconSpinner + `',
            lock: '` + iconLock + `'
        };

`
