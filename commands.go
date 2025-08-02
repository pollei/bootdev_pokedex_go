package main

import (
	"fmt"
	"os"
)

type cliCommand struct {
	name        string
	description string
	callback    func([]string) error
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
	"explore": {
		name:        "explore",
		description: "lists Pokemon's names in a location",
		callback:    commandExplore,
	},
	"catch": {
		name:        "catch",
		description: "",
		callback:    commandCatch,
	},
	"inspect": {
		name:        "inspect",
		description: "",
		callback:    commandInspect,
	},
	// help is put into by main to break "Initialization cycle"
}

func commandExit([]string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp([]string) error {
	fmt.Print("Welcome to the Pokedex!\nUsage: \n\n")
	for _, cmdItm := range cmdList {
		fmt.Printf("%s: %s\n", cmdItm.name, cmdItm.description)
	}

	// help: Displays a help message
	// exit: Exit the Pokedex")
	return nil
}

func commandMap([]string) error {
	getNamedResourceResult(&webGLOBS.localAreasList, webGLOBS.localAreasList.currIndx)
	fmt.Printf("%s\n",
		webGLOBS.localAreasList.linkedList[webGLOBS.localAreasList.currIndx].String())
	if webGLOBS.localAreasList.linkedList[webGLOBS.localAreasList.currIndx].Next != nil {
		webGLOBS.localAreasList.currIndx++
	}
	return nil
}

func commandMapb([]string) error {
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

func commandExplore(args []string) error {
	locArea, err := getExploreResult(args[1])
	if err != nil {
		fmt.Println("could not explore area")
		return err
	}
	pokeLst := locArea.Pokemon_encounters
	fmt.Printf("Exploring p%s..\nFound Pokemon:\n", args[1])
	//fmt.Printf("encounter %v %d", locArea.Pokemon_encounters, len(locArea.Pokemon_encounters))
	for _, poke := range pokeLst {
		fmt.Println(" - " + poke.Pokemon.Name)
	}
	return nil
}

func commandCatch([]string) error {
	return nil
}

func commandInspect([]string) error {
	return nil
}
