package config

import "os"

type MongoConfig struct {
	MongoURI string
	Database string
}

func NewMongoConfig() *MongoConfig {
	return &MongoConfig{
		MongoURI: os.Getenv("MONGO_URI"),
		Database: os.Getenv("DATABASE"),
	}
}