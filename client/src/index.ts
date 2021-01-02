fetch('http://localhost:6482/proto/lovelove.proto')
    .then(response => response.text())
    .then(data => console.log(data));

const webSocket = new WebSocket("ws://localhost:6482/echo");
webSocket.onopen = function (event) {
    console.log(event);
    webSocket.send("test");
}

webSocket.onclose = function (event) {
    console.log(event);
}

webSocket.onerror = function (event) {
    console.log(event);
}

webSocket.onmessage = function (event) {
    console.log(event);
}

