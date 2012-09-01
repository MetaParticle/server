package main

import (
	"github.com/MetaParticle/metaparticle/structures"	
	
	"fmt"
	"encoding/json"
)

func main() {
	mess, _ := json.Marshal(&structures.AssignMessage{1337})
	fmt.Printf("b1: %s\n", "assign:" + string(mess))
}
