package main

import (
	"github.com/MetaParticle/metaparticle/logger"
	"github.com/MetaParticle/metaparticle/entity"
	
	"flag"
	"time"
)

var (
	systemName = flag.String("s", "", "The system this server represents")
)

type System struct {
	name string
	ships map[entity.Ship]bool
	parent *System
	children map[string]*System
}

func NewSystem(name string) *System {
	return &System{ name, make(map[entity.Ship]bool), nil, make(map[string]*System)}
}

func main() {
	flag.Parse()
	
	sys := LoadSystem(*systemName)
	
	go sys.RunPhysics()
	ListenAndServe()
}

/*
  Loads the system with the given name from storage and returns it
  */
func LoadSystem(name string) *System {
	//TODO Actually load it from persistant storage.
	return NewSystem(name)
}

func (sys *System) RunPhysics() {
	ticker := time.Tick(time.Second)
	var dt time.Duration = 0
	for t := range ticker {
		logger.Logf(logger.INFO, "Time since last physics update: %G ns", dt)
		//TODO DO MAGICAL PHYSICS STUFFS ON CHILDREN AND SHIPS.
		//WILL SOMEONE PLEASE THINK OF THE CHILDREN!?!?!?!?
		dt = time.Since(t)
	}
}

func ListenAndServe() {
	//TODO LISTEN FOR MAGICAL BEINGS FROM THE GREAT UNKNOWN AND ADD THEIR SHIPS TO MAP.
}
