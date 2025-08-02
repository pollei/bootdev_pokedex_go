package main

import (
	"fmt"
	"math/rand"
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

}
*/
var cmdGLOBS = struct {
	xp     int
	caught map[string]Pokemon
}{}

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
		description: "try catching a pokemon by name",
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

func commandCatch(args []string) error {
	poke, err := getPokemonResult(args[1])
	if err != nil {
		fmt.Println("could not catch pokemon")
		return err
	}
	fmt.Printf("Throwing a Pokeball at %s...\n", args[1])
	power := rand.Intn(200 + cmdGLOBS.xp)
	if power >= poke.Base_experience {
		xpGain := max(10, (poke.Base_experience-cmdGLOBS.xp)/10)
		cmdGLOBS.xp += xpGain
		fmt.Printf("%s was caught!\n", args[1])
		cmdGLOBS.caught[args[1]] = poke
	} else {
		fmt.Printf("%s escaped!\n", args[1])
		cmdGLOBS.xp++
	}
	//fmt.Printf("base %d\n", poke.Base_experience)
	//fmt.Printf("stats %v\n", poke.Stats)
	//fmt.Printf("typess %v\n", poke.Types)
	return nil
}

func commandInspect([]string) error {
	return nil
}
