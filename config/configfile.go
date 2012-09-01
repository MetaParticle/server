package config

import (
	"os"
	"bufio"
	"strings"
	"math"
	"path/filepath"
)
type ConfigFile struct {
	*os.File //Ville göra en rak typ, men det blev så många konverteringar som failade...
}
type Config map[string]string

func GetConfigFile(path string) *ConfigFile {
	os.MkdirAll(filepath.Dir(path), 0666)
	file, err := os.OpenFile(path, os.O_CREATE | os.O_RDWR, 0666) //Ser till så att filen finns samt är läs och skrivbar.
	if err == nil {
		return &ConfigFile{file}
	}
	return nil
	
}

func (f ConfigFile) Unmarshal() Config {
	c := Config{}
	bio := bufio.NewReader(f.File)
	
	for {
		bytes, _, err := bio.ReadLine()
		if err != nil {
			break
		}
		line := string(bytes)
		lineAndComment := strings.SplitN(line, "#", 2)
		segs := strings.SplitN(lineAndComment[0], "=", 2)
		if len(segs) > 1 {
			c[segs[0]]=segs[1]
		}
	}
	
	f.Close()
	return c
}

func (f ConfigFile) Marshal(c Config) {
	for key, value := range c {
		_, err := f.File.WriteString(key + "=" + value + "\n")
		if err != nil {
			// Handle error
		}
	}
	//f.Close() //varför? Borde inte det här skötas någon annanstans?
		//typ av den som öppnat filen. 
		//Man kan tänkas vilja marshalla kontinuerligt, utan att behöva öpnna om filen varje gång.
		//ROPEN SKALLA, LÅT FILEN STANNA!

		// Gäller detsamma inte för unmarshal?
}

func (f ConfigFile) Size() int64 {
	stat, err := f.File.Stat()
	if err == nil {
		return stat.Size()
	}
	return int64(math.NaN())
}

func (f ConfigFile) Empty() bool {
	return f.Size() == 0
}


// -Oskar
// Behövs inte Ovanstående funktioner räcker gott! :) 

// -Zephyyrr
// Fine, men då vill jag ha metoder istället!
// Och utflyttat till separat paket. Det här kan komma till användning på annat håll också.

// -Oskar
// Låter fint! Antingen får man data ur config filen, eller så får man väl
// skapa en egen. Och för det skulle nog metoder passa fint.

/*
func LoadConfig(path string) Config {
	file, err := getFile(path)
	defer file.Close()
	
	return Unmarshal(file)
}

func getFile(path string) *os.File {
	return nil
}
*/
