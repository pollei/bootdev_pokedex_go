package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	//fmt.Println("Hello, World!")
	inpScan := bufio.NewScanner(os.Stdin)
	// help put  in by hand to break "Initialization cycle"
	cmdList["help"] = cliCommand{
		name: "help", description: "Displays a help message", callback: commandHelp}
	webGLOBS.localAreasList.baseUrl = "https://pokeapi.co/api/v2/location-area/"
	getNamedResourceResult(&webGLOBS.localAreasList, 0)

	for {
		fmt.Print("Pokedex > ")
		inpScan.Scan() // Reads the next line
		lineStr := inpScan.Text()
		cleanLine := cleanInput(lineStr)
		//fmt.Printf("Your command was: %s\n", cleanLine[0])
		//if "exit" == cleanLine[0] { break }
		cmd, ok := cmdList[cleanLine[0]]
		if ok {
			cmd.callback()
		} else {
			fmt.Println("Unknown command")
		}
	}
}

func cleanInput(text string) []string {
	trimStr := strings.TrimSpace(text)
	lowStr := strings.ToLower(trimStr)
	return strings.Fields(lowStr)
	//return []string{}
}
