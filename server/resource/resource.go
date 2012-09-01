package main

import (
	"fmt"
	"flag"
	"net"
	"os"
	"crypto/md5"
	
	"github.com/MetaParticle/metaparticle/cmdln"
)

var (
	dataFolder = flag.String("f", "/var/res", "Sets game data folder")
	remoteAddr  = flag.String("a", "localhost:8250", "set remote address (IP:port)")
	binFileName = flag.String("b", "clientserver", "set binary filename")
	port = flag.Int("p", 8250, "sets the port for version checking.")
	VERSION = "0.0.5 testing"
	binHash []byte
	hashRequests = 0
)

func main() {
	flag.Parse()
	
	err := os.Chdir(*dataFolder)
	if err != nil {
		cmdln.PrintStatus(-1, err.Error())
		os.Exit(1)
	}
	generateBinHash()
	
	
	cmdline, exit := cmdln.NewCommandline("RES > ")
	cmdline.PopulateDefaults()
	cmdline.AddCommand("hash", "access binary-hash functions", hashFunction)
	go cmdline.ListenAndServe()
	
	/*
	exit := make(chan bool)
	cmds := cmdln.NewCommandMap()
	cmds.AddCommand("hash", "access binary-hash functions", hashFunction)
	
	cmdline := cmdln.New(cmdln.PROMPT, cmds, exit)
	cmdline.SetPrompt("RES >")
	cmdline.PopulateDefaults()
	go cmdline.ListenAndServe()
	*/
	go listenWithBreak(supplyHash, exit)
	
	select {
		case <-exit:
	}
}

func generateBinHash() {
	bytes, err := fileToByteSlice(*binFileName)
	if err != nil {
		cmdln.PrintStatus(cmdln.ERROR, err.Error())
		os.Exit(1)
	}
	binHash = hash(bytes)
}

func hashFunction(args []string) {
	if len(args) != 0 {
		switch args[0] {
			case "print": printHash()
			case "update": generateBinHash()
				fmt.Println("New Hash generated.")
			case "requests": printHashReqs()
		}
		if args[0] != "help" {
			return
		}
	}
	fmt.Println("hash print    : Prints the current hash.")
	fmt.Println("hash update   : Updates the hash.")
	fmt.Println("hash requests : Prints the number of hash requests.")
}

func printHash() {
	fmt.Println(binHash)
}

func printHashReqs() {
	cmdln.PrintStatus(hashRequests, fmt.Sprintf("So far %d requests.", hashRequests))
}

// Helpers
func hash(b []byte) []byte {
	hasher := md5.New()
	hasher.Write(b)
	res := make([]byte, 0)
	return hasher.Sum(res)
}

func fileToByteSlice(filename string) ([]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	fi, err := file.Stat()
	if err != nil {
		return nil, err
	}
	data := make([]byte, fi.Size())
	_, err = file.Read(data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func listenWithBreak(handleFn func(net.Conn), brk chan bool) {
	addr, err := net.ResolveTCPAddr("tcp", *remoteAddr)
	if err != nil {
		cmdln.PrintStatus(cmdln.ERROR, err.Error())
		brk <- true
	}
	ln, err := net.ListenTCP("tcp", addr)
	if err != nil {
		cmdln.PrintStatus(cmdln.ERROR, err.Error())
		brk <- true
	}
	for {
		conn, err := ln.AcceptTCP()
		if err != nil {
			// handle error
			continue
		}
		hashRequests++
		go handleFn(conn)
	}
}

func supplyHash(conn net.Conn) {
	conn.Write(binHash)
	conn.Close()
}

