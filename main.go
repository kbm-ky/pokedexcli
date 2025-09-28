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
var pokedex = map[string]pokeapi.Pokemon{}

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
		"catch": {
			"catch",
			"Attempts to catch the Pokemon given its name",
			commandCatch,
		},
		"inspect": {
			"inspect",
			"Inspects the details of a caught Pokemon",
			commandInspect,
		},
		"pokedex": {
			"pokedex",
			"Prints a list of Pokemon names in your Pokedex",
			commandPokedex,
		},
	}
}

func main() {

	cmdConfig := cmdConfig{
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
		fmt.Printf("error: %v\n", err)
		return nil
	}

	for _, encounter := range area.PokemonEncounters {
		fmt.Printf(" - %v\n", encounter.Pokemon.Name)
	}
	return nil
}

func commandCatch(config *cmdConfig, args []string) error {
	if len(args) == 0 {
		fmt.Printf("Catch expects an argument, not given\n")
		return nil
	}

	target := args[0]
	fmt.Printf("Throwing a Pokeball at %s...\n", target)
	pokemon, err := pokeapi.GetPokemon(target)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return nil
	}

	user_exp, target_exp := 100, pokemon.BaseExperience
	caught := attemptCatch(user_exp, target_exp)
	if caught {
		fmt.Printf("%s was caught!\n", target)
		pokedex[target] = pokemon
	} else {
		fmt.Printf("%s escaped!\n", target)
	}

	return nil
}

func commandInspect(config *cmdConfig, args []string) error {
	if len(args) == 0 {
		fmt.Printf("Inspect expects an argument, none given")
		return nil
	}

	name := args[0]
	pokemon, ok := pokedex[name]
	if !ok {
		fmt.Printf("you have not caught that pokemon\n")
		return nil
	}

	fmt.Printf("%v\n", pokemon)
	fmt.Printf("Name: %s\n", name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)
	fmt.Printf("Stats:\n")
	for i := range pokemon.Stats {
		stat := pokemon.Stats[i]
		fmt.Printf("  -%s: %d\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Printf("Types:\n")
	for i := range pokemon.Types {
		typ := pokemon.Types[i]
		fmt.Printf("  - %s\n", typ.Type.Name)
	}

	return nil
}

func commandPokedex(config *cmdConfig, args []string) error {
	if len(pokedex) == 0 {
		fmt.Printf("Your Pokedex is empty, go catch some!\n")
		return nil
	}

	fmt.Printf("Your Pokedex:\n")
	for _, k := range slices.Sorted(maps.Keys(pokedex)) {
		fmt.Printf(" - %s\n", k)
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
