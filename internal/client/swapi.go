package client

import (
	"encoding/json"
	"github.com/bernardoms/StarWarsPlanetAPI-GO/internal/logger"
	"net/http"
)

type SwapiClient struct {
	Endpoint string
	log      *logger.Logger
}

func NewSwapiClient(endpoint string, log *logger.Logger) *SwapiClient {
	s := new(SwapiClient)
	s.Endpoint = endpoint
	s.log = log
	return s
}

func (s SwapiClient) GetPlanetByName(name string) (*SwapiPlanet, error) {
	resp, err := http.Get(s.Endpoint + "planets?search=" + name)

	if err != nil {
		s.log.LogWithFields(nil, "error", map[string]interface{}{"err": "error get planet from swapi client"}, err.Error())
		return nil, err
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	var swapi *SwapiPlanet

	err = json.NewDecoder(resp.Body).Decode(&swapi)

	if err != nil {
		s.log.LogWithFields(nil, "error", map[string]interface{}{"err": "error unmarshalling response"}, err.Error())
	}

	_ = resp.Body.Close()

	return swapi, err
}
