package main

import (
	"fmt"
	"flag"
	"github.com/MetaParticle/metaparticle/cmdln"
	. "github.com/MetaParticle/metaparticle/storage"
	. "github.com/MetaParticle/metaparticle/storage/mgodb"
)

type Person struct {
	Name string
	Phone string
}

func (p Person) String() string {
	return "{Name: "+p.Name+"\nPhone: "+p.Phone+"}"
}

func main() {
	flag.Parse()

	var db DB
	var c Set
	db = ConnectMgo(flag.Arg(0), "")
	if db == nil {
		fmt.Println("Can't connect to DB.")
		return
	}
	defer db.Close()

	exit := make(chan bool)
	cmds := cmdln.New(
		cmdln.BoldOn + "test" + cmdln.BoldOff + " > ",
		cmdln.NewCommandMap(),
		exit)
	cmds.PopulateDefaults()
	cmds.Commands.AddCommand("set", "selects a set",
	func(args []string) {
		if len(args) > 0 {
			c = db.Get(args[0])
			if c == nil {
				fmt.Println("Can't retrieve Set.")
				return
			}
		}
		fmt.Println("Set:", c)
	})

	cmds.Commands.AddCommand("get", "gets all matches of the input",
	func(args []string) {
		if c == nil {
			fmt.Println("No set chosen.")
			return
		}
		if len(args) > 1 {
			results := make([]Person, 5)
			_, err := c.Get(&results, args...)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(results)
			}
		} else {
			fmt.Println("Too few arguments.")
		}
	})

	cmds.Commands.AddCommand("getone", "gets a singe match of the input",
	func(args []string) {
		if c == nil {
			fmt.Println("No set chosen.")
			return
		}
		if len(args) > 1 {
			result := Person{}
			err := c.GetOne(&result, args...)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(result)
			}
		} else {
			fmt.Println("Too few arguments.")
		}
	})

	cmds.Commands.AddCommand("fill", "fills the example {Name:\"Ale\"}",
	func(args []string) {
		if c == nil {
			fmt.Println("No set chosen.")
			return
		}
		result := Person{}
		result.Name = "Ale"
		err := c.Fill(&result)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(result)
		}
	})

	cmds.Commands.AddCommand("defaults", "adds test entries",
	func(args []string) {
		if c == nil {
			fmt.Println("No set chosen.")
			return
		}
		err := c.Insert(&Person{"Nisse", "0906516345"}, &Person{"Ale", "0734516345"})
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Insert successful.")
		}
		result := Person{}
		c.GetOne(result, "Name", "Ale")
		fmt.Println(&result)
	})

	cmds.Commands.AddCommand("insert", "inserts a new Person",
	func(args []string) {
		if c == nil {
			fmt.Println("No set chosen.")
			return
		}
		if (len(args) < 2) {
			fmt.Println("Too few arguments.")
			return
		}
		err := c.Insert(&Person{args[0], args[1]})
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Insert successful.")
		}
	})

	go cmds.ListenAndServe()

	/*
	
	
	var c Set
	c = db.Get("people")
	if c == nil {
		fmt.Println("Can't retrieve set.")
	}
	*/
	//result := &Person{}
	//results := make([]Person, 5)
	/*
	data := &Person{Name: "Ale"}
	
	err := c.Fill(data, result)
	*/
	/*
	err := c.GetOne(result, flag.Args()...)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Name:", result.Name)
	fmt.Println("Phone:", result.Phone)
	
	n, err := c.Get(&results, flag.Args()...)
	if err != nil {
		fmt.Println(err)
		return
	}
	
	fmt.Println(n)
	fmt.Println(results)
	*/
	select {
		case <- exit: 
	}
}
