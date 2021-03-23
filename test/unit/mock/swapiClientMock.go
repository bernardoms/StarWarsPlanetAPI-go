package mock

import (
	"github.com/bernardoms/StarWarsPlanetAPI-GO/internal/client"
	"github.com/stretchr/testify/mock"
)

type SwapiClientMock struct {
	mock.Mock
}

func (m *SwapiClientMock) GetPlanetByName(name string) (*client.SwapiPlanet, error)  {
	args := m.Called(name)
	return args.Get(0).(*client.SwapiPlanet), args.Error(1)
}