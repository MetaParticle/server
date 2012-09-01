package bal

import (
	"fmt"
	"flag"
	"net"
	"net/http"

	"os"
	"os/user"
	//"os/exec"
	
	"io"

	"metaparticle/cmdln"
	"metaparticle/config"
)

const (
	NAME = "mpl, MetaParticle Launcher"
)

var (
	usr, _ = user.Current()
	dataFolder = flag.String("f", usr.HomeDir + "/.metaparticle", "Sets game data folder.")
	confName = flag.String("c", "metaparticle.conf", "Sets the configuration file.")
	VERSION = "0.0.2 testing"
)

func main() {
	flag.Parse()

	fmt.Println(NAME)
	fmt.Printf("Version: %s\n", VERSION)
	
	checkFolder(dataFolder)
	
	//configFile, newFile := getConfigFile(confName)
	configFile := config.GetConfigFile(*confName)
	if configFile != nil {
		defer configFile.Close()
	
		var confMap config.Config
		if (configFile.Empty()) {
			confMap = config.Config{}
			confMap["host"]="localhost:8250"
			confMap["bin"]="mpcsd"
			configFile.Marshal(confMap)
				
		} else {
			confMap = configFile.Unmarshal()
		}	
	
		fmt.Println(confMap)
		checkForUpdate(confMap["bin"], confMap["host"])
		download := cmdln.AskYN("Do you want to download the latest version now?")
		if !download {
			cmdln.PrintStatus(cmdln.ERROR, "The latest version")
			os.Exit(0)
		}
		res, err := http.Get("http://"+confMap["host"]+"/bin/mpcsd")
		if err != nil {
			cmdln.PrintStatus(cmdln.ERROR, "Couldn't connect to resource server.")
			return
		}
		defer res.Body.Close()
		
		//Open and promise to close destination file
		newbin, err := os.OpenFile(confMap["bin"], os.O_WRONLY | os.O_CREATE | os.O_TRUNC, 0666)
		defer newbin.Close()
		
		if (err != nil) {
			fmt.Printf("Error opening file %s for writing: %s\n", confMap["bin"], err)
			return
		}
		
		_, err = io.Copy(newbin, res.Body)
		
		if err != nil {
			fmt.Printf("Error while writing file %s: %s\n", confMap["bin"], err)
		}
		
		/* HTTP GET EXAMPLE
		res, err := http.Get("http://www.google.com/robots.txt")
		if err != nil {
			log.Fatal(err)
		}
		robots, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}
		res.Body.Close()
		fmt.Printf("%s", robots)
		*/
	}
}

func checkFolder(folder *string) {
	err := os.Chdir(*folder)
	if err != nil {
		os.Mkdir(*folder, 0x776)
		cmdln.PrintStatus(cmdln.FAULT, "Attempting to create .metaparticle")
		err := os.Chdir(*folder)
		if err != nil {
			cmdln.PrintStatus(cmdln.ERROR, err.Error())
			cmdln.PrintStatus(cmdln.ERROR, "Creation failed.")
			os.Exit(1)
		} else {
			cmdln.PrintStatus(cmdln.SUCCESS, "Great success!")
		}
	} else {
		cmdln.PrintStatus(cmdln.SUCCESS, ".metaparticle found")
	}
}

func getConfigFile(filename *string) (*os.File, bool) {
	file, err := os.Open(*filename)
	if err != nil {
		file, err = os.Create(*filename)
		cmdln.PrintStatus(cmdln.FAULT, "Attempting to create " + *filename)
		if err != nil {
			cmdln.PrintStatus(cmdln.ERROR, err.Error())
			cmdln.PrintStatus(cmdln.ERROR, "Creation failed.")
			os.Exit(1)
		} else {
			cmdln.PrintStatus(cmdln.SUCCESS, "Great success!")
			return file, true
		}
	}
	cmdln.PrintStatus(cmdln.SUCCESS, *filename + " found.")
	return file, false
}

func getHashFromServer(host string) []byte {
	hostAddr, err := net.ResolveTCPAddr("tcp", host)
	if err != nil {
		cmdln.PrintStatus(cmdln.ERROR, err.Error())
		os.Exit(1)
	}

	// DIAL
	conn, err := net.DialTCP("tcp", nil, hostAddr)
	fmt.Println(host)
	if err != nil {
		cmdln.PrintStatus(cmdln.ERROR, "Couldn't connect to MetaParticle resource server.")
		os.Exit(1)
	}
	defer conn.Close()
	
	// READ
	hash := make([]byte, GetHashSize())
	_, err = conn.Read(hash)
	if err != nil {
		cmdln.PrintStatus(cmdln.ERROR, "Couldn't read data from MetaParticle resource server.")
		os.Exit(1)
	}
	return hash
}

func checkForUpdate(binName, serverIP string) {
	// Load local binary
	bin, err := fileToByteSlice(binName)
	
	if err != nil {
		cmdln.PrintStatus(cmdln.FAULT, "Couldn't read local binary data.")
		return
	}

	// Get server hash
	hash := getHashFromServer(serverIP)
	
	// Compare hash to binary, which the function will hash.
	if CompareHash2Bin(hash, bin) {
		cmdln.PrintStatus(cmdln.FAULT, "There is a client update available.")
		return
	} else {
		cmdln.PrintStatus(cmdln.SUCCESS, "The client is up to date.")
	}
	//os.Exit(0) // All is good!
	
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
