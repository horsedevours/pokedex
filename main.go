package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
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

type locationArea struct {
	Count    int
	Next     string
	Previous string
	Results  []struct {
		Name string
		Url  string
	}
}

func commandMap(cfg *Config) error {
	fullUrl := ""
	if cfg.Next != "" {
		fullUrl = cfg.Next
	} else {
		fullUrl = baseUrl + "location-area"
	}
	req, err := http.NewRequest("GET", fullUrl, nil)
	if err != nil {
		return err
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	locs := locationArea{}
	if err := decoder.Decode(&locs); err != nil {
		return err
	}

	fmt.Println(locs.Next)

	for _, loc := range locs.Results {
		fmt.Println(loc.Name)
	}
	cfg.Next = locs.Next
	cfg.Previous = locs.Previous

	return nil
}

func commandMapb(cfg *Config) error {
	if cfg.Previous == "" {
		fmt.Println("you're on the first page")
	}
	req, err := http.NewRequest("GET", cfg.Previous, nil)
	if err != nil {
		return err
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	locs := locationArea{}
	if err := decoder.Decode(&locs); err != nil {
		return err
	}

	fmt.Println(locs.Next)

	for _, loc := range locs.Results {
		fmt.Println(loc.Name)
	}
	cfg.Next = locs.Next
	cfg.Previous = locs.Previous

	return nil
}
