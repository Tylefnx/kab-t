let socket;

export function setupWebSocket(url, onMessageCallback) {
    socket = new WebSocket(url);

    socket.onopen = () => {
        console.log('Connected to the server');
    };

    socket.onmessage = (event) => {
        const data = JSON.parse(event.data);
        console.log('Received data:', data);
        onMessageCallback(data);
    };

    socket.onclose = () => {
        console.log('Disconnected from the server');
    };
}

export function sendWebSocketMessage(message) {
    if (socket && socket.readyState === WebSocket.OPEN) {
        socket.send(JSON.stringify(message));
    }
}
