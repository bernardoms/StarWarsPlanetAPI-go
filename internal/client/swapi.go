package client

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type SwapiClient struct {
	Endpoint string
}


func NewSwapiClient(endpoint string) *SwapiClient {
	s := new(SwapiClient)
	s.Endpoint = endpoint
	return s
}

func  (s SwapiClient) GetPlanetByName(name string) (*SwapiPlanet, error) {
	resp, err := http.Get(s.Endpoint + "planets?search=" + name)

	if err != nil {
		fmt.Print(err)
		return nil, err
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	var swapi *SwapiPlanet

	err = json.NewDecoder(resp.Body).Decode(&swapi)

	if err != nil {
		fmt.Print(err)
	}

	_ = resp.Body.Close()

	return swapi, err
}