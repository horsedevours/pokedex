package main

import (
	"bufio"
	"fmt"
	"math/rand"
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

var pokedex map[string]pokeapi.Pokemon = map[string]pokeapi.Pokemon{}

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
		"catch": {
			name:        "catch",
			description: "Accepts pokemon name as argument; throws a Pokeball at the specified Pokemon",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Accepts pokemon name as argument; displays details if specified Pokemon has been caught",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Lists all of the user's caught pokemon",
			callback:    commandPokedex,
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

func commandCatch(cfg *Config, pokemon string) error {
	fmt.Printf("Throwing a Pokeball at %s...\n", pokemon)

	pkmn, err := pokeapi.GetPokemonData(baseUrl + "pokemon/" + pokemon)
	if err != nil {
		return err
	}

	battle := rand.Intn(pkmn.BaseExperience)
	if battle < 40 {
		fmt.Printf("%s was caught!\n", pkmn.Name)
		pokedex[pkmn.Name] = pkmn
	} else {
		fmt.Printf("%s escaped!\n", pkmn.Name)
	}

	return nil
}

func commandInspect(cfg *Config, pokemon string) error {
	pkmn, ok := pokedex[pokemon]
	if !ok {
		fmt.Printf("You ain't caught no %s!\n", pokemon)
		return nil
	}

	fmt.Printf("Name: %s\n", pkmn.Name)
	fmt.Printf("Height: %d\n", pkmn.Height)
	fmt.Printf("Weight: %d\n", pkmn.Weight)
	fmt.Println("Stats:")
	for _, stat := range pkmn.Stats {
		fmt.Printf("  -%s: %d\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Println("Types:")
	for _, typ := range pkmn.Types {
		fmt.Printf("  -%s\n", typ.Type.Name)
	}

	return nil
}

func commandPokedex(cfg *Config, dummy string) error {
	if len(pokedex) == 0 {
		fmt.Println("You ain't caught nothin'")
		return nil
	}

	fmt.Println("Your pokedex:")
	for p, _ := range pokedex {
		fmt.Printf(" - %s\n", p)
	}
	return nil
}
