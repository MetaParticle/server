package structures

var world = NewWorld()


//This is a world.
//It represents a world of players.
//It has a list of networks and a channel for broadcasting to all players.
type World struct {
	nets [](*Network)
	broadcast chan interface{}
}

//This creates a new world an creates its first network.
func NewWorld() *World {
	w := new(World)
	w.nets = make([]*Network, 1)
	w.nets[0] = NewNetwork(w)
	w.broadcast = make(chan interface{}, 256)
	go w.launch()
	return w
}

//This gets the default world.
func GetWorld() *World {
	return world
}

//Gets the network i of the world.
//if the network does not exist, nil is returned.
func (w *World) GetNetwork(i uint) *Network {
	return w.nets[0]
}

//This method launches a world and lets it start handeling broadcasts.
func (w *World) launch() {
	for {
		message := <-w.broadcast
		for _, n := range w.nets {
			n.broadcast <- message
		}
	}
}
