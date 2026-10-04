package management

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

func TestManagementPageAssetContract(t *testing.T) {
	expectedNames := []string{"Overview", "History", "Diagnostics", "Config", "Help"}
	if len(managementPageFeatures) != len(expectedNames) {
		t.Fatalf("expected %d feature assets, got %d", len(expectedNames), len(managementPageFeatures))
	}

	for index, feature := range managementPageFeatures {
		if feature.name != expectedNames[index] {
			t.Errorf("feature %d is %q, want %q", index, feature.name, expectedNames[index])
		}
		if strings.TrimSpace(feature.markup) == "" {
			t.Errorf("%s feature has no markup asset", feature.name)
		}
	}

	if got := strings.Count(StatusHTML, "<section id=\"panel"); got != len(expectedNames) {
		t.Fatalf("assembled page contains %d feature panels, want %d", got, len(expectedNames))
	}
	if StatusHTML != assembleManagementPage() {
		t.Fatal("StatusHTML is not the deterministic output of the page assembly seam")
	}
}

func TestOverviewUsesCPAEmailAsCredentialTitle(t *testing.T) {
	if !strings.Contains(StatusHTML, `const credDisplayName = item.email || "Credential";`) {
		t.Fatal("overview credential title must use the CPA Host email")
	}
	if strings.Contains(StatusHTML, `item.name || item.account || item.auth_index`) {
		t.Fatal("overview credential title still uses a non-email fallback chain")
	}
}

func TestOverview_SingleProbeButtonAndStatusDisabledBadge(t *testing.T) {
	if !strings.Contains(StatusHTML, `probeSingleCredential`) {
		t.Fatal("overview must include probeSingleCredential script function")
	}
	if !strings.Contains(StatusHTML, `data-credential-action=\"probe\" data-auth-index=\"" + escapeHTML(authIdx)`) {
		t.Fatal("overview card must render single-credential probe button")
	}
	if !strings.Contains(StatusHTML, `statusBadge = "<span class=\"badge badge-subtle\">Disabled</span>";`) {
		t.Fatal("overview card must render statusDisabled badge for manually disabled accounts")
	}
	if !strings.Contains(StatusHTML, `item.target && item.target.disabled && !isManuallyDisabled`) {
		t.Fatal("overview depleted count must exclude manually disabled accounts")
	}
}

// Server data must never be spliced into an inline event handler: HTML
// entity escaping is undone before the handler's JavaScript is parsed, so a
// quote in an auth index or email would break out of the string literal.
func TestStatusHTML_NoServerDataInInlineHandlers(t *testing.T) {
	inlineHandler := regexp.MustCompile(`on[a-z]+=\\?"[^"]*"\s*\+`)
	if match := inlineHandler.FindString(StatusHTML); match != "" {
		t.Fatalf("inline event handler built by string concatenation: %s", match)
	}
}

func TestStatusHTML_EnglishOnly(t *testing.T) {
	if hanRE := regexp.MustCompile(`\p{Han}`); hanRE.MatchString(StatusHTML) {
		t.Fatalf("assembled page contains Han characters: %q", hanRE.FindAllString(StatusHTML, 10))
	}
	if !strings.Contains(StatusHTML, `<html lang="en">`) {
		t.Error("document language must be English")
	}
	for _, leftover := range []string{"data-i18n", "I18N", "currentLang", "toggleLanguage", "applyLanguage", "LANG_STORAGE_KEY", "langLabel", "zh-CN", "en-US\"]"} {
		if strings.Contains(StatusHTML, leftover) {
			t.Errorf("assembled page still carries language-switch scaffolding %q", leftover)
		}
	}
	if regexp.MustCompile(`\bt\(\s*["']`).MatchString(extractPageScript(t)) {
		t.Error("script still resolves strings through a translation lookup")
	}

	// Static markup must already be the final English text so nothing is
	// swapped in after load.
	markup := StatusHTML[:strings.LastIndex(StatusHTML, "<script>")]
	for _, text := range []string{
		"Connecting to management session",
		"Antigravity credentials",
		"Total credentials",
		"Fetch quota",
		"Reset priority",
		"Recent runs",
		"Scheduler and health",
		"Auto periodic scheduling",
		"Reset to defaults",
		"How priorities are computed",
	} {
		if !strings.Contains(markup, text) {
			t.Errorf("static markup is missing English text %q", text)
		}
	}
}

func TestStatusHTML_InlineHandlersResolve(t *testing.T) {
	script := extractPageScript(t)
	functionNames := map[string]struct{}{}
	functionRE := regexp.MustCompile(`\bfunction\s+([A-Za-z_$][A-Za-z0-9_$]*)\s*\(`)
	for _, match := range functionRE.FindAllStringSubmatch(script, -1) {
		functionNames[match[1]] = struct{}{}
	}

	callRE := regexp.MustCompile(`\b([A-Za-z_$][A-Za-z0-9_$]*)\s*\(`)
	handlerRE := regexp.MustCompile(`(?i)\bon[a-z]+\s*=\s*"([^"]+)"`)
	for _, match := range handlerRE.FindAllStringSubmatch(StatusHTML, -1) {
		for _, call := range callRE.FindAllStringSubmatch(match[1], -1) {
			if _, ok := functionNames[call[1]]; !ok {
				t.Errorf("inline handler %q calls unavailable function %q", match[1], call[1])
			}
		}
	}
}

func TestStatusHTML_StaticDOMReferencesResolve(t *testing.T) {
	idNames := map[string]struct{}{}
	idRE := regexp.MustCompile(`\bid="([A-Za-z][A-Za-z0-9_-]*)"`)
	for _, match := range idRE.FindAllStringSubmatch(StatusHTML, -1) {
		idNames[match[1]] = struct{}{}
	}

	script := extractPageScript(t)
	getElementRE := regexp.MustCompile(`document\.getElementById\(\s*["']([A-Za-z][A-Za-z0-9_-]*)["']\s*\)`)
	selectorRE := regexp.MustCompile(`document\.querySelector(?:All)?\(\s*["']#([A-Za-z][A-Za-z0-9_-]*)["']\s*\)`)
	for _, expression := range []*regexp.Regexp{getElementRE, selectorRE} {
		for _, match := range expression.FindAllStringSubmatch(script, -1) {
			if _, ok := idNames[match[1]]; !ok {
				t.Errorf("script references missing assembled DOM id %q", match[1])
			}
		}
	}
}

func TestStatusHTML_OperationalWorkflowContract(t *testing.T) {
	workflowSnippets := []string{
		"function triggerProbe()",
		"RUN_PATH + \"?mode=probe",
		"function triggerApplyWithConfirm()",
		"showModal(\"apply-confirm\"",
		"function executeDirectApply()",
		"RUN_PATH + \"?mode=apply",
		"function triggerReset()",
		"RESET_PATH",
		"function fetchDynamicConfig()",
		"function saveDynamicConfig()",
		"CONFIG_PATH",
		"function resetDynamicConfigToDefaults()",
		"Reset configuration to defaults?",
	}
	for _, snippet := range workflowSnippets {
		if !strings.Contains(StatusHTML, snippet) {
			t.Errorf("assembled page is missing workflow contract %q", snippet)
		}
	}

	applyConfirm := strings.Index(StatusHTML, `showModal("apply-confirm"`)
	directApply := strings.Index(StatusHTML, `RUN_PATH + "?mode=apply`)
	if applyConfirm == -1 || directApply == -1 || applyConfirm >= directApply {
		t.Fatal("assembled apply workflow does not expose confirmation before direct apply")
	}
}

func TestStatusHTML_JavaScriptSyntax(t *testing.T) {
	node, err := exec.LookPath("node")
	if err != nil {
		t.Skip("node is not available; JavaScript syntax is checked in the Node-enabled verification environment")
	}

	scriptPath := filepath.Join(t.TempDir(), "management-page.js")
	if err := os.WriteFile(scriptPath, []byte(extractPageScript(t)), 0o600); err != nil {
		t.Fatalf("write JavaScript fixture: %v", err)
	}
	cmd := exec.Command(node, "--check", scriptPath)
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("assembled management-page script failed node --check: %v\n%s", err, output)
	}
}

func TestApplyPreviewWorkflow(t *testing.T) {
	node, err := exec.LookPath("node")
	if err != nil {
		t.Skip("node is not available; apply preview behavior is checked in the Node-enabled verification environment")
	}

	testCases := []struct {
		name       string
		changes    string
		items      string
		wantModal  bool
		wantMode   string
		wantCount  int
		wantToasts int
	}{
		{
			name:      "projected target change is previewed",
			changes:   `[]`,
			items:     `[{auth_index:"credential-1",name:"Credential 1",current:{priority:99,priority_missing:false,disabled:false},target:{priority:98,priority_missing:false,disabled:false}}]`,
			wantModal: true,
			wantMode:  "projected",
			wantCount: 1,
		},
		{
			name:       "synchronized state reports optimal",
			changes:    `[]`,
			items:      `[{auth_index:"credential-1",name:"Credential 1",current:{priority:99,priority_missing:false,disabled:false},target:{priority:99,priority_missing:false,disabled:false}}]`,
			wantToasts: 1,
		},
		{
			name:      "write qualified changes retain pending preview",
			changes:   `[{auth_index:"credential-1",name:"Credential 1",current:{priority:99,priority_missing:false,disabled:false},target:{priority:98,priority_missing:false,disabled:false}}]`,
			items:     `[{auth_index:"credential-1",name:"Credential 1",current:{priority:99,priority_missing:false,disabled:false},target:{priority:98,priority_missing:false,disabled:false}}]`,
			wantModal: true,
			wantMode:  "pending",
			wantCount: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			harness := templateScriptOverviewActionsCore + templateScriptModals + `
const toastCalls = [];
function element(initial) {
    return Object.assign({hidden:false,textContent:"",innerHTML:"",children:[],appendChild:function(child){this.children.push(child);}}, initial || {});
}
const elements = {
    modelGroupSelect: element({value:"gemini"}),
    diffModal: element({hidden:true}),
    modalTitle: element(),
    modalSummary: element(),
    modalDiffList: element(),
    btnModalApply: element()
};
const document = {
    getElementById: function(id) {
        return elements[id] || null;
    },
    createElement: function() { return element(); }
};
function showToast(message, kind) { toastCalls.push({message, kind}); }
function escapeHTML(value) { return String(value); }
dynamicConfig = { antigravity_model_group: "gemini" };
latestSnapshot = { active_model_group: "gemini", groups: { gemini: { changes: ` + tc.changes + `, items: ` + tc.items + ` } } };

(async function() {
    await triggerApplyWithConfirm();
    const expected = ` + applyPreviewExpectationJSON(tc.wantModal, tc.wantMode, tc.wantCount, tc.wantToasts) + `;
    if (elements.diffModal.hidden === expected.modal) {
        throw new Error("unexpected modal visibility: " + elements.diffModal.hidden);
    }
    if (toastCalls.length !== expected.toasts) {
        throw new Error("unexpected toast count: " + toastCalls.length);
    }
    if (expected.modal) {
        const expectedSummaryKey = expected.mode === "projected" ? "Current target preview" : "Credentials to update";
        if (elements.modalSummary.textContent.indexOf(expectedSummaryKey) !== 0) {
            throw new Error("unexpected preview summary: " + elements.modalSummary.textContent);
        }
        if (elements.modalDiffList.children.length !== expected.count) {
            throw new Error("unexpected rendered preview count: " + elements.modalDiffList.children.length);
        }
    }
})().catch(function(err) {
    console.error(err.stack || err.message);
    process.exit(1);
});
`
			runNodeFixture(t, node, "apply-preview.js", harness)
		})
	}
}

func TestSamplesModalUsesSelectedModelGroupOnly(t *testing.T) {
	node, err := exec.LookPath("node")
	if err != nil {
		t.Skip("node is not available; samples modal behavior is checked in the Node-enabled verification environment")
	}

	harness := templateScriptOverviewSamples + `
function element(initial) {
    return Object.assign({hidden:true,textContent:"",innerHTML:""}, initial || {});
}
const elements = {
    samplesModal: element(),
    samplesModalTitle: element(),
    samplesModalSubtitle: element(),
    samplesModalBody: element(),
    modelGroupSelect: element({value:"claude_gpt"})
};
const document = {
    getElementById: function(id) { return elements[id] || null; }
};
const SAMPLES_PATH = "/samples";
function escapeHTML(value) { return String(value); }
function formatDateTime(value) { return value ? new Date(value).toISOString() : "-"; }
function quotaFillClass(percent) { return percent <= 10 ? "meter-fill-danger" : "meter-fill-healthy"; }
async function apiFetch() {
    return {
        groups: {
            gemini: {name:"Gemini models", samples:[{observed_at:"2026-08-23T02:00:00Z",short_window_rem:100,long_window_rem:100}]},
            claude_gpt: {name:"Claude & GPT models", samples:[
                {observed_at:"2026-08-23T02:00:00Z",short_window_rem:0,long_window_rem:45},
                {observed_at:"2026-08-23T02:15:00Z",short_window_rem:100,long_window_rem:45}
            ]}
        }
    };
}
(async function() {
    await openSamplesModal("auth-test", "Test Account");
    const html = elements.samplesModalBody.innerHTML;
    if (html.indexOf("Gemini") >= 0) throw new Error("selected Claude/GPT view rendered Gemini data");
    if ((html.match(/class="sample-group"/g) || []).length !== 2) throw new Error("selected group sample count mismatch");
    if (html.indexOf("0%") < 0) throw new Error("zero quota was not rendered as 0%");
    if (html.indexOf("100%") > html.indexOf("0%")) throw new Error("samples are not sorted newest first");
})().catch(function(err) {
    console.error(err.stack || err.message);
    process.exit(1);
});
`
	runNodeFixture(t, node, "samples-modal.js", harness)
}

func TestProbeHistoryDetailsRenderBothModelGroups(t *testing.T) {
	node, err := exec.LookPath("node")
	if err != nil {
		t.Skip("node is not available; probe history behavior is checked in the Node-enabled verification environment")
	}

	harness := templateScriptHistoryDetails + templateScriptModals + `
const calls = [];
function element(initial) {
    return Object.assign({hidden:true,textContent:"",innerHTML:"",children:[],appendChild:function(child){this.children.push(child);}}, initial || {});
}
const elements = {
    diffModal: element(),
    modalTitle: element(),
    modalSummary: element(),
    modalDiffList: element(),
    btnModalApply: element()
};
const document = {
    getElementById: function(id) { return elements[id] || null; },
    createElement: function() { return element(); }
};
const SAMPLES_PATH = "/samples";
const latestDiagnostics = {run_history:[{kind:"probe",probe_round_id:"round-1"}]};
function showToast(message, kind) { throw new Error(message + " (" + kind + ")"); }
function escapeHTML(value) { return String(value); }
function apiFetch(path) {
    calls.push(path);
    if (path.indexOf("model_group=gemini") >= 0) {
        return Promise.resolve({records:[{email:"gemini@example.com",sample:{observed_at:"2026-08-23T02:00:00Z",short_window_rem:0,long_window_rem:25}}]});
    }
    return Promise.resolve({records:[{email:"claude@example.com",sample:{observed_at:"2026-08-23T02:01:00Z",short_window_rem:75,long_window_rem:0}}]});
}
(async function() {
    showHistoryDetails(0);
    await new Promise(function(resolve) { setTimeout(resolve, 0); });
    if (calls.length !== 2) throw new Error("probe history did not request both model groups");
    if (!calls.some(function(path) { return path.indexOf("model_group=gemini") >= 0; })) throw new Error("Gemini probe samples were not requested");
    if (!calls.some(function(path) { return path.indexOf("model_group=claude_gpt") >= 0; })) throw new Error("Claude/GPT probe samples were not requested");
    if (elements.modalDiffList.innerHTML.indexOf("history-comparison-table") < 0) throw new Error("probe history comparison table is missing");
    if (elements.modalDiffList.innerHTML.indexOf("history-comparison-quota-only") < 0) throw new Error("probe history must use the quota-only table layout");
    if (elements.modalDiffList.innerHTML.indexOf("Gemini quota this probe") < 0 || elements.modalDiffList.innerHTML.indexOf("Claude/GPT quota this probe") < 0) throw new Error("probe history model-group columns are missing");
    if (elements.modalDiffList.innerHTML.indexOf("Priority change") >= 0) throw new Error("probe history must not render a priority column");
    if (elements.modalDiffList.innerHTML.indexOf("gemini@example.com") < 0 || elements.modalDiffList.innerHTML.indexOf("claude@example.com") < 0) throw new Error("probe history accounts are missing");
})().catch(function(err) {
    console.error(err.stack || err.message);
    process.exit(1);
});
`
	runNodeFixture(t, node, "probe-history.js", harness)
}

func TestAutoScheduleHistoryCombinesProbeAndHostDetails(t *testing.T) {
	node, err := exec.LookPath("node")
	if err != nil {
		t.Skip("node is not available; automatic schedule history is checked in the Node-enabled verification environment")
	}

	harness := templateScriptHistoryDetails + templateScriptModals + `
function element(initial) {
    return Object.assign({hidden:true,textContent:"",innerHTML:"",children:[],appendChild:function(child){this.children.push(child);}}, initial || {});
}
const elements = {
    diffModal: element(),
    modalTitle: element(),
    modalSummary: element(),
    modalDiffList: element(),
    btnModalApply: element()
};
const document = {
    getElementById: function(id) { return elements[id] || null; },
    createElement: function() { return element(); }
};
const SAMPLES_PATH = "/samples";
const latestDiagnostics = {run_history:[{
    kind:"auto_apply",
    trigger:"auto_apply",
    probe_round_id:"round-auto",
    snapshot:{changes:[{email:"account@example.com",current:{priority:98,disabled:false},target:{priority:99,disabled:false},reason:"weekly urgency"}]},
    attempted:1,
    succeeded:1
}]};
function showToast(message, kind) { throw new Error(message + " (" + kind + ")"); }
function escapeHTML(value) { return String(value); }
function apiFetch(path) {
    return Promise.resolve({records:[{email:"account@example.com",sample:{observed_at:"2026-08-24T02:00:00Z",short_window_rem:90,long_window_rem:55}}]});
}
function extractChanges(result) { return result && Array.isArray(result.changes) ? result.changes : []; }
function formatReason(reason) { return reason || ""; }
(async function() {
    showHistoryDetails(0);
    await new Promise(function(resolve) { setTimeout(resolve, 0); });
    if (elements.modalTitle.textContent !== "Automatic schedule details") throw new Error("automatic schedule title is missing");
    if (elements.modalDiffList.innerHTML.indexOf("history-comparison-table") < 0) throw new Error("automatic schedule comparison table is missing");
    if (elements.modalDiffList.innerHTML.indexOf("Priority change") < 0) throw new Error("automatic schedule priority column is missing");
    if (elements.modalDiffList.innerHTML.indexOf("account@example.com") < 0) throw new Error("automatic schedule account is missing");
    if ((elements.modalDiffList.innerHTML.match(/history-comparison-row/g) || []).length !== 2) throw new Error("Gemini, Claude/GPT, and priority data were not merged into one account row");
    if (elements.modalDiffList.innerHTML.indexOf("weekly urgency") < 0) throw new Error("automatic schedule priority details are missing");
})().catch(function(err) {
    console.error(err.stack || err.message);
    process.exit(1);
});
`
	runNodeFixture(t, node, "auto-schedule-history.js", harness)
}

func TestStatusHTML_NeverPromptsForOrPersistsManagementKey(t *testing.T) {
	for _, forbidden := range []string{
		"keyModal", "manualKeyInput", "openKeyModal", "saveKeyAndRefresh", "setSavedKey",
		"management_key", "management-key", "URLSearchParams", "sessionStorage", "localStorage.setItem",
		`type="password"`,
	} {
		if strings.Contains(StatusHTML, forbidden) {
			t.Errorf("page must not prompt for, read from the URL, or persist the Management Key (found %q)", forbidden)
		}
	}
	for _, required := range []string{`"cli-proxy-auth"`, `"enc::v1::"`, `"cli-proxy-api-webui::secure-storage"`, `localStorage.getItem("isLoggedIn") !== "true"`} {
		if !strings.Contains(StatusHTML, required) {
			t.Errorf("page must read the CPAMC session like key-model-access (missing %q)", required)
		}
	}
}

func TestStatusHTML_CPAMCSessionGate(t *testing.T) {
	node, err := exec.LookPath("node")
	if err != nil {
		t.Skip("node is not available; the CPAMC session gate is checked in the Node-enabled verification environment")
	}

	harness := `
function makeElement(id) {
    const classes = new Set();
    const attributes = {};
    return {
        id: id, hidden: false, textContent: "", innerHTML: "", value: "", disabled: false, title: "", className: "",
        dataset: {}, style: {}, children: [],
        classList: {
            add: function(name) { classes.add(name); },
            remove: function(name) { classes.delete(name); },
            contains: function(name) { return classes.has(name); },
            toggle: function(name, force) { const on = force === undefined ? !classes.has(name) : Boolean(force); if (on) classes.add(name); else classes.delete(name); return on; }
        },
        setAttribute: function(name, value) { attributes[name] = String(value); },
        getAttribute: function(name) { return Object.prototype.hasOwnProperty.call(attributes, name) ? attributes[name] : null; },
        removeAttribute: function(name) { delete attributes[name]; },
        appendChild: function(child) { this.children.push(child); },
        querySelector: function() { return null; },
        querySelectorAll: function() { return []; },
        addEventListener: function() {},
        remove: function() {},
        focus: function() {}
    };
}
const elements = {};
const document = {
    hidden: false,
    documentElement: makeElement("html"),
    body: makeElement("body"),
    getElementById: function(id) {
        if (!elements[id]) elements[id] = makeElement(id);
        return elements[id];
    },
    querySelector: function() { return null; },
    querySelectorAll: function() { return []; },
    addEventListener: function() {},
    createElement: function() { return makeElement(""); }
};
document.getElementById("app").hidden = true;
document.getElementById("retrySessionButton").hidden = true;

const storageListeners = [];
const window = {
    location: { host: "cpa.example:8317" },
    matchMedia: function() { return { matches: false }; },
    addEventListener: function(type, listener) { if (type === "storage") storageListeners.push(listener); },
    scrollTo: function() {}
};
window.self = window;
window.top = window;
const navigator = { userAgent: "SessionGateTest/1.0" };
function requestAnimationFrame(callback) { callback(); }

const store = new Map();
const localStorage = {
    getItem: function(name) { return store.has(name) ? store.get(name) : null; },
    setItem: function(name, value) { store.set(name, String(value)); },
    removeItem: function(name) { store.delete(name); }
};

// Encodes a value exactly like CPAMC secureStorage v1.
function cpamcObfuscate(value) {
    const key = new TextEncoder().encode("cli-proxy-api-webui::secure-storage|" + window.location.host + "|" + navigator.userAgent);
    const plain = new TextEncoder().encode(value);
    let binary = "";
    for (let index = 0; index < plain.length; index += 1) binary += String.fromCharCode(plain[index] ^ key[index % key.length]);
    return "enc::v1::" + btoa(binary);
}

const requests = [];
let respondWith = 200;
global.fetch = async function(path, options) {
    requests.push({ path: path, headers: (options && options.headers) || {} });
    const status = respondWith;
    return {
        status: status,
        ok: status >= 200 && status < 300,
        statusText: status === 401 ? "Unauthorized" : "OK",
        text: async function() { return status === 401 ? "{\"error\":\"invalid management key\"}" : JSON.stringify({ paused: false, active_model_group: "gemini" }); }
    };
};
function settle() { return new Promise(function(resolve) { setTimeout(resolve, 30); }); }
function assert(condition, message) { if (!condition) throw new Error(message); }
function gateVisible() { return elements.authGate && !elements.authGate.hidden && elements.app.hidden; }
` + extractPageScript(t) + `
(async function() {
    // 1. No CPAMC session: the gate explains what to do, no request is sent and no key is asked for.
    await settle();
    assert(requests.length === 0, "no request may be sent without a CPAMC session, got " + requests.length);
    assert(gateVisible(), "the session gate must be shown without a CPAMC session");
    assert(elements.authTitle.textContent === "No reusable CPAMC session found", "unexpected gate title: " + elements.authTitle.textContent);
    assert(elements.authMessage.textContent.indexOf("Remember password") >= 0, "gate must explain the Remember password requirement");
    assert(elements.retrySessionButton.hidden === false, "gate must offer a retry button");

    // 2. Logged in but not "Remember password" style session: isLoggedIn alone is not enough.
    store.set("isLoggedIn", "false");
    store.set("cli-proxy-auth", cpamcObfuscate(JSON.stringify({ state: { managementKey: "ignored" } })));
    assert(readCPAMCManagementKey() === "", "a logged-out CPAMC session must not yield a key");

    // 3. Obfuscated CPAMC session: requests carry the key, the dashboard opens, nothing is written to storage.
    store.set("isLoggedIn", "true");
    store.set("cli-proxy-auth", cpamcObfuscate(JSON.stringify({ state: { managementKey: " secret-key " }, version: 0 })));
    const storedBefore = JSON.stringify(Array.from(store.entries()));
    await startSession();
    await settle();
    assert(requests.length > 0, "a CPAMC session must trigger authenticated requests");
    requests.forEach(function(request) {
        assert(request.headers.Authorization === "Bearer secret-key", "request " + request.path + " lacks the Bearer key");
        assert(request.headers["X-Management-Key"] === "secret-key", "request " + request.path + " lacks X-Management-Key");
    });
    assert(!elements.app.hidden && elements.authGate.hidden, "the dashboard must open once the session is verified");
    assert(JSON.stringify(Array.from(store.entries())) === storedBefore, "the page must not write to localStorage");

    // 4. Legacy plain-text key is still honoured.
    store.set("cli-proxy-auth", "{}");
    store.set("managementKey", JSON.stringify("legacy-key"));
    assert(readCPAMCManagementKey() === "legacy-key", "legacy CPAMC managementKey must be read");
    store.delete("managementKey");
    store.set("cli-proxy-auth", cpamcObfuscate(JSON.stringify({ state: { managementKey: "secret-key" } })));

    // 5. CPA rejects the key: gate returns, and nothing else hits the network while blocked.
    respondWith = 401;
    requests.length = 0;
    try { await apiFetch(SNAPSHOT_PATH); } catch (_) {}
    assert(isAuthBlocked, "a 401 must block further requests");
    assert(gateVisible(), "a 401 must bring the session gate back");
    assert(elements.authTitle.textContent === "Management session unavailable", "unexpected 401 gate title: " + elements.authTitle.textContent);
    const afterReject = requests.length;
    try { await apiFetch(SNAPSHOT_PATH); } catch (_) {}
    switchTab("diagnostics");
    await settle();
    assert(requests.length === afterReject, "blocked session must not issue requests");

    // 6. CPAMC logs in again in another document: the page reconnects by itself.
    respondWith = 200;
    store.set("cli-proxy-auth", cpamcObfuscate(JSON.stringify({ state: { managementKey: "fresh-key" } })));
    storageListeners.forEach(function(listener) { listener({ key: "cli-proxy-auth" }); });
    await settle();
    assert(!isAuthBlocked && !elements.app.hidden, "a renewed CPAMC session must reconnect automatically");
    assert(requests[requests.length - 1].headers.Authorization === "Bearer fresh-key", "reconnect must use the renewed key");

    // 7. CPAMC logs out: the page locks without waiting for a failed request.
    store.delete("isLoggedIn");
    storageListeners.forEach(function(listener) { listener({ key: "isLoggedIn" }); });
    assert(gateVisible(), "logging out of CPAMC must lock the page");
    assert(elements.authTitle.textContent === "CPAMC session ended", "unexpected logout gate title: " + elements.authTitle.textContent);

    // 8. Theme follows CPAMC's stored theme and storage updates.
    store.set("cli-proxy-theme", JSON.stringify({ state: { theme: "dark" } }));
    storageListeners.forEach(function(listener) { listener({ key: "cli-proxy-theme" }); });
    assert(document.documentElement.getAttribute("data-theme") === "dark", "dark CPAMC theme must apply");
    store.set("cli-proxy-theme", JSON.stringify({ state: { theme: "light" } }));
    syncTheme();
    assert(document.documentElement.getAttribute("data-theme") === null, "light CPAMC theme uses the default palette");
    process.exit(0);
})().catch(function(err) {
    console.error(err.stack || err.message);
    process.exit(1);
});
`
	runNodeFixture(t, node, "cpamc-session-gate.js", harness)
}

func applyPreviewExpectationJSON(wantModal bool, wantMode string, wantCount, wantToasts int) string {
	return fmt.Sprintf(`{modal:%t,mode:%q,count:%d,toasts:%d}`, wantModal, wantMode, wantCount, wantToasts)
}

func runNodeFixture(t *testing.T, node, filename, script string) {
	t.Helper()
	scriptPath := filepath.Join(t.TempDir(), filename)
	if err := os.WriteFile(scriptPath, []byte(script), 0o600); err != nil {
		t.Fatalf("write JavaScript fixture: %v", err)
	}
	cmd := exec.Command(node, scriptPath)
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("JavaScript fixture failed: %v\n%s", err, output)
	}
}

// extractPageScript returns every inline script of the page in document order.
func extractPageScript(t *testing.T) string {
	t.Helper()
	var scripts strings.Builder
	rest := StatusHTML
	for {
		start := strings.Index(rest, "<script>")
		if start == -1 {
			break
		}
		end := strings.Index(rest[start:], "</script>")
		if end == -1 {
			t.Fatal("assembled page contains an unterminated inline script")
		}
		scripts.WriteString(rest[start+len("<script>") : start+end])
		scripts.WriteString("\n")
		rest = rest[start+end+len("</script>"):]
	}
	if scripts.Len() == 0 {
		t.Fatal("assembled page does not contain an inline script")
	}
	return scripts.String()
}

func TestStatusHTML_ThemeResolvedBeforeFirstPaint(t *testing.T) {
	head := StatusHTML[:strings.Index(StatusHTML, "<body>")]
	themeScript := strings.Index(head, templateScriptThemeBootstrap)
	if themeScript == -1 {
		t.Fatal("the theme bootstrap script must be inside <head>")
	}
	if style := strings.Index(head, "<style>"); style == -1 || themeScript > style {
		t.Fatal("the theme bootstrap script must run before the stylesheet")
	}
	if !strings.Contains(templateScriptThemeBootstrap, "\n        syncTheme();\n    </script>") {
		t.Fatal("the theme bootstrap script must apply the theme immediately, not on a later event")
	}
}

func TestThemeBootstrapAppliesCPAMCThemeSynchronously(t *testing.T) {
	node, err := exec.LookPath("node")
	if err != nil {
		t.Skip("node is not available; theme bootstrap behavior is checked in the Node-enabled verification environment")
	}
	script := strings.TrimSuffix(strings.TrimPrefix(strings.TrimSpace(templateScriptThemeBootstrap), "<script>"), "</script>")
	cases := []struct {
		name, setup, want string
	}{
		{"stored dark theme", `store["cli-proxy-theme"] = JSON.stringify({state:{theme:"dark"}});`, "dark"},
		{"stored white theme", `store["cli-proxy-theme"] = JSON.stringify({state:{theme:"white"}});`, "white"},
		{"stored auto follows dark system", `store["cli-proxy-theme"] = JSON.stringify({state:{theme:"auto"}}); systemDark = true;`, "dark"},
		{"no stored theme follows dark system", `systemDark = true;`, "dark"},
		{"stored light theme", `store["cli-proxy-theme"] = JSON.stringify({state:{theme:"light"}});`, "null"},
		{"dark CPAMC parent wins over storage", `store["cli-proxy-theme"] = JSON.stringify({state:{theme:"light"}}); parentTheme = "dark";`, "dark"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			harness := `
const store = {};
let systemDark = false;
let parentTheme = undefined;
const attributes = {};
const classes = {};
const documentElement = {
    getAttribute: function(name) { return name in attributes ? attributes[name] : null; },
    setAttribute: function(name, value) { attributes[name] = String(value); },
    removeAttribute: function(name) { delete attributes[name]; },
    classList: { toggle: function(name, on) { classes[name] = Boolean(on); } }
};
const document = { documentElement: documentElement };
const localStorage = { getItem: function(name) { return name in store ? store[name] : null; } };
const window = { matchMedia: function() { return { matches: systemDark }; } };
` + tc.setup + `
if (parentTheme !== undefined) {
    window.self = window;
    window.top = {};
    window.parent = { document: { documentElement: { getAttribute: function() { return parentTheme; } } } };
} else {
    window.self = window;
    window.top = window;
}
` + script + `
const got = documentElement.getAttribute("data-theme");
if (String(got) !== ` + strconv.Quote(tc.want) + `) {
    console.error("data-theme = " + got + ", want ` + tc.want + `");
    process.exit(1);
}
if (classes["is-embedded"] !== (parentTheme !== undefined)) {
    console.error("is-embedded class not set before first paint");
    process.exit(1);
}
`
			runNodeFixture(t, node, "theme-bootstrap.js", harness)
		})
	}
}

func TestStatusHTML_CPAThemeAlignmentAndDevServerHelper(t *testing.T) {
	if strings.Contains(StatusHTML, "id=\"btnThemeToggle\"") || strings.Contains(StatusHTML, "toggleTheme()") {
		t.Fatal("the page follows CPAMC's theme and must not expose its own theme toggle")
	}

	script := extractPageScript(t)
	for _, snippet := range []string{
		"window.__setTheme = window.setTheme = function",
		`const CPAMC_THEME_KEY = "cli-proxy-theme";`,
		`new MutationObserver(syncTheme).observe(window.parent.document.documentElement, { attributes: true, attributeFilter: ["data-theme"] });`,
		`if (event.key === CPAMC_THEME_KEY) syncTheme();`,
	} {
		if !strings.Contains(script, snippet) {
			t.Errorf("theme handling must mirror key-model-access (missing %q)", snippet)
		}
	}
	for _, selector := range []string{`:root[data-theme="white"] {`, `:root[data-theme="dark"] {`} {
		if !strings.Contains(StatusHTML, selector) {
			t.Errorf("stylesheet is missing the CPAMC theme palette %s", selector)
		}
	}
}

func TestHistoryComparisonStacksOnNarrowScreens(t *testing.T) {
	for _, snippet := range []string{
		`data-label=\"Gemini quota\"`,
		`data-label=\"Claude/GPT quota\"`,
		`data-label=\"Priority change\"`,
		`.history-comparison-header { display: none; }`,
		`.history-cell::before { content: attr(data-label);`,
	} {
		if !strings.Contains(StatusHTML, snippet) {
			t.Errorf("narrow history comparison layout is missing %q", snippet)
		}
	}
}
func TestProbeOutcomeToastAndRenderFailureTruth(t *testing.T) {
	node, err := exec.LookPath("node")
	if err != nil {
		t.Skip("node is not available; probe outcome behavior is checked in the Node-enabled verification environment")
	}

	harness := templateScriptOverviewActionsCore + templateScriptOverviewRender + templateScriptProbeSchedule + `
const toastCalls = [];
function element(initial) {
    return Object.assign({
        hidden: false,
        disabled: false,
        textContent: "",
        innerHTML: "",
        children: [],
        attributes: {},
        dataset: {},
        style: {},
        classList: {
            classes: [],
            add: function(c) { this.classes.push(c); },
            remove: function(c) { this.classes = this.classes.filter(x => x !== c); },
            contains: function(c) { return this.classes.includes(c); },
            toggle: function(c, on) { if (on) this.add(c); else this.remove(c); }
        },
        appendChild: function(child) { this.children.push(child); },
        setAttribute: function(k, v) { this.attributes[k] = v; },
        getAttribute: function(k) { return this.attributes[k] || null; },
        querySelector: function() { return null; },
        querySelectorAll: function() { return []; },
        addEventListener: function() {}
    }, initial || {});
}

const elements = {
    modelGroupSelect: element({value: "gemini"}),
    valLastAudit: element(),
    valNextProbe: element(),
    valTotalCreds: element(),
    valTotalDesc: element(),
    valBoosted: element(),
    valDepleted: element(),
    credentialsContainer: element(),
    btnProbe: element()
};

const document = {
    getElementById: function(id) { return elements[id] || null; },
    createElement: function(tag) { return element({tagName: tag}); },
    querySelector: function() { return null; },
    querySelectorAll: function() { return []; }
};

function showToast(msg, type) { toastCalls.push({msg: msg, type: type}); }
function escapeHTML(s) { return String(s == null ? "" : s).replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;").replace(/"/g, "&quot;").replace(/'/g, "&#39;"); }
function formatCountdown() { return "01h 00m"; }
const ICONS = new Proxy({}, { get: function() { return ""; } });

dynamicConfig = { auto_apply: false, antigravity_model_group: "gemini" };
latestDiagnostics = { latest_audit: "test audit" };
scheduleConfig = null;
isAuthBlocked = false;
userSelectedModelGroup = false;
const RUN_PATH = "/run";
const SNAPSHOT_PATH = "/snapshot/latest";
const DIAGNOSTICS_PATH = "/diagnostics";

async function apiFetch(path) {
    if (path.indexOf("mode=probe") >= 0) {
        return {
            attempted: 0,
            succeeded: 0,
            probe_outcomes: [
                {
                    auth_index: "fail-403",
                    model_group: "gemini",
                    status: "failed",
                    error: "retrieve quota summary status 403",
                    http_status_code: 403
                }
            ]
        };
    }
    return {};
}

async function fetchSnapshot() {}
async function fetchDiagnostics() {}
async function refreshDashboard() {}
function isAutoApplyEnabled() { return false; }

(async function() {
    // 1. A failed single-credential probe must surface as an error with its HTTP status.
    const label = element({textContent: "Probe"});
    const probeBtn = element({querySelector: function() { return label; }});
    await probeSingleCredential("fail-403", probeBtn);
    if (toastCalls.length !== 1) throw new Error("expected exactly 1 toast, got " + toastCalls.length);
    if (toastCalls[0].type !== "error") throw new Error("expected error toast for failed probe, got: " + toastCalls[0].type);
    if (toastCalls[0].msg.indexOf("403") < 0 || toastCalls[0].msg.indexOf("Probe failed") < 0) {
        throw new Error("expected toast to mention Probe failed and 403, got: " + toastCalls[0].msg);
    }

    // 2. Card rendering for a failed probe and for a 429 cooldown.
    latestSnapshot = {
        active_model_group: "gemini",
        groups: {
            gemini: {
                items: [
                    {
                        auth_index: "fail-403",
                        email: "test-403@example.com",
                        reason: "probe failed: retrieve quota summary status 403; historical: fresh remaining positive",
                        evidence_fresh: false,
                        current: { priority: 99, disabled: false },
                        target: { priority: 99, disabled: false },
                        r5h: 1.0,
                        r7d: 1.0,
                        urgency: 0.98,
                        cycle_burn_rate: 0.16
                    },
                    {
                        auth_index: "cooldown-429",
                        email: "test-429@example.com",
                        reason: "probe failed: retrieve quota summary status 429; historical: 429 rate limit cooldown",
                        evidence_fresh: false,
                        current: { priority: -1, disabled: false },
                        target: { priority: -1, disabled: false },
                        r5h: 0.8,
                        r7d: 0.9,
                        urgency: 1.1,
                        cycle_burn_rate: 0.15
                    }
                ]
            }
        }
    };

    renderDashboard();

    const renderedCards = elements.credentialsContainer.children;
    if (renderedCards.length !== 2) throw new Error("expected 2 rendered cards, got " + renderedCards.length);

    const card0HTML = renderedCards[0].innerHTML;
    if (card0HTML.indexOf("badge-danger\">Probe failed") < 0) throw new Error("card 0 should have a badge-danger 'Probe failed': " + card0HTML);
    if (card0HTML.indexOf("Positive balance") >= 0) throw new Error("card 0 must NOT show 'Positive balance' when the probe failed");
    if (card0HTML.indexOf("Probe failed: retrieve quota summary status 403") < 0) throw new Error("card 0 should show the probe failure detail: " + card0HTML);
    if (card0HTML.indexOf("meter-stale") < 0 || card0HTML.indexOf("Historical") < 0) throw new Error("card 0 should mark non-fresh meters as Historical: " + card0HTML);

    const card1HTML = renderedCards[1].innerHTML;
    if (card1HTML.indexOf("badge-warning\">429 cooldown") < 0) throw new Error("card 1 should prioritise the '429 cooldown' badge over probe failed: " + card1HTML);
    if (card1HTML.indexOf("Probe failed: retrieve quota summary status 429") < 0) throw new Error("card 1 reason line should show the failure cause: " + card1HTML);

    // 3. A batch probe with failures reports how many failed.
    await triggerProbe();
    if (typeof probeCooldownTimer !== "undefined") clearInterval(probeCooldownTimer);
    const lastToast = toastCalls[toastCalls.length - 1];
    if (lastToast.type !== "error" || lastToast.msg.indexOf("Quota probe finished with 1 failed credential") < 0) {
        throw new Error("expected batch probe failure toast, got: " + JSON.stringify(lastToast));
    }
})().catch(function(err) {
    console.error(err.stack || err.message);
    process.exit(1);
});
`

	runNodeFixture(t, node, "probe-outcome-truth.js", harness)
}
