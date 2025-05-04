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

var cache = pokecache.NewCache(5 * time.Second)

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
