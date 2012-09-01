package logger

import (
	"os"
	"io"
)

var stdlog *Loggie

func init() {
	stdlog = NewLoggie("Loggie", os.Stdout)
	stdlog.SetLogLevel(STDLVL)
}

//Returns the loggies level as an int
func GetCurrentLogLevel() loglevel {
	return stdlog.loglvl
}

//Sets to loggies loglevel
func SetLogLevel(level loglevel) {
	stdlog.SetLogLevel(level)
}

//Makes loggie write to the new writer aswell. 
func AddDestination(w io.Writer) {
	stdlog.AddDestination(w)
}

//Prints a log message if the level is below the current Loglevel.
//Messege is printed in format "Log <level>: <message>" and 
//ends with a new line.
func Log(level loglevel, message string) {
	stdlog.Log(level, message)
}


//Log function using the Printf function of log rather than Print.
//Ends with a newline
func Logf(level loglevel, message string, params ...interface{}) {
	stdlog.Logf(level, message, params...)
}

func Printf(message string, params ...interface{}) {
	stdlog.Printf(message, params...)
}

func Println(message string) {
	stdlog.Println(message)
}

func Nl() {
	stdlog.Nl()
}

func Panic(message string,  params ...interface{}) {
	stdlog.Panic(message, params...)
}

func Error(message string, params ...interface{}) {
	stdlog.Error(message, params...)
}
