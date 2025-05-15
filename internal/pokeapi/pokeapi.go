package pokeapi

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/horsedevours/pokedex/internal/pokecache"
)

type locationArea struct {
	Count    int
	Next     string
	Previous string
	Results  []struct {
		Name string
		Url  string
	}
}

type areaPokemon struct {
	Encounters []struct {
		Pokemon struct {
			Name string
			Url  string
		}
	} `json:"pokemon_encounters"`
}

type Pokemon struct {
	Name           string
	Height         int
	Weight         int
	Stats          []PokemonStat
	Types          []PokemonType
	BaseExperience int `json:"base_experience"`
}

type PokemonStat struct {
	Stat struct {
		Name string
	}
	BaseStat int `json:"base_stat"`
}

type PokemonType struct {
	Type struct {
		Name string
	}
}

var cache = pokecache.NewCache(5 * time.Minute)

func GetLocationAreas(url string) (locationArea, error) {
	locs := locationArea{}
	if data, ok := cache.Get(url); ok {
		if err := json.Unmarshal(data, &locs); err != nil {
			return locationArea{}, err
		}
		return locs, nil
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return locationArea{}, err
	}

	client := &http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return locationArea{}, err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return locationArea{}, err
	}

	if err := json.Unmarshal(data, &locs); err != nil {
		return locationArea{}, err
	}

	cache.Add(url, data)
	return locs, nil
}

func GetAreaPokemon(url string) (areaPokemon, error) {
	ap := areaPokemon{}
	if data, ok := cache.Get(url); ok {
		if err := json.Unmarshal(data, &ap); err != nil {
			return areaPokemon{}, err
		}
		return ap, nil
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return areaPokemon{}, err
	}

	client := &http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return areaPokemon{}, err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return areaPokemon{}, err
	}

	if err := json.Unmarshal(data, &ap); err != nil {
		return areaPokemon{}, err
	}

	cache.Add(url, data)
	return ap, nil
}

func GetPokemonData(url string) (Pokemon, error) {
	pkmn := Pokemon{}
	if data, ok := cache.Get(url); ok {
		if err := json.Unmarshal(data, &pkmn); err != nil {
			return Pokemon{}, err
		}
		return pkmn, nil
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return Pokemon{}, err
	}

	client := &http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return Pokemon{}, err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return Pokemon{}, err
	}

	if err := json.Unmarshal(data, &pkmn); err != nil {
		return Pokemon{}, err
	}

	cache.Add(url, data)
	return pkmn, nil
}
