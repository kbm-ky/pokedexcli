package pokeapi

type NamedAPIResourceList[T any] struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []T    `json:"results"`
}

type Location struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Region      any    `json:"region"`
	Names       []any  `json:"names"`
	GameIndices []any  `json:"game_indices"`
	Areas       []any  `json:"areas"`
}

type LocationArea struct {
	Id                   int                `json:"id"`
	Name                 string             `json:"name"`
	GameIndex            int                `json:"game_index"`
	EncounterMethodRates []any              `json:"encounter_method_rates"`
	Location             Location           `json:"location"`
	Names                []any              `json:"names"`
	PokemonEncounters    []PokemonEncounter `json:"pokemon_encounters"`
}

type PokemonEncounter struct {
	Pokemon        Pokemon `json:"pokemon"`
	VersionDetails []any   `json:"version_details"`
}

type Pokemon struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}
