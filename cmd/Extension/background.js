// Dropbox OAuth configuration
const DROPBOX_CLIENT_ID = 'your_dropbox_app_key'; // Replace with your Dropbox app key
const REDIRECT_URI = chrome.identity.getRedirectURL('oauth');
const YOUR_API_ENDPOINT = 'https://your-api-domain.com/api/submit'; // Replace with your API endpoint

class DropboxAuth {
    constructor() {
        this.authUrl = `https://www.dropbox.com/oauth2/authorize?client_id=${DROPBOX_CLIENT_ID}&response_type=code&redirect_uri=${encodeURIComponent(REDIRECT_URI)}`;
    }

    async authenticate() {
        try {
            const responseUrl = await chrome.identity.launchWebAuthFlow({
                url: this.authUrl,
                interactive: true
            });

            const code = this.extractCodeFromUrl(responseUrl);
            if (!code) {
                throw new Error('Authorization code not found');
            }

            const tokens = await this.exchangeCodeForTokens(code);
            await this.storeTokens(tokens);
            return tokens;
        } catch (error) {
            console.error('Authentication failed:', error);
            throw error;
        }
    }

    extractCodeFromUrl(url) {
        const match = url.match(/code=([^&]+)/);
        return match ? match[1] : null;
    }

    async exchangeCodeForTokens(code) {
        const response = await fetch('https://api.dropboxapi.com/oauth2/token', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded',
            },
            body: new URLSearchParams({
                code: code,
                grant_type: 'authorization_code',
                client_id: DROPBOX_CLIENT_ID,
                redirect_uri: REDIRECT_URI
            })
        });

        if (!response.ok) {
            throw new Error(`Token exchange failed: ${response.statusText}`);
        }

        const data = await response.json();
        return {
            access_token: data.access_token,
            refresh_token: data.refresh_token,
            expires_at: Date.now() + (data.expires_in * 1000)
        };
    }

    async storeTokens(tokens) {
        await chrome.storage.local.set({
            dropbox_tokens: tokens
        });
    }

    async getStoredTokens() {
        const result = await chrome.storage.local.get('dropbox_tokens');
        return result.dropbox_tokens || null;
    }

    async refreshAccessToken(refreshToken) {
        const response = await fetch('https://api.dropboxapi.com/oauth2/token', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded',
            },
            body: new URLSearchParams({
                grant_type: 'refresh_token',
                refresh_token: refreshToken,
                client_id: DROPBOX_CLIENT_ID
            })
        });

        if (!response.ok) {
            throw new Error(`Token refresh failed: ${response.statusText}`);
        }

        const data = await response.json();
        const newTokens = {
            access_token: data.access_token,
            refresh_token: refreshToken, // Refresh token doesn't change
            expires_at: Date.now() + (data.expires_in * 1000)
        };

        await this.storeTokens(newTokens);
        return newTokens;
    }

    async getValidAccessToken() {
        let tokens = await this.getStoredTokens();

        if (!tokens) {
            // No tokens stored, need to authenticate
            tokens = await this.authenticate();
        } else if (Date.now() >= tokens.expires_at - 60000) { // Refresh if expires in less than 1 minute
            // Token expired or about to expire, refresh it
            tokens = await this.refreshAccessToken(tokens.refresh_token);
        }

        return tokens.access_token;
    }
}

const dropboxAuth = new DropboxAuth();

// Handle messages from popup
chrome.runtime.onMessage.addListener((message, sender, sendResponse) => {
    if (message.action === 'sendHtmlToApi') {
        handleSendHtml(message.tabId)
            .then(result => sendResponse({ success: true, result }))
            .catch(error => sendResponse({ success: false, error: error.message }));
        return true; // Keep message channel open for async response
    }

    if (message.action === 'authenticate') {
        dropboxAuth.authenticate()
            .then(tokens => sendResponse({ success: true, tokens }))
            .catch(error => sendResponse({ success: false, error: error.message }));
        return true;
    }

    if (message.action === 'checkAuth') {
        dropboxAuth.getStoredTokens()
            .then(tokens => sendResponse({ success: true, authenticated: !!tokens }))
            .catch(error => sendResponse({ success: false, error: error.message }));
        return true;
    }
});

async function handleSendHtml(tabId) {
    try {
        // Get valid access token
        const accessToken = await dropboxAuth.getValidAccessToken();

        // Get HTML content from the active tab
        const results = await chrome.scripting.executeScript({
            target: { tabId: tabId },
            function: () => document.documentElement.outerHTML
        });

        const htmlContent = results[0].result;

        // Get tab info for additional context
        const tab = await chrome.tabs.get(tabId);

        // Send to your API
        const response = await fetch(YOUR_API_ENDPOINT, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${accessToken}`
            },
            body: JSON.stringify({
                html: htmlContent,
                url: tab.url,
                title: tab.title,
                timestamp: new Date().toISOString()
            })
        });

        if (!response.ok) {
            throw new Error(`API request failed: ${response.statusText}`);
        }

        const result = await response.json();
        return result;

    } catch (error) {
        console.error('Failed to send HTML to API:', error);
        throw error;
    }
}