package structures

import (
	"github.com/MetaParticle/metaparticle/logger"
	"code.google.com/p/go.net/websocket"
)

//This is a network.
//It's primary purpose is to group players into factions
//All messages is shared within a network.
//The network communicates with the world using the worlds broadcast channel
type Network struct {
	connecting chan Connection
	disconnecting chan Connection
	broadcast chan interface{}
	conns map[Connection]bool
	world *World
}

//Creates a new network for the world w.
func NewNetwork(w *World) *Network {
	n := new(Network)
	//n.connecting = make(chan Connection, 16)
	//n.disconnecting = make(chan Connection, 16)
	n.broadcast = make(chan interface{}, 64)
	n.conns = make(map[Connection]bool)
	n.world = w
	go n.start()
	return n
}

//DOES NOTHING
func (n Network) GetConnecting() chan Connection {
	return n.connecting
}

//DOES NOTHING
func (n Network) GetDisconnecting() chan Connection {
	return n.disconnecting
}

//DOES NOTHING
func getNetwork(w *World, ws *websocket.Conn) *Network {
	return nil
}

//Writes the slice to all registered connections.
func (net Network) Write(p []byte) (n int, err error) {
	for c := range net.conns {
				m, e :=  c.Write(p)
				n += m
				if e != nil {
					logger.Logf(5, "Error: %s", e)
					logger.Logf(6, "p: %s", string(p))
					go c.Close()
				}
	}
	return
}

//Adds a new connection to the network.
func (net Network) Register(c Connection) {
	net.conns[c] = true
	logger.Logf(6, "net.conns: %v", net.conns)
}

//Removes a connection from the network.
func (net Network) Unregister(c Connection) {
	delete(net.conns, c)
	close(c.GetChannel())
	c.Close()
	net.world.broadcast <- map[string] interface{} {
		"Remove": map[string] interface{} {
			"Id": c.GetPlayer().Id,
		},
	}
}

//Starts the network and handeles the messages on the channels.
func (n *Network) start() {
	for {
		select {
			case c := <-n.connecting:
				n.Register(c)
			case c := <-n.disconnecting:
				n.Unregister(c)
			case message := <-n.broadcast:
				for c := range n.conns {
						select {
							case c.GetChannel() <- message:
							default: 
								n.Unregister(c)
						}
				}
		}
	}
}
