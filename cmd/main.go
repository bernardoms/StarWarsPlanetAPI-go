package main

import (
	"fmt"
	config2 "github.com/bernardoms/StarWarsPlanetAPI-GO/config"
	"github.com/bernardoms/StarWarsPlanetAPI-GO/internal/client"
	"github.com/bernardoms/StarWarsPlanetAPI-GO/internal/handler"
	"github.com/bernardoms/StarWarsPlanetAPI-GO/internal/repository"
	"github.com/gorilla/mux"
	newrelic "github.com/newrelic/go-agent"
	"github.com/newrelic/go-agent/_integrations/nrgorilla/v1"
	"log"
	"net/http"
	"os"
)

func main() {

	app, errNewRelic := newrelic.NewApplication(
		newrelic.NewConfig(os.Getenv("NEWRELIC_APP"), os.Getenv("NEWRELIC_LICENSE")),
	)

	if errNewRelic != nil {
		log.Print("Error starting new relic agent")
	}

	config := config2.NewMongoConfig()

	mongo := repository.NewSession(*config)

	swapiClient := client.NewSwapiClient("https://swapi.dev/api/")

	r := mux.NewRouter()

	planetHandler := handler.NewPlanetHandler(mongo, swapiClient)

	nrgorilla.InstrumentRoutes(r, app)

	r.HandleFunc("/v1/planets", planetHandler.SavePlanet).Methods("POST")
	r.HandleFunc("/v1/planets", planetHandler.GetPlanets).Methods("GET")
	r.HandleFunc("/v1/planets/{planetId}", planetHandler.GetPlanetById).Methods("GET")
	r.HandleFunc("/v1/planets/{planetId}", planetHandler.RemovePlanetById).Methods("DELETE")

	fmt.Printf("running server on %d", 8080)

	err := http.ListenAndServe(":8080", r)

	if err != nil {
		fmt.Printf("error to open port %s with error %s", "8080", err)
	}
}
