  var secureCb;
  var secureCbLabel;
  var wsUri;
  var consoleLog;
  var connectBut;
  var disconnectBut;
  var sendMessage;
  var sendBut;
  var clearLogBut;

  function doConnect()
  {
    if (uri.indexOf("?") == -1) {
        uri += "?encoding=text";
    } else {
        uri += "&encoding=text";
    }
    websocket = new WebSocket(uri);
    websocket.onopen = function(evt) { onOpen(evt) };
    websocket.onclose = function(evt) { onClose(evt) };
    websocket.onmessage = function(evt) { onMessage(evt) };
    websocket.onerror = function(evt) { onError(evt) };
  }
  
  function doDisconnect()
  {
    websocket.close()
  }
  
  function doSend()
  {
    logToConsole("SENT: " + sendMessage.value);
    websocket.send(sendMessage.value);
  }

  function logToConsole(message)
  {
    var pre = document.createElement("p");
    pre.style.wordWrap = "break-word";
    pre.innerHTML = getSecureTag()+message;
    consoleLog.appendChild(pre);

    while (consoleLog.childNodes.length > 50)
    {
      consoleLog.removeChild(consoleLog.firstChild);
    }
    
    consoleLog.scrollTop = consoleLog.scrollHeight;
  }
  
  function onOpen(evt)
  {
    logToConsole("CONNECTED");
    setGuiConnected(true);
  }
  
  function onClose(evt)
  {
    logToConsole("DISCONNECTED");
    setGuiConnected(false);
  }
  
  function onMessage(evt)
  {
    logToConsole('<span style="color: blue;">RESPONSE: ' + evt.data+'</span>');
  }

  function onError(evt)
  {
    logToConsole('<span style="color: red;">ERROR:</span> ' + evt.data);
  }
  
  function setGuiConnected(isConnected)
  {
    wsUri.disabled = isConnected;
    connectBut.disabled = isConnected;
    disconnectBut.disabled = !isConnected;
    sendMessage.disabled = !isConnected;
    sendBut.disabled = !isConnected;
    secureCb.disabled = isConnected;
    var labelColor = "black";
    if (isConnected)
    {
      labelColor = "#999999";
    }
     secureCbLabel.style.color = labelColor;
    
  }
	
	function clearLog()
	{
		while (consoleLog.childNodes.length > 0)
		{
			consoleLog.removeChild(consoleLog.lastChild);
		}
	}
doConnect();
