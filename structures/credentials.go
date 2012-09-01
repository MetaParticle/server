package structures

import (
	//"encoding/json"
	
	"code.google.com/p/go.net/websocket"
)

var idGen = make(chan uint32, 32)

func init() {
	go GenNewIds(idGen)
}

//Structure used to group id and password when logging in.
type Credentials struct {
	Id uint32
	Password string
}

//Gets credetials using a websocket connection.
func GetWSCreds(ws *websocket.Conn) (creds *Credentials, err error) {
	/*
	dec := json.NewDecoder(ws)
	err = dec.Decode(creds) //creds has its init values if malformed data is recived according to the specs.
	*/
	return &Credentials{uint32(<-idGen), "hunter02"}, nil
}

//Safe function to get the id and password from a Credetials.
func (cred Credentials) Split() (uint32, string) {
	return cred.Id, cred.Password
}

//Id generator.
func GenNewIds(c chan<- uint32) {
	i := uint32(1)
	for {
		c <- i
		i++
	}
}
