package client

import (
	client2 "github.com/bernardoms/StarWarsPlanetAPI-GO/internal/client"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestShouldReturnPlanetWithSuccessFromGet(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/planets" {
				w.Header().Add("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{
    "count": 1, 
    "next": null, 
    "previous": null, 
    "results": [
        {
            "name": "Alderaan", 
            "rotation_period": "24", 
            "orbital_period": "364", 
            "diameter": "12500", 
            "climate": "temperate", 
            "gravity": "1 standard", 
            "terrain": "grasslands, mountains", 
            "surface_water": "40", 
            "population": "2000000000", 
            "residents": [
                "http://swapi.dev/api/people/5/", 
                "http://swapi.dev/api/people/68/", 
                "http://swapi.dev/api/people/81/"
            ], 
            "films": [
                "http://swapi.dev/api/films/1/", 
                "http://swapi.dev/api/films/6/"
            ], 
            "created": "2014-12-10T11:35:48.479000Z", 
            "edited": "2014-12-20T20:58:18.420000Z", 
            "url": "http://swapi.dev/api/planets/2/"
        }
    ]
}`))
			}
		}),
	)
	defer ts.Close()
	client := client2.NewSwapiClient(ts.URL + "/")
	planet, err := client.GetPlanetByName("Aldebaran")

	assert.Equal(t, nil, err)
	assert.Equal(t, "Alderaan", planet.Results[0].Name)
	assert.Equal(t, "12500", planet.Results[0].Diameter)
	assert.Equal(t, "1 standard", planet.Results[0].Gravity)
	assert.Equal(t, "2000000000", planet.Results[0].Population)
	assert.Equal(t, "temperate", planet.Results[0].Climate)
	assert.Equal(t, "grasslands, mountains", planet.Results[0].Terrain)
	assert.Equal(t, 2, len(planet.Results[0].Films))
}

func TestShouldReturnPlanetNilWhenApiReturnsNotFound(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/planets" {
				w.Header().Set("Content-Length", "1")
				w.WriteHeader(http.StatusNotFound)
			}
		}))
	client := client2.NewSwapiClient(ts.URL + "/")

	planet, err := client.GetPlanetByName("Aldebaran")

	assert.NoError(t, err)
	assert.Nil(t, planet)
}
