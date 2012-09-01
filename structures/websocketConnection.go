package structures

import (
	"github.com/MetaParticle/metaparticle/logger"
	"github.com/MetaParticle/metaparticle/entity"
	
	//"encoding/json"
	
	"code.google.com/p/go.net/websocket"
)

// Implementation of Connection using websockets
type WSConnection struct {
	p *entity.Player
	n *Network
	connection *websocket.Conn
	send chan interface{}
	online bool
}

// Used to create a new WSConnection.
func NewWSConnection(player *entity.Player, n *Network, ws *websocket.Conn, send chan interface{}) *WSConnection {
	wsc := new(WSConnection)
	wsc.p = player
	wsc.n = n
	wsc.connection = ws 
	wsc.send = send
	wsc.online = true
	return wsc
}

// Gets the player from the connection.
func (wsc WSConnection) GetPlayer() *entity.Player {
	return wsc.p
}

// Gets the channel used to communicate with the player.
func (wsc WSConnection) GetChannel() chan interface{} {
	return wsc.send
}

// This makes Connection implement io.Writer
func (wsc WSConnection) Write(p []byte) (n int, err error) {
	return wsc.connection.Write(p)
}

// Goable blocking function that looks for messages on the channel and writes it in JSON format.
func (wsc WSConnection) Writer() {
	for message := range wsc.send {
		err := websocket.JSON.Send(wsc.connection, message)
		if err != nil {
			logger.Logf(5, "Error on write: %s", err.Error())
			break
		}
	}
}

// Blocking function that reads JSON objects from the connection and broadcasts them on the network.
func (wsc WSConnection) Reader() {
	for {
		var data interface{}
		err := websocket.JSON.Receive(wsc.connection, &data)
		logger.Logf(10, "Data from player: %d", wsc.p.Id)
		logger.Logf(10, "%v", data)
		wsc.n.world.broadcast <- &data
		if err != nil {
			logger.Logf(5, "Error on read: %s", err.Error())
			break
		}
	}
	go wsc.Close()
}

// Closes the connection and frees system resources connected to it.
func (wsc WSConnection) Close() {
	if wsc.online {
		wsc.online = false
		logger.Logf(5, "Closing connection to: %s", wsc.connection.RemoteAddr())
		wsc.connection.Close()
	}
}
