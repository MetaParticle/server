package entity

import (
	"crypto/sha512"
	//"hash"
	//"fmt"
	"math"
)

var (
	hasher = sha512.New()
)

//This is a player.
//It stores all information uniqe for this player.
type Player struct {
	Id      uint32
	Name    string
	Ship string
	Email string
	Hashpass string
	//Insert more stuff as necessary.
}

//Using a userid and password, gets the player linked to that account.
//If credentials is incorrect, nil is returned.
func GetPlayer(userid uint32, password string) *Player {
	return GetTestPlayer()
}

func GetTestPlayer() *Player {
	p := new(Player)
	p.Name = "Mall Groda"
	p.Id = 1
	p.Ship = "Caelia"
	p.Hashpass = hashPassword("hunter02")
	return p
}

func hashPassword(password string) string {
	buf := extendToBlock(hasher.BlockSize(), password)
	//fmt.Println(buf)
	hasher.Write(buf)
	hashbuf:= make([]byte, hasher.Size())
	defer hasher.Reset()
	return string(hasher.Sum(hashbuf))
}

func extendToBlock(size int, s string) []byte {
	blocks := int(math.Ceil(float64(len(s))/float64(size)))
	buf := make([]byte, 0,blocks*size)
	fillsize := cap(buf)-len(s)
	buf = append(buf, []byte(s)...)
	fill(&buf, fillsize, s)
	return buf
}

func fill(buf *[]byte, fillsize int, s string) {
	times := fillsize/len(s)
	rest := fillsize%len(s)
	
	//fmt.Printf("times: %d, rest: %d\n", times, rest)
	
	for i := 0; i < times; i++ {
		*buf = append(*buf, []byte(s)...)
	}
	*buf = append(*buf, []byte(s[:rest])...)
	//fmt.Println("Size:", len(*buf))
}

