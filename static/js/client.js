var Client = (function () {

    "use strict";

    return {
        connect: function (onConnect, onMessage) {
            if (!window.WebSocket) {
                alert("Your browser does not support WebSockets. :|");
                return;
            }

            // Let's open a web socket to the server
            var proto = "ws";
            if (window.location.protocol == 'https:') {
                proto = "wss";
            }

            var conn = new window.WebSocket(proto + "://" + document.location.host + "/ws/");
            conn.binaryType = "arraybuffer";

            conn.onopen = function (data) {
                console.log("connection was opened to '" + data.currentTarget.url + '"');
                onConnect();
            };
            conn.onerror = function () {
                console.log("connection error");
            };
            conn.onmessage = function (evt) {
                onMessage(new DataStream(evt.data));
            };
            conn.onclose = function (event) {
                console.log("connection was closed to '" + event.currentTarget.url + '"');
            };
        }
    }
})();
