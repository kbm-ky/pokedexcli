package main

import (
	"bufio"
	"fmt"
	"maps"
	"os"
	"slices"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
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
	}
}

func main() {

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
				err := cmdStruct.callback()
				if err != nil {
					fmt.Printf("error: %v\n", err)
					os.Exit(1)
				}
			}
		}
	}
}

func commandExit() error {
	fmt.Printf("Closing the Pokedex... Goodbye!\n")
	os.Exit(0)
	return nil
}

func commandHelp() error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Printf("Usage:\n\n")
	for _, k := range slices.Sorted(maps.Keys(commands)) {
		fmt.Printf("%s: %s\n", k, commands[k].description)
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
