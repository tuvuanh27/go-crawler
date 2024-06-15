package mongodriver

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoConfig struct {
	URI string `mapstructure:"uri"`
}

type Mongo struct {
	Client *mongo.Client
	config *MongoConfig
}

func NewMongo(config *MongoConfig) (*mongo.Client, error) {
	var client *mongo.Client
	var err error

	// Set client options
	clientOptions := options.Client().ApplyURI(config.URI)

	// Connect to MongoDB
	client, err = mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}

	// Check the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (m *Mongo) Close() error {
	return m.Client.Disconnect(context.Background())
}
