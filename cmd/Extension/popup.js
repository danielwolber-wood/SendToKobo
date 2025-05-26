document.addEventListener('DOMContentLoaded', async () => {
    const authStatusEl = document.getElementById('authStatus');
    const authSectionEl = document.getElementById('authSection');
    const mainSectionEl = document.getElementById('mainSection');
    const statusEl = document.getElementById('status');
    const authenticateBtn = document.getElementById('authenticateBtn');
    const sendHtmlBtn = document.getElementById('sendHtmlBtn');
    const reauthBtn = document.getElementById('reauthBtn');

    // Check authentication status on load
    await checkAuthStatus();

    // Event listeners
    authenticateBtn.addEventListener('click', handleAuthenticate);
    sendHtmlBtn.addEventListener('click', handleSendHtml);
    reauthBtn.addEventListener('click', handleAuthenticate);

    async function checkAuthStatus() {
        try {
            const response = await sendMessage({ action: 'checkAuth' });

            if (response.success && response.authenticated) {
                showAuthenticatedState();
            } else {
                showUnauthenticatedState();
            }
        } catch (error) {
            showError('Failed to check authentication status');
            showUnauthenticatedState();
        }
    }

    async function handleAuthenticate() {
        showLoading(authenticateBtn, 'Connecting...');
        hideStatus();

        try {
            const response = await sendMessage({ action: 'authenticate' });

            if (response.success) {
                showSuccess('Successfully connected to Dropbox!');
                showAuthenticatedState();
            } else {
                showError(`Authentication failed: ${response.error}`);
            }
        } catch (error) {
            showError('Authentication failed');
        } finally {
            hideLoading(authenticateBtn, 'Connect to Dropbox');
        }
    }

    async function handleSendHtml() {
        showLoading(sendHtmlBtn, 'Sending...');
        hideStatus();

        try {
            // Get current active tab
            const [tab] = await chrome.tabs.query({ active: true, currentWindow: true });

            const response = await sendMessage({
                action: 'sendHtmlToApi',
                tabId: tab.id
            });

            if (response.success) {
                showSuccess('HTML sent successfully!');
            } else {
                showError(`Failed to send HTML: ${response.error}`);
            }
        } catch (error) {
            showError('Failed to send HTML');
        } finally {
            hideLoading(sendHtmlBtn, 'Send Current Page');
        }
    }

    function showAuthenticatedState() {
        authStatusEl.textContent = '✓ Connected to Dropbox';
        authStatusEl.className = 'auth-status authenticated';
        authSectionEl.style.display = 'none';
        mainSectionEl.style.display = 'block';
    }

    function showUnauthenticatedState() {
        authStatusEl.textContent = '⚠ Not connected to Dropbox';
        authStatusEl.className = 'auth-status';
        authSectionEl.style.display = 'block';
        mainSectionEl.style.display = 'none';
    }

    function showStatus(message, type) {
        statusEl.textContent = message;
        statusEl.className = `status ${type}`;
        statusEl.style.display = 'block';
    }

    function showSuccess(message) {
        showStatus(message, 'success');
    }

    function showError(message) {
        showStatus(message, 'error');
    }

    function showInfo(message) {
        showStatus(message, 'info');
    }

    function hideStatus() {
        statusEl.style.display = 'none';
    }

    function showLoading(button, loadingText) {
        button.textContent = loadingText;
        button.classList.add('loading');
    }

    function hideLoading(button, originalText) {
        button.textContent = originalText;
        button.classList.remove('loading');
    }

    function sendMessage(message) {
        return new Promise((resolve) => {
            chrome.runtime.sendMessage(message, resolve);
        });
    }
});