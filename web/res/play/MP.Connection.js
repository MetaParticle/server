MP.Connection = function(ip, port) {
	this.ip = ip || "localhost";
	this.port = port || "8080";
	return new WebSocket("ws://" + this.ip + ":" + this.port + "/ws");
}

MP.Connection.prototype = {
	constructor: MP.Connection
}
	
	/*
	hasSupport: function() {
	    if (window["WebSocket"]) {
			return true;
		} else {
			return false;
		}
	},
	connectTo: function(ip, port) {
		this.conn = new WebSocket("ws://localhost:8080/ws");
		this.conn.onclose = function (evt) {
			//appendLog($("<div><b>Connection closed.</b></div>"))
		}
		this.conn.onmessage = function (evt) {
			//appendLog($("<div/>").text(evt.data))
		}
	}
}
/*	connectTo: function(ip, port) {
		port = port || 8256;

		try {
			this.conn = new WebSocket( "ws://" + ip + ":" + port + "/ws" );
			this.conn.onopen = function(evt) { onOpen(evt); };
			this.conn.onmessage = function(evt) { onMessage(evt); };
			this.conn.onclose = function(evt) { onClose(evt); };
			this.conn.onerror = function(evt) { onError(evt); };
		} catch (exception) {
			message('<p class=\"err\">EXCEPTION:' + exception);
		}
	},
	close: function() {
		if ( this.conn !== undefined ) {
			this.conn.close();
			this.conn = undefined;
		}
	}
}

function message(text) {
	$('#log').append(text + "</p>");
}

function onOpen(evt) {
	writeToScreen("CONNECTED");
	doSend("WebSocket rocks");
}
function onClose(evt) { 
writeToScreen("DISCONNECTED");
} 
function onMessage(evt) {
writeToScreen('<span style="color: blue;">RESPONSE: ' + evt.data+'</span>');
conn.close();
} 
function onError(evt) {
writeToScreen('<span style="color: red;">ERROR:</span> ' + evt.data);
}
function writeToScreen(string) {
	$('#log').append(string);
}
*/
/*
$(function () {

    var conn;
    var msg = $("#msg");
    var log = $("#log");

    function appendLog(msg) {
        var d = log[0]
        var doScroll = d.scrollTop == d.scrollHeight - d.clientHeight;
        msg.appendTo(log)
        if (doScroll) {
            d.scrollTop = d.scrollHeight - d.clientHeight;
        }
    }

    $("#form").submit(function () {
        if (!conn) {
            return false;
        }
        if (!msg.val()) {
            return false;
        }
        conn.send(msg.val());
        msg.val("");
        return false
    });

    if (window["WebSocket"]) {
        conn = new WebSocket("ws://{{$}}/ws");
        conn.onclose = function (evt) {
            appendLog($("<div><b>Connection closed.</b></div>"))
        }
        conn.onmessage = function (evt) {
            appendLog($("<div/>").text(evt.data))
        }
    } else {
        appendLog($("<div><b>Your browser does not support WebSockets.</b></div>"))
    }
});
*/