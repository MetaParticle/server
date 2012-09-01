package main

import (
	"github.com/MetaParticle/metaparticle/logger"
	ws "github.com/MetaParticle/metaparticle/websocketserver"
	www "github.com/MetaParticle/metaparticle/wwwserver"
	"github.com/MetaParticle/metaparticle/cmdln"
	"github.com/MetaParticle/metaparticle/config"
	
	"os"
	"os/user"
	"flag"
	"strconv"
)

const (
	VERSION_NUMBER = 4.3 //Update as necessary.
)

var (
	//logger
	loglvl = int(logger.IMPORTANT)
	l logger.LoggerCloser

	//websocketserver
	wsport = ws.STDPORT

	//wwwserver
	port = www.STDPORT//Port 80 needs root btw
	root =  www.LOCALROOT
	
	usr, _ = user.Current()
	dataFolder = flag.String("f", usr.HomeDir + "/.metaparticle", "Sets game data folder.")
	confName = flag.String("c", "metaparticle_meta.conf", "Sets the configuration file.")
	
)

//Function used to print error messages if http servers crash.
func CrashAndBurn(err error, eport int) {
	l.Logf(logger.NORMAL, "ListenAndServe port %d: %s", eport, err.Error())
	l.Log(logger.NORMAL, "Please check that no other instance is already running on this machine.")
}

func init() {
	os.MkdirAll(*dataFolder + "/logs", 0666)
	var err error
	l, err = logger.NewFileLogger("meta")
	if err != nil {
		l = logger.NewLoggie("meta", os.Stdout)
		logger.Logf(logger.MINOR, "Could not open logfile: %s\n Continuing without logging to file.", err)
	}
	logger.Logf(0, "l: %v", l)
	l.SetLogLevel(logger.GetLogLevel(loglvl))
}

func LoadConfig() {
	os.MkdirAll(*dataFolder, 0666)
	configFile := config.GetConfigFile(*dataFolder + "/" + *confName)
	if configFile != nil {
		var c config.Config
		defer configFile.Close()
		if (configFile.Empty()) {
			c = make(map[string]string)
			c["root"]=root
			c["port"]=strconv.Itoa(port)
			c["wsport"]=strconv.Itoa(wsport)
			c["loglevel"]=strconv.Itoa(loglvl)
			configFile.Marshal(c)
				
		} else {
			c = configFile.Unmarshal()
		}
		
		setVars(c)
	}
}

func setVars(c config.Config) {
	root = c["root"]
	t, err := strconv.Atoi(c["port"])
	if err == nil {port=t} else {logger.Logf(logger.NORMAL, "Error while setting port to %v: %s", t, err.Error())}
	
	t, err = strconv.Atoi(c["wsport"])
	if err == nil {wsport=t} else {logger.Logf(logger.NORMAL, "Error while setting wsport to %v: %s", t, err.Error())}
	
	t, err = strconv.Atoi(c["loglevel"])
	if err == nil {loglvl=t} else {logger.Logf(logger.NORMAL, "Error while setting loglevel to %v: %s", t, err.Error())}
}


func main() {
	//Startup code
	
	flag.Parse()
	LoadConfig()

	defer l.Close()
    
	l.Logf(logger.NORMAL, "Meta Particle Server %G", VERSION_NUMBER)
	l.Logf(logger.NORMAL, "Loglevel: %d", logger.GetCurrentLogLevel())
	l.Logf(logger.NORMAL, "WWW Root: %s", root)
	//logger.Logf(logger.NORMAL, "Resource Root: %s", *resroot)
	l.Logf(logger.NORMAL, "WWW Port: %d", port)
	l.Logf(logger.NORMAL, "WebSocket port: %d", wsport)
	l.Println("\n")

	wwwerr := make(chan error)
	wserr := make(chan error)
	quitchan := make(chan bool)
	
	defer func() {
		if r := recover(); r != nil {
			l.Logf(logger.SEVERE, "Uncaught panic found in main.")
		}
	}()

	go www.ListenAndServe(root, port, wwwerr)
	go ws.ListenAndServe(wsport, wserr)
	
	cmdline:=cmdln.New(cmdln.PROMPT, cmdln.NewCommandMap(), quitchan)
	cmdline.PopulateDefaults()
	go cmdline.ListenAndServe()
	
	select {
		case err := <- wwwerr: CrashAndBurn(err, port)
		case err := <- wserr: CrashAndBurn(err, wsport)
		case <- quitchan: l.Log(logger.NORMAL, "Quit command recived. Shutting down.")
	}
	
}
