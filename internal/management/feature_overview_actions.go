package management

const templateScriptOverviewActionsCore = `        let latestSnapshot = null;
        let userSelectedModelGroup = false;

        function formatReason(reason, isBoosted, isDisabled) {
            if (!reason) {
                if (isDisabled) return "Weekly depleted";
                if (isBoosted) return "Boosted";
                return "Active";
            }
            var lower = reason.toLowerCase();
            if (lower.indexOf("boost") >= 0) return "Boosted";
            if (lower.indexOf("remaining positive") >= 0) return "Positive balance";
            if (lower.indexOf("weekly depleted") >= 0) return "Weekly depleted";
            if (lower.indexOf("short") >= 0 && lower.indexOf("depleted") >= 0) return "5h depleted";
            if (lower.indexOf("429") >= 0 || lower.indexOf("cooldown") >= 0) return "429 cooldown";
            if (lower.indexOf("predicted") >= 0) return "Predicted";
            if (lower.indexOf("in sync") >= 0 || lower.indexOf("optimal") >= 0) return "In sync";
            if (lower.indexOf("disabled on host") >= 0) return "Disabled";
            return reason;
        }

        async function fetchSnapshot(options) {
            if (isAuthBlocked) return;
            const settings = options || {};
            try {
                const data = await apiFetch(SNAPSHOT_PATH);
                latestSnapshot = data;
                syncModelGroupSelect(data && data.active_model_group);
                renderDashboard();
            } catch (err) {
                if (!settings.silent && !isAuthBlocked) showToast(err.message, "error");
            }
        }

        async function fetchDiagnostics(options) {
            if (isAuthBlocked) return;
            const settings = options || {};
            try {
                const data = await apiFetch(DIAGNOSTICS_PATH);
                latestDiagnostics = data;
                renderDashboard();
                renderHistory();
                renderDiagnostics();
                renderScheduleStatus();
            } catch (err) {
                if (!settings.silent && !isAuthBlocked) showToast(err.message, "error");
            }
        }

        async function syncHost(options) {
            if (isAuthBlocked) return;
            const settings = options || {};
            try {
				const data = await apiFetch(SYNC_PATH, { method: "POST" });
                latestSnapshot = data;
                renderDashboard();
                if (settings.notifySync) showToast("Synchronized with CPA host auth files", "success");
            } catch (err) {
                if (isAuthBlocked) throw err;
                await fetchSnapshot(settings);
            }
        }

        async function refreshDashboard(options) {
            if (isAuthBlocked) return;
            const settings = options || {};
            const btn = document.getElementById("btnRefresh");
            if (btn && settings.notifySync) btn.disabled = true;
            try {
                if (settings.withSync) {
                    await syncHost(settings);
                } else {
                    await fetchSnapshot(settings);
                }
                if (!isAuthBlocked) {
                    await fetchDiagnostics(settings);
                }
            } finally {
                if (btn && settings.notifySync) btn.disabled = false;
            }
        }

        async function probeSingleCredential(authIndex, btn) {
            if (isAuthBlocked || !authIndex) return;
            const label = btn ? btn.querySelector("span") : null;
            const originalText = label ? label.textContent : "";
            if (btn) btn.disabled = true;
            if (label) label.textContent = "Probing…";
            try {
                const groupSelect = document.getElementById("modelGroupSelect");
                const modelGroup = groupSelect ? groupSelect.value : "gemini";
                const url = RUN_PATH + "?mode=probe&antigravity_model_group=" + encodeURIComponent(modelGroup) + "&auth_index=" + encodeURIComponent(authIndex);
                await apiFetch(url, { method: "POST" });
                await fetchSnapshot({ silent: true });
                showToast("Credential quota updated", "success");
            } catch (err) {
                showToast(err.message, "error");
            } finally {
                if (btn) btn.disabled = false;
                if (label) label.textContent = originalText || "Probe";
            }
        }

        function extractChanges(result) {
            if (!result) return [];
            if (Array.isArray(result.changes) && result.changes.length > 0) {
                return result.changes;
            }
            if (result.snapshot && Array.isArray(result.snapshot.changes) && result.snapshot.changes.length > 0) {
                return result.snapshot.changes;
            }
            if (result.snapshot && Array.isArray(result.snapshot.items) && result.snapshot.items.length > 0) {
                return result.snapshot.items;
            }
            if (Array.isArray(result.items) && result.items.length > 0) {
                return result.items;
            }
            return [];
        }

        function hasProjectedStateChange(item) {
            if (!item || !item.current || !item.target) return false;

            const currentDisabled = Boolean(item.current.disabled);
            const targetDisabled = Boolean(item.target.disabled);
            if (currentDisabled !== targetDisabled) return true;
            if (currentDisabled && targetDisabled) return false;

            const currentMissing = Boolean(item.current.priority_missing) || item.current.priority === undefined || item.current.priority === null;
            const targetMissing = Boolean(item.target.priority_missing) || item.target.priority === undefined || item.target.priority === null;
            if (currentMissing !== targetMissing) return true;
            if (currentMissing && targetMissing) return false;

            return String(item.current.priority) !== String(item.target.priority);
        }

        function buildApplyPreview(groupData) {
            const pendingChanges = Array.isArray(groupData.changes) ? groupData.changes : [];
            if (pendingChanges.length > 0) {
                return { mode: "pending", changes: pendingChanges };
            }

            const projectedChanges = (Array.isArray(groupData.items) ? groupData.items : []).filter(hasProjectedStateChange);
            return { mode: "projected", changes: projectedChanges };
        }

        async function triggerApplyWithConfirm() {
            const groupSelect = document.getElementById("modelGroupSelect");
            const selectedGroup = groupSelect ? groupSelect.value : "gemini";
            const activeGroup = (dynamicConfig && dynamicConfig.antigravity_model_group) || (latestSnapshot && latestSnapshot.active_model_group) || "gemini";
            if (selectedGroup !== activeGroup) {
                showToast("Only the active model group can apply. Active: " + activeGroup, "info");
                return;
            }
            const groupData = (latestSnapshot && latestSnapshot.groups && latestSnapshot.groups[selectedGroup]) || {};
            const items = groupData.items || [];
            if (items.length === 0) {
                showToast("No credentials found to apply", "info");
                return;
            }
            const preview = buildApplyPreview(groupData);
            if (preview.changes.length === 0) {
                showToast("All credentials are in their optimal state; nothing to apply", "info");
                return;
            }

            // Show changes preview modal with direct confirm button
            showModal("apply-confirm", {
                snapshot: groupData,
                changes: preview.changes,
                preview_mode: preview.mode,
                items: items
            }, false);
        }

        async function executeDirectApply() {
            closeModal();
            const groupSelect = document.getElementById("modelGroupSelect");
            const group = groupSelect ? groupSelect.value : "gemini";
            const btnApp = document.getElementById("btnApply");
            if (btnApp) btnApp.disabled = true;

            try {
                const path = RUN_PATH + "?mode=apply&antigravity_model_group=" + encodeURIComponent(group);
                const previewID = latestSnapshot && latestSnapshot.preview_id ? latestSnapshot.preview_id : "";
                const applyPath = previewID ? path + "&preview_id=" + encodeURIComponent(previewID) : path;
                const result = await apiFetch(applyPath, { method: "POST" });
                const succeeded = (result && result.succeeded !== undefined) ? result.succeeded : 0;
                const attempted = (result && result.attempted !== undefined) ? result.attempted : 0;
                if (succeeded > 0 || attempted > 0) {
                    showToast("Applied credential priorities to the CPA host", "success");
                } else {
                    showToast("All credentials are in their optimal state; nothing to apply", "info");
                }
                await refreshDashboard({ silent: true });
            } catch (err) {
                showToast(err.message, "error");
            } finally {
                if (btnApp) btnApp.disabled = false;
            }
        }

`
