package main

import (
	"bufio"
	"fmt"
	"maps"
	"os"
	"slices"
	"strings"

	"github.com/kbm-ky/pokedexcli/internal/pokeapi"
)

type cliCommand struct {
	name        string
	description string
	callback    func(config *cmdConfig, args []string) error
}

type cmdConfig struct {
	Next     string
	Previous string
}

var commands map[string]cliCommand

func init() {
	commands = map[string]cliCommand{

		"exit": {
			"exit",
			"Exit the Pokedex",
			commandExit,
		},
		"help": {
			"help",
			"Displays a help message",
			commandHelp,
		},
		"map": {
			"map",
			"Displays next 20 locations in the Pokemon world",
			commandMap,
		},
		"mapb": {
			"mapb",
			"Displays the previous 20 locations in the Pokemon world",
			commandMapBack,
		},
		"explore": {
			"explore",
			"Explores the given location name",
			commandExplore,
		},
	}
}

func main() {

	cmdConfig := cmdConfig{
		// Next: "https://pokeapi.co/api/v2/location-area/",
		Next: pokeapi.LocationAreaEndpoint,
	}
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Printf("Pokedex > ")

		if ok := scanner.Scan(); !ok {
			if err := scanner.Err(); err != nil {
				fmt.Printf("error: %v", err)
				os.Exit(1)
			} else {
				break
			}
		}

		text := scanner.Text()
		cleaned := cleanInput(text)
		if len(cleaned) > 0 {
			command := cleaned[0]
			cmdStruct, ok := commands[command]
			if !ok {
				fmt.Printf("Unknown command\n")
			} else {
				err := cmdStruct.callback(&cmdConfig, cleaned[1:])
				if err != nil {
					fmt.Printf("error: %v\n", err)
					os.Exit(1)
				}
			}
		}
	}
}

func commandExit(_ *cmdConfig, _ []string) error {
	fmt.Printf("Closing the Pokedex... Goodbye!\n")
	os.Exit(0)
	return nil
}

func commandHelp(_ *cmdConfig, _ []string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Printf("Usage:\n\n")

	maxNameLen := 0
	for k := range commands {
		nameLen := len(k)
		if nameLen > maxNameLen {
			maxNameLen = nameLen
		}
	}
	format := "%-*s %s\n"
	for _, k := range slices.Sorted(maps.Keys(commands)) {
		fmt.Printf(format, maxNameLen, k, commands[k].description)
	}
	return nil
}

func commandMap(config *cmdConfig, _ []string) error {
	if config.Next == "" {
		fmt.Printf("You are at the end of the list.\n")
		return nil
	}

	list, err := pokeapi.GetLocationAreaList(config.Next)
	if err != nil {
		return nil
	}

	for _, area := range list.Results {
		fmt.Printf("%s\n", area.Name)
	}

	config.Next = list.Next
	config.Previous = list.Previous

	return nil
}

func commandMapBack(config *cmdConfig, _ []string) error {
	if config.Previous == "" {
		fmt.Printf("You are on the first page.\n")
		return nil
	}

	list, err := pokeapi.GetLocationAreaList(config.Previous)
	if err != nil {
		return nil
	}

	for _, area := range list.Results {
		fmt.Printf("%s\n", area.Name)
	}

	config.Next = list.Next
	config.Previous = list.Previous

	return nil
}

func commandExplore(config *cmdConfig, args []string) error {
	if len(args) == 0 {
		fmt.Printf("Explore expects an argument, not given\n")
		return nil
	}

	location := args[0]
	fmt.Printf("Exploring %s...\n", location)
	area, err := pokeapi.GetLocationArea(location)
	if err != nil {
		return err
	}

	for _, encounter := range area.PokemonEncounters {
		fmt.Printf(" - %v\n", encounter.Pokemon.Name)
	}
	return nil
}

func cleanInput(text string) []string {
	cleaned := []string{}

	for field := range strings.FieldsSeq(text) {
		cleaned = append(cleaned, strings.ToLower(field))
	}

	return cleaned
}
