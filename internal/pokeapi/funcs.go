package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/kbm-ky/pokedexcli/internal/pokecache"
)

var cache *pokecache.Cache

func init() {
	cache = pokecache.NewCache(5 * time.Second)
}

func Get(url string) ([]byte, error) {
	data, ok := cache.Get(url)
	if !ok {
		res, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()

		if res.StatusCode > 299 {
			return nil, fmt.Errorf("bad response code: %d", res.StatusCode)
		}

		data, err = io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		cache.Add(url, data)
	}
	return data, nil
}

func GetLocationArea(url string) (NamedAPIResourceList[LocationArea], error) {
	var list NamedAPIResourceList[LocationArea]

	data, err := Get(url)
	if err != nil {
		return NamedAPIResourceList[LocationArea]{}, err
	}

	err = json.Unmarshal(data, &list)
	if err != nil {
		return NamedAPIResourceList[LocationArea]{}, err
	}

	return list, nil
}
