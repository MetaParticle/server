package structures

import "github.com/MetaParticle/metaparticle/entity"

// Interface used to define what a connection is.
// This allows a connection to be made using any type of connection that can reliably be established.
type Connection interface {
		GetPlayer() *entity.Player
		GetChannel() (chan interface{})
		Write(p []byte) (n int, err error)
		Writer()
		Reader()
		Close()
}
