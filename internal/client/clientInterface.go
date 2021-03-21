package client

type SwapiPlanet struct {
	Results [] Results `json:"results"`
}

type Results struct {
	Name           string   `json:"name"`
	Diameter       string   `json:"diameter"`
	Gravity        string   `json:"gravity"`
	Population     string   `json:"population"`
	Climate        string   `json:"climate"`
	Terrain        string   `json:"terrain"`
	Created        string   `json:"created"`
	Edited         string   `json:"edited"`
	Url            string   `json:"url"`
	Residents      []string `json:"residents"`
	Films          []string `json:"films"`
}

type SwapiClientInterface interface {
	GetPlanetByName(name string) (*SwapiPlanet, error)
}