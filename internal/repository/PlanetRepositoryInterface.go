package repository

import "go.mongodb.org/mongo-driver/bson/primitive"

type Planet struct {
	Id                 primitive.ObjectID `json:"-" bson:"_id"`
	Name               string             `bson:"name"`
	Weather            string             `bson:"weather"`
	Land               string             `bson:"land"`
	AppearanceQuantity int                `bson:"appearanceQuantity"`
}

type PlanetRepositoryInterface interface {
	FindById(id primitive.ObjectID) (*Planet, error)
	Save(planet *Planet) (*Planet, error)
	FindAll(filter Filter) (*[]Planet, error)
	Delete(id primitive.ObjectID) error
}