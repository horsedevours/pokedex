package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/horsedevours/pokedex/internal/pokeapi"
)

const baseUrl string = "https://pokeapi.co/api/v2/"

type cliCommand struct {
	name        string
	description string
	callback    func(*Config) error
}

type Config struct {
	Next     string
	Previous string
}

var registry map[string]cliCommand

func main() {
	registry = map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "Displays next 20 locations",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays previous 20 locations",
			callback:    commandMapb,
		},
	}

	cfg := Config{}

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := scanner.Text()
		inputSlice := cleanInput(input)

		if cmd, ok := registry[inputSlice[0]]; ok {
			cmd.callback(&cfg)
			continue
		}
		fmt.Println("Unknown command")
	}
}

func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}

func commandHelp(cfg *Config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	for _, cmd := range registry {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return fmt.Errorf("Failed to help for some reason...")
}

func commandExit(cfg *Config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return fmt.Errorf("Failed to exit for some reason...")
}

func commandMap(cfg *Config) error {
	fullUrl := ""
	if cfg.Next != "" {
		fullUrl = cfg.Next
	} else {
		fullUrl = baseUrl + "location-area"
	}

	locs, err := pokeapi.GetLocationAreas(fullUrl)
	if err != nil {
		return err
	}

	cfg.Next = locs.Next
	cfg.Previous = locs.Previous

	for _, loc := range locs.Results {
		fmt.Println(loc.Name)
	}

	return nil
}

func commandMapb(cfg *Config) error {
	if cfg.Previous == "" {
		fmt.Println("you're on the first page")
	}

	locs, err := pokeapi.GetLocationAreas(cfg.Previous)
	if err != nil {
		return err
	}

	cfg.Next = locs.Next
	cfg.Previous = locs.Previous

	for _, loc := range locs.Results {
		fmt.Println(loc.Name)
	}

	return nil
}
