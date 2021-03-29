package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/bernardoms/StarWarsPlanetAPI-GO/internal/client"
	"github.com/bernardoms/StarWarsPlanetAPI-GO/internal/handler"
	"github.com/bernardoms/StarWarsPlanetAPI-GO/internal/repository"
	"github.com/bernardoms/StarWarsPlanetAPI-GO/test/unit/mock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	mock2 "github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestShouldGetPlanetByIdWithSuccess(t *testing.T) {
	mongoMock := new(mock.MongoMock)
	swapiMock := new(mock.SwapiClientMock)
	mockLogger := new(mock.LoggerMock)

	h := handler.NewPlanetHandler(mongoMock, swapiMock, mockLogger)

	returnedPlanet := repository.Planet{Name: "Aldebaran", Land: "Dry", Weather: "Dry", AppearanceQuantity: 2}

	id, _ := primitive.ObjectIDFromHex("5ea7208049e00ddb76994ede")

	mongoMock.On("FindById", id).Return(&returnedPlanet, nil)

	r, _ := http.NewRequest("GET", "/v1/planets/5ea7208049e00ddb76994ede", nil)
	w := httptest.NewRecorder()

	vars := map[string]string{
		"planetId": "5ea7208049e00ddb76994ede",
	}

	r = mux.SetURLVars(r, vars)

	h.GetPlanetById(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "{\"Name\":\"Aldebaran\",\"Weather\":\"Dry\",\"Land\":\"Dry\",\"AppearanceQuantity\":2}", w.Body.String())
}

func TestShouldGetPlanetByIdReturnBadRequestInvalidId(t *testing.T) {
	mongoMock := new(mock.MongoMock)
	swapiMock := new(mock.SwapiClientMock)
	mockLogger := new(mock.LoggerMock)

	mockLogger.On("LogWithFields", mock2.Anything, mock2.Anything, mock2.Anything, mock2.Anything).Return(nil)

	h := handler.NewPlanetHandler(mongoMock, swapiMock, mockLogger)

	r, _ := http.NewRequest("GET", "/v1/planets/123", nil)
	w := httptest.NewRecorder()

	vars := map[string]string{
		"planetId": "123",
	}

	r = mux.SetURLVars(r, vars)

	h.GetPlanetById(w, r)

	mockLogger.AssertNumberOfCalls(t, "LogWithFields", 1)
	mongoMock.AssertNumberOfCalls(t, "FindById", 0)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "{\"description\":\"planet id is not a valid id\"}", w.Body.String())
}

func TestShouldGetPlanetByIdReturnInternalServerErrorWhenProblemWithRepository(t *testing.T) {
	mongoMock := new(mock.MongoMock)
	swapiMock := new(mock.SwapiClientMock)
	mockLogger := new(mock.LoggerMock)

	mockLogger.On("LogWithFields", mock2.Anything, mock2.Anything, mock2.Anything, mock2.Anything).Return(nil)

	h := handler.NewPlanetHandler(mongoMock, swapiMock, mockLogger)

	id, _ := primitive.ObjectIDFromHex("5ea7208049e00ddb76994ede")

	mongoMock.On("FindById", id).Return(&repository.Planet{}, errors.New("error on repository"))

	r, _ := http.NewRequest("GET", "/v1/planets/5ea7208049e00ddb76994ede", nil)
	w := httptest.NewRecorder()

	vars := map[string]string{
		"planetId": "5ea7208049e00ddb76994ede",
	}

	r = mux.SetURLVars(r, vars)

	h.GetPlanetById(w, r)

	mockLogger.AssertNumberOfCalls(t, "LogWithFields", 1)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "{\"description\":\"error on repository\"}", w.Body.String())
}

func TestShouldReturnAllPlanetsWithoutFilter(t *testing.T) {
	mongoMock := new(mock.MongoMock)
	swapiMock := new(mock.SwapiClientMock)
	mockLogger := new(mock.LoggerMock)

	h := handler.NewPlanetHandler(mongoMock, swapiMock, mockLogger)

	returnedPlanets := []repository.Planet{{Name: "Aldebaran", Land: "Dry", Weather: "Dry", AppearanceQuantity: 2},
		{Name: "Tattoine", Land: "Dry", Weather: "Dry", AppearanceQuantity: 1}}

	mongoMock.On("FindAll", repository.Filter{}).Return(&returnedPlanets, nil)

	r, _ := http.NewRequest("GET", "/v1/planets", nil)
	w := httptest.NewRecorder()

	h.GetPlanets(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "[{\"Name\":\"Aldebaran\",\"Weather\":\"Dry\",\"Land\":\"Dry\",\"AppearanceQuantity\":2},"+
		"{\"Name\":\"Tattoine\",\"Weather\":\"Dry\",\"Land\":\"Dry\",\"AppearanceQuantity\":1}]", w.Body.String())
}

func TestShouldReturnAllPlanetsWithoutFilterEmptyList(t *testing.T) {
	mongoMock := new(mock.MongoMock)
	swapiMock := new(mock.SwapiClientMock)
	mockLogger := new(mock.LoggerMock)

	h := handler.NewPlanetHandler(mongoMock, swapiMock, mockLogger)

	returnedPlanets := make([]repository.Planet, 0)

	mongoMock.On("FindAll", repository.Filter{}).Return(&returnedPlanets, nil)

	r, _ := http.NewRequest("GET", "/v1/planets", nil)
	w := httptest.NewRecorder()

	h.GetPlanets(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "[]", w.Body.String())
}

func TestShouldReturnAllPlanetsWithFilter(t *testing.T) {
	mongoMock := new(mock.MongoMock)
	swapiMock := new(mock.SwapiClientMock)
	mockLogger := new(mock.LoggerMock)

	h := handler.NewPlanetHandler(mongoMock, swapiMock, mockLogger)

	returnedPlanets := []repository.Planet{{Name: "Aldebaran", Land: "Dry", Weather: "Dry", AppearanceQuantity: 2}}

	mongoMock.On("FindAll", repository.Filter{Name: "Aldebaran"}).Return(&returnedPlanets, nil)

	r, _ := http.NewRequest("GET", "/v1/planets?name=Aldebaran", nil)

	w := httptest.NewRecorder()

	h.GetPlanets(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "[{\"Name\":\"Aldebaran\",\"Weather\":\"Dry\",\"Land\":\"Dry\",\"AppearanceQuantity\":2}]", w.Body.String())
}

func TestShouldRemovePlanetByIdWithSuccess(t *testing.T) {
	mongoMock := new(mock.MongoMock)
	swapiMock := new(mock.SwapiClientMock)
	mockLogger := new(mock.LoggerMock)

	h := handler.NewPlanetHandler(mongoMock, swapiMock, mockLogger)

	id, _ := primitive.ObjectIDFromHex("5ea7208049e00ddb76994ede")

	vars := map[string]string{
		"planetId": "5ea7208049e00ddb76994ede",
	}

	mongoMock.On("Delete", id).Return(nil)

	r, _ := http.NewRequest("GET", "/v1/planets/5ea7208049e00ddb76994ede", nil)

	r = mux.SetURLVars(r, vars)

	w := httptest.NewRecorder()

	h.RemovePlanetById(w, r)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestShouldReturnBadRequestDeletePlanetInvalidId(t *testing.T) {
	mongoMock := new(mock.MongoMock)
	swapiMock := new(mock.SwapiClientMock)
	mockLogger := new(mock.LoggerMock)

	mockLogger.On("LogWithFields", mock2.Anything, mock2.Anything, mock2.Anything, mock2.Anything).Return(nil)

	h := handler.NewPlanetHandler(mongoMock, swapiMock, mockLogger)

	vars := map[string]string{
		"planetId": "123",
	}

	r, _ := http.NewRequest("GET", "/v1/123", nil)

	r = mux.SetURLVars(r, vars)

	w := httptest.NewRecorder()

	h.RemovePlanetById(w, r)

	mongoMock.AssertNumberOfCalls(t, "Delete", 0)
	mockLogger.AssertNumberOfCalls(t, "LogWithFields", 1)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "{\"description\":\"planet id is not a valid id\"}", w.Body.String())
}

func TestShouldReturnInternalServerErrorWhenThereIsErrorOnRepostioryWhenDeleting(t *testing.T) {
	mongoMock := new(mock.MongoMock)
	swapiMock := new(mock.SwapiClientMock)
	mockLogger := new(mock.LoggerMock)
	mockLogger.On("LogWithFields", mock2.Anything, mock2.Anything, mock2.Anything, mock2.Anything).Return(nil)

	h := handler.NewPlanetHandler(mongoMock, swapiMock, mockLogger)

	vars := map[string]string{
		"planetId": "5ea7208049e00ddb76994ede",
	}

	id, _ := primitive.ObjectIDFromHex("5ea7208049e00ddb76994ede")

	mongoMock.On("Delete", id).Return(errors.New("error on repository"))

	r, _ := http.NewRequest("GET", "/v1/123", nil)

	r = mux.SetURLVars(r, vars)

	w := httptest.NewRecorder()

	h.RemovePlanetById(w, r)

	mockLogger.AssertNumberOfCalls(t, "LogWithFields", 1)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "{\"description\":\"error on repository\"}", w.Body.String())
}

func TestShouldReturnCreatedWhenCreatePlanetWithSuccess(t *testing.T) {
	mongoMock := new(mock.MongoMock)
	swapiMock := new(mock.SwapiClientMock)
	mockLogger := new(mock.LoggerMock)

	h := handler.NewPlanetHandler(mongoMock, swapiMock, mockLogger)

	id, _ := primitive.ObjectIDFromHex("5ea7208049e00ddb76994ede")

	films := make([]string, 0)
	films = append(films, "film 1")
	films = append(films, "film 2")

	swapiResponse := client.SwapiPlanet{Results: []client.Results{{Films: films}}}
	savedPlanet := repository.Planet{Id: id}

	mongoMock.On("Save", mock2.Anything).Return(&savedPlanet, nil)
	swapiMock.On("GetPlanetByName", "Aldebaran").Return(&swapiResponse, nil)

	planetRequest := handler.PlanetRequest{Name: "Aldebaran", Land: "dessert", Weather: "rain"}

	reqBodyBytes := new(bytes.Buffer)

	_ = json.NewEncoder(reqBodyBytes).Encode(planetRequest)

	r, _ := http.NewRequest("POST", "/v1/planets", reqBodyBytes)

	w := httptest.NewRecorder()

	h.SavePlanet(w, r)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, "v1/planets/5ea7208049e00ddb76994ede", w.Header().Get("Location"))
}

func TestShouldReturnNotFoundWhenPlanetNotExist(t *testing.T) {
	mongoMock := new(mock.MongoMock)
	swapiMock := new(mock.SwapiClientMock)
	mockLogger := new(mock.LoggerMock)
	mockLogger.On("LogWithFields", mock2.Anything, mock2.Anything, mock2.Anything, mock2.Anything).Return(nil)

	h := handler.NewPlanetHandler(mongoMock, swapiMock, mockLogger)
	films := make([]string, 0)
	films = append(films, "film 1")
	films = append(films, "film 2")

	var swapi *client.SwapiPlanet

	swapiMock.On("GetPlanetByName", "Aldebaran").Return(swapi, nil)

	planetRequest := handler.PlanetRequest{Name: "Aldebaran", Land: "dessert", Weather: "rain"}

	reqBodyBytes := new(bytes.Buffer)

	_ = json.NewEncoder(reqBodyBytes).Encode(planetRequest)

	r, _ := http.NewRequest("POST", "/v1/planets", reqBodyBytes)

	w := httptest.NewRecorder()

	h.SavePlanet(w, r)

	mockLogger.AssertNumberOfCalls(t, "LogWithFields", 1)
	mongoMock.AssertNumberOfCalls(t, "Save", 0)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestShouldReturnServerErrorWhenThereIsAnErrorCallingClient(t *testing.T) {
	mongoMock := new(mock.MongoMock)
	swapiMock := new(mock.SwapiClientMock)
	mockLogger := new(mock.LoggerMock)
	mockLogger.On("LogWithFields", mock2.Anything, mock2.Anything, mock2.Anything, mock2.Anything).Return(nil)

	h := handler.NewPlanetHandler(mongoMock, swapiMock, mockLogger)

	var swapi *client.SwapiPlanet

	swapiMock.On("GetPlanetByName", "Aldebaran").Return(swapi, errors.New("error calling client"))

	planetRequest := handler.PlanetRequest{Name: "Aldebaran", Land: "dessert", Weather: "rain"}

	reqBodyBytes := new(bytes.Buffer)

	_ = json.NewEncoder(reqBodyBytes).Encode(planetRequest)

	r, _ := http.NewRequest("POST", "/v1/planets", reqBodyBytes)

	w := httptest.NewRecorder()

	h.SavePlanet(w, r)

	mongoMock.AssertNumberOfCalls(t, "Save", 0)
	mockLogger.AssertNumberOfCalls(t, "LogWithFields", 1)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "{\"description\":\"error calling client\"}", w.Body.String())
}

func TestShouldReturnServerErrorWhenThereIsAnErrorCallingRepository(t *testing.T) {
	mongoMock := new(mock.MongoMock)
	swapiMock := new(mock.SwapiClientMock)
	mockLogger := new(mock.LoggerMock)
	mockLogger.On("LogWithFields", mock2.Anything, mock2.Anything, mock2.Anything, mock2.Anything).Return(nil)

	h := handler.NewPlanetHandler(mongoMock, swapiMock, mockLogger)

	films := make([]string, 0)
	films = append(films, "film 1")
	films = append(films, "film 2")

	var emptyResponse *repository.Planet
	swapiResponse := client.SwapiPlanet{Results: []client.Results{{Films: films}}}

	swapiMock.On("GetPlanetByName", "Aldebaran").Return(&swapiResponse, nil)
	mongoMock.On("Save", mock2.Anything).Return(emptyResponse, errors.New("error on repository"))

	planetRequest := handler.PlanetRequest{Name: "Aldebaran", Land: "dessert", Weather: "rain"}

	reqBodyBytes := new(bytes.Buffer)

	_ = json.NewEncoder(reqBodyBytes).Encode(planetRequest)

	r, _ := http.NewRequest("POST", "/v1/planets", reqBodyBytes)

	w := httptest.NewRecorder()

	h.SavePlanet(w, r)

	mockLogger.AssertNumberOfCalls(t, "LogWithFields", 1)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "{\"description\":\"error on repository\"}", w.Body.String())
}
