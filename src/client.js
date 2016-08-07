const DataStream = require('./datastream.js');

let conn;

module.exports = {
  connect(onMessage) {
    if (!window.WebSocket) {
      alert('Your browser does not support WebSockets. :|'); // eslint-disable-line
      return;
    }

    // Let's open a web socket to the server
    let proto = 'ws';
    if (window.location.protocol === 'https:') {
      proto = 'wss';
    }

    conn = new window.WebSocket(`${proto}://${document.location.host}/ws/`);
    conn.binaryType = 'arraybuffer';

    conn.onopen = function socketOpen(data) {
      console.log(`connection was opened to "${data.currentTarget.url}"`); // eslint-disable-line
    };

    conn.onerror = function socketOnError() {
      console.log('connection error'); // eslint-disable-line no-console
    };

    conn.onmessage = function socketOnMessage(evt) {
      onMessage(new DataStream(evt.data));
    };

    conn.onclose = function socketOnClose(event) {
      console.log('connection was closed to "' + event.currentTarget.url + '"'); // eslint-disable-line
    };
  },

  send(message) {
    conn.send(message);
  },
};
