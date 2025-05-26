// Content script - runs on every webpage
// This script can be used to enhance HTML extraction if needed

// Listen for messages from background script if additional processing is needed
chrome.runtime.onMessage.addListener((message, sender, sendResponse) => {
    if (message.action === 'getEnhancedHtml') {
        // You can add custom HTML processing logic here if needed
        // For example, clean up the HTML, remove certain elements, etc.

        const cleanHtml = document.documentElement.outerHTML;
        sendResponse({ html: cleanHtml });
    }
});

// Optional: Add visual feedback when HTML is being captured
function showCaptureIndicator() {
    const indicator = document.createElement('div');
    indicator.id = 'html-capture-indicator';
    indicator.style.cssText = `
    position: fixed;
    top: 20px;
    right: 20px;
    background: #0061ff;
    color: white;
    padding: 10px 15px;
    border-radius: 6px;
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
    font-size: 14px;
    z-index: 10000;
    box-shadow: 0 4px 12px rgba(0,0,0,0.2);
    animation: slideIn 0.3s ease-out;
  `;

    // Add CSS animation
    if (!document.getElementById('capture-indicator-style')) {
        const style = document.createElement('style');
        style.id = 'capture-indicator-style';
        style.textContent = `
      @keyframes slideIn {
        from { transform: translateX(100%); opacity: 0; }
        to { transform: translateX(0); opacity: 1; }
      }
    `;
        document.head.appendChild(style);
    }

    indicator.textContent = 'Capturing page...';
    document.body.appendChild(indicator);

    // Remove indicator after 2 seconds
    setTimeout(() => {
        if (indicator.parentNode) {
            indicator.style.animation = 'slideIn 0.3s ease-out reverse';
            setTimeout(() => {
                if (indicator.parentNode) {
                    indicator.remove();
                }
            }, 300);
        }
    }, 2000);
}

// You can call this function when capture starts if you want visual feedback
// showCaptureIndicator();