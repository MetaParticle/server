package logger

import (
	"log"
	"fmt"
	"os"
	"os/user"
	"io"
)

const (
	//The standard logging level.
	//This means that IMPORTANT classified logs is produced
	STDLVL = loglevel(5)
)

var (
	LOGFOLDER = "log"
)

func init() {
	user, err := user.Current()
	if err == nil {
		LOGFOLDER = user.HomeDir + "/.metaparticle/logs"
	}
}

func GetLogLevel(lvl int) loglevel {
	return loglevel(lvl)
}

type LoggerCloser interface {
	GetCurrentLogLevel() loglevel
	SetLogLevel(loglevel)
	AddDestination(w io.Writer)
	Log(loglevel, string)
	Logf(loglevel, string, ...interface{})
	Printf(string, ...interface{})
	Println(string)
	Nl()
	
	Panic(string, ...interface{})
	Error(string, ...interface{})
	
	io.Closer
}

type Loggie struct {
	name string
	loglvl loglevel
	dst io.Writer
	logger *log.Logger
}

func NewLoggie(name string, dst io.Writer) (l *Loggie){
	l = new(Loggie)
	l.name = name
	l.dst = dst
	l.loglvl = STDLVL
	l.logger = log.New(l.dst, "[" + name + "]", log.LstdFlags)
	return
}

//Returns the loggies level as an int
func (l *Loggie) GetCurrentLogLevel() loglevel {
	return l.loglvl
}

//Sets to loggies loglevel
func (l *Loggie) SetLogLevel(level loglevel) {
	l.loglvl = level
}

//Makes loggie write to the new writer aswell. 
func (l *Loggie) AddDestination(w io.Writer) {
	l.dst = io.MultiWriter(l.dst, w)
	l.logger = log.New(l.dst, "[" + l.name + "]", log.LstdFlags)
}

//Prints a log message if the level is below the current Loglevel.
//Messege is printed in format "Log <level>: <message>" and 
//ends with a new line.
func (l Loggie) Log(level loglevel, message string) {
	if l.loglvl >= level {
		l.logger.Printf("Log %v: %s\n", level, message)
	}
}


//Log function using the Printf function of log rather than Print.
//Ends with a newline
func (l Loggie) Logf(level loglevel, message string, params ...interface{}) {
		fmessage := fmt.Sprintf(message, params...)
		l.Log(level, fmessage)
}

func (l Loggie) Printf(message string, params ...interface{}) {
	fmessage := fmt.Sprintf(message, params...)
	fmt.Fprintf(l.dst, fmessage + "\n")
}

func (l Loggie) Println(message string) {
	fmt.Fprintln(l.dst, message)
}

func (l Loggie) Nl() {
	fmt.Fprintln(l.dst, "")
}

func (l Loggie) Panic(message string,  params ...interface{}) {
	fmess := fmt.Sprintf(message, params...)
	l.Log(PANIC, fmess)
	panic(fmess)
}

func (l Loggie) Error(message string, params ...interface{}) {
	fmess := fmt.Sprintf(message, params...)
	l.Log(ERROR, fmess)
	os.Exit(1)
}

func (l Loggie) Close() error {
	//Does nothing. Only to match LoggerCloser interface.
	return nil
}
