package main

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strings"
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
	"pokedex": {
		name:        "pokedex",
		description: "",
		callback:    commandPokedex,
	},
	// help is put into by main to break "Initialization cycle"
}

// protect against names that could corrupt url if appended
// certainly no ../../  &=/#?%
func dirtyNameRune(r rune) bool {
	if r >= '0' && r <= '9' {
		return false
	}
	if r >= 'a' && r <= 'z' {
		return false
	}
	if r == '-' {
		return false
	}
	return true
}
func dirtyName(text string) bool {
	return strings.ContainsFunc(text, dirtyNameRune)
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
	if len(args) <= 1 {
		return errors.New("not enough arguments")
	}
	if dirtyName(args[1]) {
		fmt.Println("not a proper area name")
		return errors.New("dirty name")
	}
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
	if len(args) <= 1 {
		return errors.New("not enough arguments")
	}
	if dirtyName(args[1]) {
		fmt.Println("not a proper pokemon name")
		return errors.New("dirty name")
	}
	poke, err := getPokemonResult(args[1])
	if err != nil {
		fmt.Printf("could not find pokemon named \"%s\" to catch\n", args[1])
		return err
	}
	fmt.Printf("Throwing a Pokeball at %s...\n", args[1])
	lowXpMaxPower := poke.Base_experience + 5 + poke.Base_experience/144
	normMaxPower := cmdGLOBS.xp + 40 + rand.Intn(max(50, cmdGLOBS.xp, poke.Base_experience/2))
	power := rand.Intn(max(lowXpMaxPower, normMaxPower))
	var closeXpGain int
	if (power+75 > poke.Base_experience) && (power-75 < poke.Base_experience) {
		closeXpGain = rand.Intn(power/3+15+poke.Base_experience/7) + 1
	}
	if power >= poke.Base_experience {
		xpGain := max(10, (poke.Base_experience-cmdGLOBS.xp)/8, closeXpGain)
		cmdGLOBS.xp += xpGain
		fmt.Printf("%s was caught!\n", args[1])
		cmdGLOBS.caught[args[1]] = poke
	} else {
		xpGain := max(1, (poke.Base_experience-cmdGLOBS.xp)/13, closeXpGain)
		fmt.Printf("%s escaped!\n", args[1])
		cmdGLOBS.xp += xpGain
	}
	//fmt.Printf("base %d\n", poke.Base_experience)
	//fmt.Printf("stats %v\n", poke.Stats)
	//fmt.Printf("typess %v\n", poke.Types)
	return nil
}

func commandInspect(args []string) error {
	if len(args) <= 1 {
		return errors.New("not enough arguments")
	}
	if dirtyName(args[1]) {
		fmt.Println("not a proper pokemon name")
		return errors.New("dirty name")
	}
	poke, ok := cmdGLOBS.caught[args[1]]
	if ok {
		fmt.Printf("Name: %s\n", poke.Name)
		fmt.Printf("Height: %d\n", poke.Height)
		fmt.Printf("Weight: %d\n", poke.Weight)
		fmt.Println("Stats:")
		for _, stat := range poke.Stats {
			fmt.Printf(" - %s: %d\n", stat.Stat.Name, stat.Base_stat)
		}
		fmt.Println("Types:")
		for _, pType := range poke.Types {
			fmt.Printf(" - %s\n", pType.Type.Name)
		}

	} else {
		fmt.Println("you have not caught that pokemon")
	}
	return nil
}

func commandPokedex(args []string) error {
	fmt.Println("Your Pokedex:")
	for _, poke := range cmdGLOBS.caught {
		fmt.Printf(" - %s\n", poke.Name)
	}
	return nil
}
