package management

// shellModalMarkup contains modal foundations shared by multiple features.
const shellModalMarkup = `
    <!-- Diff & Run Result Modal -->
    <div id="diffModal" class="modal-backdrop" hidden onclick="closeModalFromBackdrop(event)">
        <div id="diffModalPanel" class="modal" role="dialog" aria-modal="true" aria-labelledby="modalTitle">
            <div class="modal-head">
                <div>
                    <h2 id="modalTitle">Schedule change preview</h2>
                    <p id="modalSummary"></p>
                </div>
                <button type="button" class="icon-button" onclick="closeModal()" aria-label="Close">` + iconClose + `</button>
            </div>
            <div id="modalDiffList" class="modal-body"></div>
            <div class="modal-foot">
                <button type="button" class="button secondary compact" onclick="closeModal()">Close</button>
                <button type="button" class="button primary compact" id="btnModalApply" onclick="executeDirectApply()">Confirm apply</button>
            </div>
        </div>
    </div>

    <!-- Themed Confirm Modal -->
    <div id="confirmModal" class="modal-backdrop" hidden>
        <div class="modal narrow" role="dialog" aria-modal="true" aria-labelledby="confirmTitle">
            <div class="confirm-body">
                <h2 class="confirm-title" id="confirmTitle"></h2>
                <p class="confirm-message" id="confirmMessage"></p>
            </div>
            <div class="modal-foot">
                <button type="button" class="button secondary compact" id="confirmCancelBtn">Cancel</button>
                <button type="button" class="button danger-solid compact" id="confirmOkBtn">Confirm</button>
            </div>
        </div>
    </div>

    <div id="toastRoot" class="toast-region" aria-live="polite" aria-atomic="true"></div>
`
