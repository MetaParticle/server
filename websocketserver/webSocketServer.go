package websocketserver

import (
	"github.com/MetaParticle/metaparticle/logger"
	"github.com/MetaParticle/metaparticle/entity"
	"github.com/MetaParticle/metaparticle/structures"

	"fmt"
	"net/http"
	//"encoding/json"
	

	"code.google.com/p/go.net/websocket"
)

const (
	//The standard port for websocket communications
	STDPORT = 8256
)

func getNetwork(w *structures.World, ws *websocket.Conn) *structures.Network {
	return w.GetNetwork(0)
}

//Using a websocket connection, credetials is accuired and
//used to get a player from the database.
//If no player is found in the database, a new player is created.
func GetPlayer(ws *websocket.Conn) *entity.Player {
	cred, _ := structures.GetWSCreds(ws)
	return entity.GetPlayer( cred.Split() )
}

// Echo the data received on the Web Socket.
func GenNewConnectionHandeler(w *structures.World) func(*websocket.Conn) {
	logger.Log(5, "WS: Handlerfunction registered for the world!")
    return func(ws *websocket.Conn) {
        logger.Logf(3, "WS: Connection established from %s!", ws.RemoteAddr().String())
        n := getNetwork(w, ws)
        p := GetPlayer(ws)
        c := structures.NewWSConnection(p, n, ws, make(chan interface{}, 7))
        n.Register(c)
        defer func() { n.GetDisconnecting() <- c }()
        ws.Write(structures.NewAssignMessage(p))
        go c.Writer()
        c.Reader()
    }
}

//This starts the websocket server.
//It should be started using a goroutine due to it's blocking nature.
func ListenAndServe(port int, errchan chan error) {
    w := structures.GetWorld()
    //registration of websocket handelers
    wsMux := http.NewServeMux()
    wsMux.Handle("/", websocket.Handler( GenNewConnectionHandeler(w) ))
    
    //Listen on websocket port
    errchan <- http.ListenAndServe(":"+fmt.Sprint(port), wsMux)
    return
}
