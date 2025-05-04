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
	callback    func(*Config, string) error
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
		"explore": {
			name:        "explore",
			description: "Accepts area name as argument; displays Pokemon inhabiting the specified area",
			callback:    commandExplore,
		},
	}

	cfg := Config{}

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := scanner.Text()
		inputSlice := cleanInput(input)
		if len(inputSlice) == 0 {
			continue
		}

		extraArg := ""
		if len(inputSlice) > 1 {
			extraArg = inputSlice[1]
		}

		if cmd, ok := registry[inputSlice[0]]; ok {
			cmd.callback(&cfg, extraArg)
			continue
		}
		fmt.Println("Unknown command")
	}
}

func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}

func commandHelp(cfg *Config, dummy string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	for _, cmd := range registry {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return fmt.Errorf("Failed to help for some reason...")
}

func commandExit(cfg *Config, dummy string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return fmt.Errorf("Failed to exit for some reason...")
}

func commandMap(cfg *Config, dummy string) error {
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

func commandMapb(cfg *Config, dummy string) error {
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

func commandExplore(cfg *Config, area string) error {
	fmt.Printf("Exploring %s area...\n", area)
	areaPokemon, err := pokeapi.GetAreaPokemon(baseUrl + "location-area/" + area)
	if err != nil {
		return err
	}

	fmt.Println("Found Pokemon:")
	for _, poke := range areaPokemon.Encounters {
		fmt.Printf("- %s\n", poke.Pokemon.Name)
	}

	return nil
}
