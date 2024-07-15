let socket = new WebSocket("ws://localhost:8080/ws");

socket.onopen = function () {
    console.log("Connection established");
};

if (navigator.geolocation) {
    navigator.geolocation.watchPosition(
        (position) => {
            const { latitude, longitude } = position.coords;

            const data = {
                latitude: latitude,
                longitude: longitude
            };

            socket.send(JSON.stringify(data));
        },
        (error) => {
            console.error('Error watching position: ', error);
        },
        {
            enableHighAccuracy: true,
            timeout: 5000,
            maximumAge: 0
        });
} else {
    console.error('Geolocation is not supported by this browser.');
}

socket.onclose = function (event) {
    if (event.wasClean) {
        document.getElementById('status').innerText = `Connection closed cleanly, code=${event.code} reason=${event.reason}`;
        console.log(`Connection closed cleanly, code=${event.code} reason=${event.reason}`);
    } else {
        document.getElementById('status').innerText = "Connection died";
        console.log('Connection died');
    }
};

socket.onerror = function (error) {
    document.getElementById('status').innerText = `WebSocket error: ${error.message}`;
    console.log(`WebSocket error: ${error.message}`);
};


const map = L.map("map").setView([0, 0], 16);

L.tileLayer("https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png", {
    attribution: "Tridip Dam"
}).addTo(map);

const markers = {}

socket.onmessage = (e) => {
    const data = JSON.parse(e.data);
    if (data.clientId) {
        if (markers[data.clientId]) {
            markers[data.clientId].setLatLng([data.latitude, data.longitude]);
            console.log(`Updated marker for client ${data.clientId}:`, data.latitude, data.longitude);
        } else {
            markers[data.clientId] = L.marker([data.latitude, data.longitude]).addTo(map).bindPopup(`Client ${data.clientId}`);
            console.log(`Added marker for client ${data.clientId}:`, data.latitude, data.longitude);
        }

        const markerBounds = L.latLngBounds(Object.values(markers).map(marker => marker.getLatLng()));
        map.fitBounds(markerBounds, { padding: [50, 50] });
    } else {
        console.error('Received data without clientId:', data);
    }
};




