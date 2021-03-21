package mock

import (
	"github.com/bernardoms/StarWarsPlanetAPI-GO/internal/repository"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MongoMock struct {
	mock.Mock
}

func (m *MongoMock) FindById(id primitive.ObjectID) (*repository.Planet, error) {
	args := m.Called(id)
	return args.Get(0).(*repository.Planet), args.Error(1)
}

func (m *MongoMock) FindAll(filter repository.Filter) (*[]repository.Planet, error) {
	args := m.Called(filter)
	return args.Get(0).(*[]repository.Planet), args.Error(1)
}

func (m *MongoMock) Save(planet *repository.Planet) (*repository.Planet, error) {
	args := m.Called(planet)
	return args.Get(0).(*repository.Planet), args.Error(1)
}

func (m *MongoMock) Delete(id primitive.ObjectID) error {
	args := m.Called(id)
	return args.Error(0)
}
