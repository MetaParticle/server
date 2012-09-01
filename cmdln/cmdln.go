package cmdln

import (
	"github.com/MetaParticle/metaparticle/logger"
	
	"fmt"
	"strings"
	//"strconv"
	"bufio"
	"os"
	//"io"
)

const (
	PROMPT = BoldOn + "MetaParticle" + BoldOff + " > "
	ERROR = -1
	FAULT = 0
	SUCCESS = 1
)

type CommandLine struct {
	ComMap
	Prompt string
	
	exit chan bool
	stop bool
}

func NewCommandline(prompt string) (cmdline *CommandLine, exit chan bool) {
	exit = make(chan bool)
	cmdline = &CommandLine{make(ComMap), prompt, exit, false}
	return
}

func New(prompt string, cmds *ComMap, quitchan chan bool) (cmdline *CommandLine) {
	cmdline = &CommandLine{*cmds, prompt, quitchan, false}
	//go io.Copy(cmdline, os.Stdout)
	
	/*
	 Jag vill att all output skall gå genom den här, 
	 så att prompten kan tryckas ut efteråt.
	 Det verkar inte vilja gå just nu... synd.
	 */
	
	return
}

func (c *CommandLine) ListenAndServe() {
	fmt.Println("For help using the commandline, use \"help\".")
	reader := bufio.NewReader(os.Stdin)
	
	var line string	
	var cmd *command
	for !c.stop {	
		c.printPrompt()
		line = readLine(reader); 
		input := strings.Split(line, " ")
		
		cmd = c.ComMap[input[0]]
		if cmd != nil {
			logger.Logf(1, "Console used command: %s", line)
			cmd.function(input[1:])
		} else {
			c.printHelp()
		}
	}
	fmt.Println("End of the line! All passangers, prepare for landing!")
	c.exit <- true
}

func (c CommandLine) Write(b []byte) (n int, err error) {
	fmt.Print(string(b))
	c.printPrompt()
	return
}

func (c CommandLine) Printfln(mess string, args ...interface{}) {
	fmt.Printf(mess + "\n", args...)
	c.printPrompt()
}

func (c CommandLine) printPrompt() {
	fmt.Print(c.Prompt)
}

func (c *CommandLine) SetPrompt(prompt string) {
	c.Prompt = prompt
}

func (c *CommandLine) PopulateDefaults() {
	c.ComMap["help"] = &command{"Prints this help text", func(args []string){c.printHelp()}}
	c.ComMap["quit"] = &command{"Shuts down the server nicely", func(args []string){c.stop=true}}
}

type commandFunc func([]string)

type command struct {
	desc string
	function commandFunc
}

type ComMap map[string]*command

func NewCommandMap() *ComMap {
	cmdm := make(ComMap)
	return &cmdm
}

/*
 * fn = func(args []string)
 * returns an error if the command already was in the map, else nil
 */
func (c ComMap) AddCommand(cmd, desc string, fn commandFunc) error {
	if c[cmd] == nil {
		c[cmd] = &command{desc, fn}
		return nil
	}
	return nil // error
}

/*
 * returns an error if there was no command to remove, else nil
 */
func (c ComMap) RemoveCommand(cmd string) error {
	if c[cmd] == nil {
		return nil // error
	}
	delete(c, cmd)
	return nil
}

func readLine(reader *bufio.Reader) string {
	data, _, err := reader.ReadLine()
	if err != nil {
		logger.Println("Error reading input. Quitting.")
		return "quit"
	}
	return string(data)
}

func (c ComMap) printHelp() {
	fmt.Println("This is the helptext for the MetaParticle server.")
	fmt.Println(FgYellow + "+----------------")
	fmt.Println("| Avaliable commands:" + Reset)
	
	for name, cmd := range c {
		fmt.Printf(FgYellow + "| " + FgGreen + "%s: " + Reset + "%s.\n", name, cmd.desc)
	}
	
	fmt.Println(FgYellow + "+----------------" + Reset)
	fmt.Println()	
}

func PrintStatus(ok int, s string) {
	switch {
		case ok > 0: fmt.Print(FgGreen)
		case ok == 0: fmt.Print(FgYellow)
		default: fmt.Print(FgRed)
	}
	fmt.Println(s + Reset)
}

func AskYN(question string) bool {
	fmt.Print(question + " (Y/n) ")
	reader := bufio.NewReader(os.Stdin)
	answer, _, err := reader.ReadLine()
	if err != nil {
		return false
	}
	return strings.ContainsAny(strings.ToLower(string(answer)), "yes")
}
