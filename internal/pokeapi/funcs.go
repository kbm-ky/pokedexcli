package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/kbm-ky/pokedexcli/internal/pokecache"
)

var cache *pokecache.Cache

func init() {
	cache = pokecache.NewCache(5 * time.Second)
}

func GetLocationArea(url string) (NamedAPIResourceList[LocationArea], error) {
	var list NamedAPIResourceList[LocationArea]

	data, ok := cache.Get(url)
	if !ok {
		log.Printf("NOT using cache\n")
		res, err := http.Get(url)
		if err != nil {
			return list, err
		}
		defer res.Body.Close()

		if res.StatusCode > 299 {
			return list, fmt.Errorf("bad response code: %d", res.StatusCode)
		}

		data, err = io.ReadAll(res.Body)
		if err != nil {
			return list, err
		}

		cache.Add(url, data)
	} else {
		// log.Printf("Using cache\n")
	}

	err := json.Unmarshal(data, &list)
	if err != nil {
		return NamedAPIResourceList[LocationArea]{}, err
	}

	return list, nil
}
