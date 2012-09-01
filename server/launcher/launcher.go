package main

import (
	"fmt"
	"flag"
	"net"
	"net/http"

	"os"
	"os/user"
	"os/exec"
	
	"io"

	"errors"

	"github.com/MetaParticle/metaparticle/cmdln"
	"github.com/MetaParticle/metaparticle/config"
)

const (
	NAME = "mpl, MetaParticle Launcher"
)

var (
	usr, _ = user.Current()
	dataFolder = flag.String("f", usr.HomeDir + "/.metaparticle", "Sets game data folder.")
	confName = flag.String("c", "metaparticle.conf", "Sets the configuration file.")
	VERSION = "0.0.17 testing"
)

func main() {
	flag.Parse()

	fmt.Println(NAME)
	fmt.Printf("Version: %s\n", VERSION)
	
	err := cd(dataFolder)
	if err != nil {
		cmdln.PrintStatus(cmdln.ERROR, err.Error())
		return
	}
	cmdln.PrintStatus(cmdln.SUCCESS, "~/.metaparticle found.")
	
	conf, err := getConfMap(*confName)
	if err != nil {
		cmdln.PrintStatus(cmdln.ERROR, err.Error())
		return
	}
	
	different, err := updateExists(conf["host"], conf["bin"])
	if err != nil {
		cmdln.PrintStatus(cmdln.ERROR, err.Error())
		return
	}
	
	if different {
		cmdln.PrintStatus(cmdln.FAULT, "New version available.")
		if !cmdln.AskYN("Do you want to download the latest version now?") {
			cmdln.PrintStatus(cmdln.ERROR, "To play, you must have the latest version.")
			return
		}
		err := updateFile(conf["host"], conf["bin"])
		if err != nil {
			cmdln.PrintStatus(cmdln.ERROR, err.Error())
			return
		}
	}
	cmdln.PrintStatus(cmdln.SUCCESS, conf["bin"] + " is up to date.")
	
	cmd := exec.Command(conf["bin"])
	cmd.Run()
}

func cd(folder *string) error {
	err := os.Chdir(*folder)
	if err != nil {
		os.Mkdir(*folder, os.ModePerm)
		cmdln.PrintStatus(cmdln.FAULT, "Attempting to create " + *folder)
		return os.Chdir(*folder)
	}
	return nil
}

func getConfMap(fileName string) (config.Config, error) {
	confFile := config.GetConfigFile(fileName)
	if confFile == nil {
		return nil, errors.New("Couldn't open config file.")
	}
	defer confFile.Close()
	
	var conf config.Config
	if (confFile.Empty()) {
		conf = config.Config{}
		conf["host"]="localhost:8250"
		conf["bin"]="mpcsd"
		confFile.Marshal(conf)
	} else {
		conf = confFile.Unmarshal()
	}
	return conf, nil
}

func updateExists(server, name string) (bool, error) {
	localHash, err := getHashFromFile(name)
	if err != nil {
		switch err.(type) {
			case *os.PathError:
				return true, nil			
			default:
				return false, err
		}
	}
	serverHash, err := getHashFromServer(server)
	if err != nil {
		return false, err
	}
	return !compareHash(localHash, serverHash), nil
}

func getHashFromFile(fileName string) ([]byte, error) {
	file, err := os.Open(fileName)
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
	return hash(data), nil
}

func getHashFromServer(host string) ([]byte, error) {
	hostAddr, err := net.ResolveTCPAddr("tcp", host)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialTCP("tcp", nil, hostAddr)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	
	hash := make([]byte, 16)
	_, err = conn.Read(hash)
	if err != nil {
		return nil, err
	}
	return hash, nil
}

func updateFile(server, filename string) error {
	res, err := http.Get("http://"+server+"/bin/mpcsd")
	if err != nil {
		return err
	}
	defer res.Body.Close()
	
	newbin, err := os.OpenFile(filename, os.O_WRONLY | os.O_CREATE | os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	defer newbin.Close()
	
	_, err = io.Copy(newbin, res.Body)
	return err
}
