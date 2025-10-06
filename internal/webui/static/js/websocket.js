// WebSocket connection manager for real-time updates

class WebSocketManager {
    constructor(url) {
        this.url = url;
        this.ws = null;
        this.reconnectInterval = 3000;
        this.reconnectTimer = null;
        this.listeners = {};
        this.isConnected = false;
    }

    connect() {
        try {
            const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
            const wsUrl = `${protocol}//${window.location.host}${this.url}`;

            this.ws = new WebSocket(wsUrl);

            this.ws.onopen = () => {
                console.log('WebSocket connected');
                this.isConnected = true;
                this.updateStatus('connected');
                if (this.reconnectTimer) {
                    clearTimeout(this.reconnectTimer);
                    this.reconnectTimer = null;
                }
            };

            this.ws.onclose = () => {
                console.log('WebSocket disconnected');
                this.isConnected = false;
                this.updateStatus('disconnected');
                this.reconnect();
            };

            this.ws.onerror = (error) => {
                console.error('WebSocket error:', error);
                this.updateStatus('error');
            };

            this.ws.onmessage = (event) => {
                try {
                    const message = JSON.parse(event.data);
                    this.handleMessage(message);
                } catch (error) {
                    console.error('Error parsing WebSocket message:', error);
                }
            };
        } catch (error) {
            console.error('Error connecting to WebSocket:', error);
            this.reconnect();
        }
    }

    reconnect() {
        if (this.reconnectTimer) return;

        this.reconnectTimer = setTimeout(() => {
            console.log('Reconnecting to WebSocket...');
            this.connect();
        }, this.reconnectInterval);
    }

    disconnect() {
        if (this.reconnectTimer) {
            clearTimeout(this.reconnectTimer);
            this.reconnectTimer = null;
        }
        if (this.ws) {
            this.ws.close();
            this.ws = null;
        }
    }

    handleMessage(message) {
        const { type, data } = message;

        // Call all registered listeners for this message type
        if (this.listeners[type]) {
            this.listeners[type].forEach(callback => {
                try {
                    callback(data);
                } catch (error) {
                    console.error(`Error in listener for ${type}:`, error);
                }
            });
        }

        // Call global listeners
        if (this.listeners['*']) {
            this.listeners['*'].forEach(callback => {
                try {
                    callback(message);
                } catch (error) {
                    console.error('Error in global listener:', error);
                }
            });
        }
    }

    on(type, callback) {
        if (!this.listeners[type]) {
            this.listeners[type] = [];
        }
        this.listeners[type].push(callback);
    }

    off(type, callback) {
        if (!this.listeners[type]) return;

        const index = this.listeners[type].indexOf(callback);
        if (index > -1) {
            this.listeners[type].splice(index, 1);
        }
    }

    updateStatus(status) {
        const statusIndicator = document.getElementById('ws-status');
        const statusText = document.getElementById('ws-status-text');

        if (!statusIndicator || !statusText) return;

        switch (status) {
            case 'connected':
                statusIndicator.style.color = '#198754';
                statusText.textContent = 'Connected';
                break;
            case 'disconnected':
                statusIndicator.style.color = '#6c757d';
                statusText.textContent = 'Disconnected';
                break;
            case 'error':
                statusIndicator.style.color = '#dc3545';
                statusText.textContent = 'Error';
                break;
            default:
                statusIndicator.style.color = '#ffc107';
                statusText.textContent = 'Connecting...';
        }
    }

    send(message) {
        if (this.ws && this.isConnected) {
            this.ws.send(JSON.stringify(message));
        } else {
            console.error('WebSocket is not connected');
        }
    }
}

// Create global WebSocket manager instance
const wsManager = new WebSocketManager('/api/v1/ws');

// Auto-connect on page load
window.addEventListener('load', () => {
    wsManager.connect();
});

// Clean up on page unload
window.addEventListener('beforeunload', () => {
    wsManager.disconnect();
});
