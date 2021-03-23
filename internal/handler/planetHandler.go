package handler

import (
	"encoding/json"
	"github.com/bernardoms/StarWarsPlanetAPI-GO/internal/client"
	"github.com/bernardoms/StarWarsPlanetAPI-GO/internal/repository"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
)

type PlanetRequest struct {
	Name               string             `bson:"name"`
	Weather            string             `bson:"weather"`
	Land               string             `bson:"land"`
}

type ResponseError struct {
	Description string `json:"description"`
}

type PlanetHandler struct {
	swapiClient client.SwapiClientInterface
	repository  repository.PlanetRepositoryInterface
}

var decoder = schema.NewDecoder()

func NewPlanetHandler(mongo repository.PlanetRepositoryInterface,
	swapiClient client.SwapiClientInterface) *PlanetHandler{

	planetHandler := new(PlanetHandler)

	planetHandler.swapiClient = swapiClient
	planetHandler.repository = mongo

	return planetHandler
}

func (p *PlanetHandler) GetPlanets(w http.ResponseWriter, r *http.Request) {
	filter := new(repository.Filter)
	err := decoder.Decode(filter, r.URL.Query())

	planets, err := p.repository.FindAll(*filter)

	if err != nil {
		respondWithJson(w, http.StatusInternalServerError, ResponseError{Description: err.Error()})
		return
	}

	respondWithJson(w, http.StatusOK, planets)
}

func (p *PlanetHandler) GetPlanetById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	objectId, err := primitive.ObjectIDFromHex(vars["planetId"])
	if err != nil {
		respondWithJson(w, http.StatusBadRequest, ResponseError{Description: "planet id is not a valid id"})
		return
	}
	
	foundPlanet, err := p.repository.FindById(objectId)

	if err != nil {
		//f := map[string]interface{}{"msg": err}
		//u.Logger.LogWithFields(r, "error", f)
		respondWithJson(w, http.StatusInternalServerError, ResponseError{Description: err.Error()})
		return
	}

	respondWithJson(w, http.StatusOK, foundPlanet)
}

func (p *PlanetHandler) RemovePlanetById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	objectId, err := primitive.ObjectIDFromHex(vars["planetId"])
	if err != nil {
		respondWithJson(w, http.StatusBadRequest, ResponseError{Description: "planet id is not a valid id"})
		return
	}

	err = p.repository.Delete(objectId)

	if err != nil {
		//f := map[string]interface{}{"msg": err}
		//u.Logger.LogWithFields(r, "error", f)
		respondWithJson(w, http.StatusInternalServerError, ResponseError{Description: err.Error()})
		return
	}

	respondWithEmpty(w, http.StatusNoContent, "")
}

func (p *PlanetHandler) SavePlanet(w http.ResponseWriter, r *http.Request) {

	var planetRequest PlanetRequest

	err := json.NewDecoder(r.Body).Decode(&planetRequest)

	if err != nil {
		log.Println("error unmarshalling the request body", err)
	}

	planets, err := p.swapiClient.GetPlanetByName(planetRequest.Name)

	if err != nil {
		log.Println("error getting planets from swapi api", err)
		respondWithJson(w, http.StatusInternalServerError, ResponseError{Description: err.Error()})
		return
	}

	if planets == nil {
		log.Println("planet not found on star wars world")

		respondWithEmpty(w, http.StatusNotFound, "")

		return
	}

	planet := new(repository.Planet)

	planet.Id = primitive.NewObjectID()
	planet.Name = planetRequest.Name
	planet.Land = planetRequest.Land
	planet.Weather = planetRequest.Weather
	planet.AppearanceQuantity = len(planets.Results[0].Films)

	savedPlanet, err := p.repository.Save(planet)

	if err != nil {
		log.Println("error saving planet", err)
		respondWithJson(w, http.StatusInternalServerError, ResponseError{Description: err.Error()})
		return
	}

	respondWithEmpty(w, http.StatusCreated, "v1/planets/" + savedPlanet.Id.Hex())
}

func respondWithEmpty(w http.ResponseWriter, code int, location string) {
	w.Header().Set("Content-Type", "application/json")
	if location != "" {
		w.Header().Set("Location", location)
	}
	w.WriteHeader(code)
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(code)
	_, _ = w.Write(response)
}