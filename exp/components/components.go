package components



type Component struct {
	inputs []<-chan Flow
	listeners map[int]FlowWriter
	function func(Flow) Flow
}

func (c Component) registerListener(place int, f FlowWriter) bool {
	if (c.listeners[place] == nil) {
		c.listeners[place] = f
		return true
	}
	return false
}

func (c Component) Send(place int, f *Flow) {
	if (c.listeners[place] != nil) {
		c.listeners[place].Send(f)
	}
}

type FlowWriter interface {
	Send(*Flow)
}

type ChanIO chan Flow

func (c ChanIO) Send(f Flow)  {
	c<-f
}

func (c ChanIO) Write()

type Flow struct {
	kind string
	amount byte
}
