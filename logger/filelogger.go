package logger

import (
	"os"
	"io"
)

type FileLogger struct {
	*Loggie
	*os.File
}

func NewFileLogger(name string) (*FileLogger, error) {
	os.MkdirAll(LOGFOLDER, 438)
	f, err := os.OpenFile(LOGFOLDER + "/" + name + ".log", os.O_WRONLY | os.O_CREATE | os.O_APPEND, 438)
	if err != nil {
		return nil, err
	}
	f.Write([]byte("\n--------\n"))
	dst := io.MultiWriter(os.Stdout, f)
	return &FileLogger{NewLoggie(name, dst), f}, nil
}

func (fl FileLogger) Close() error {
	return fl.File.Close()
}

func (fl FileLogger) Rotate() {
	
}
