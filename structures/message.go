package structures

import (
	"github.com/MetaParticle/metaparticle/entity"
	"encoding/json"
)

//A Message is a simple message.
//It has a sender id together with the data
//The data is JSON encoded information that is either a chatmessage or a updatepackage
type Message struct {
	SenderId uint32
	Data string
}

//Structure to send to new player.
type AssignMessage struct {
 	Id uint32
}

//Creates a new AssignMessage for a player.
func NewAssignMessage(p *entity.Player) []byte {
	mess, _ := json.Marshal(&AssignMessage{p.Id})
	return []byte("{\"Assign\":" + string(mess) + "}")
}
