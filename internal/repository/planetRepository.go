package repository

import (
	"context"
	"github.com/bernardoms/StarWarsPlanetAPI-GO/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type sessionCreator struct {
	Client *mongo.Client
}

type Mongo struct {
	collection *mongo.Collection
	session *mongo.Client
}

type Filter struct {
	Name     string `schema:"name"`
}

func NewSession(config config.MongoConfig) *Mongo {

	mo := new(Mongo)

	client, err := mongo.NewClient(options.Client().ApplyURI(config.MongoURI))

	m := sessionCreator{
		Client: client,
	}

	ctx := context.TODO()

	if err != nil {
		log.Print("Error on creating mongo client ", err)
	}

	err = m.Client.Connect(ctx)

	mo.session = client

	if err != nil {
		log.Print("Error on connecting to database ", err)
	}

	mo.getCollection(config)

	return mo
}

func (m *Mongo) getCollection(config config.MongoConfig) {
	c := m.session.Database(config.Database).Collection(config.Database)
	m.collection = c
}

func (m *Mongo) Save(planet *Planet) (*Planet, error) {
	_, err := m.collection.InsertOne(context.TODO(), &planet)

	return planet, err
}

func (m *Mongo) FindAll(filter Filter) (*[]Planet, error) {
	planet := make([]Planet, 0)

	result, err := m.collection.Find(context.TODO(), mountFilter(filter))

	if err == nil && result != nil {
		err = result.All(context.TODO(), &planet)
	}

	return &planet, err
}

func (m *Mongo) FindById(id primitive.ObjectID) (*Planet, error) {
	var result *Planet

	cur, err := m.collection.Find(context.TODO(), bson.M{"_id": id})

	if err != nil {
		return &Planet{}, err
	}

	if cur.Next(context.TODO()) != false {
		err = cur.Decode(&result)
	}

	return result, err
}

func (m *Mongo) Delete(id primitive.ObjectID) error {
	filter := bson.M{"_id": id}

	_, err := m.collection.DeleteOne(context.TODO(), filter)

	return err
}

func mountFilter(filter Filter) bson.M{
	f := bson.M{}

	if filter.Name != "" {
		f["name"] = filter.Name
	}

	return f
}
