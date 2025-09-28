package pokeapi

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func GetLocationArea(page string) (NamedAPIResourceList[LocationArea], error) {
	var list NamedAPIResourceList[LocationArea]
	res, err := http.Get(page)
	if err != nil {
		return list, err
	}
	defer res.Body.Close()

	if res.StatusCode > 299 {
		return list, fmt.Errorf("bad response code: %d", res.StatusCode)
	}

	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&list); err != nil {
		return NamedAPIResourceList[LocationArea]{}, err
	}

	return list, nil
}
