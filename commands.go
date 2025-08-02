package main

import (
	"fmt"
	"os"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

/*
type cliConfigUrls struct {
	previous string
	next     string
}
type cliConfigEmpty struct {
}
type cliConfig interface {
	//cliConfigUrls | cliConfigEmpty

} */

var cmdList = map[string]cliCommand{
	"exit": {
		name:        "exit",
		description: "Exit the Pokedex",
		callback:    commandExit,
	},
	"map": {
		name:        "map",
		description: "displays the names of 20 location areas in the Pokemon world",
		callback:    commandMap,
	},
	"mapb": {
		name:        "mapb",
		description: "displays the previous names of 20 location areas in the Pokemon world",
		callback:    commandMapb,
	},
	// help is put into by main to break "Initialization cycle"
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp() error {
	fmt.Print("Welcome to the Pokedex!\nUsage: \n\n")
	for _, cmdItm := range cmdList {
		fmt.Printf("%s: %s\n", cmdItm.name, cmdItm.description)
	}

	// help: Displays a help message
	// exit: Exit the Pokedex")
	return nil
}

func commandMap() error {
	getNamedResourceResult(&webGLOBS.localAreasList, webGLOBS.localAreasList.currIndx)
	fmt.Printf("%s\n",
		webGLOBS.localAreasList.linkedList[webGLOBS.localAreasList.currIndx].String())
	if webGLOBS.localAreasList.linkedList[webGLOBS.localAreasList.currIndx].Next != nil {
		webGLOBS.localAreasList.currIndx++
	}
	return nil
}

func commandMapb() error {
	if webGLOBS.localAreasList.currIndx > 0 {
		webGLOBS.localAreasList.currIndx--
	} else {
		fmt.Println("you're on the first page")
		return nil
	}
	fmt.Printf("%s\n",
		(webGLOBS.localAreasList.linkedList[webGLOBS.localAreasList.currIndx]).String())
	return nil
}
