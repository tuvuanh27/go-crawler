package mongodriver

import (
	"context"
	"github.com/tuvuanh27/go-crawler/internal/pkg/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoConfig struct {
	URI          string `mapstructure:"uri"`
	DatabaseName string `mapstructure:"databaseName"`
}

var ProductCollection = "products"
var UserCollection = "users"

type Mongo struct {
	Client *mongo.Client
	config *MongoConfig
}

func setupUniqueIndex(db *mongo.Database, collectionName string, indexKeys bson.D, indexName string, logger logger.ILogger) error {
	collection := db.Collection(collectionName)

	// Define the index model
	indexModel := mongo.IndexModel{
		Keys:    indexKeys, // Use the passed index keys
		Options: options.Index().SetUnique(true).SetName(indexName),
	}

	// Create the index
	createdIndexName, err := collection.Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		logger.Error("Failed to create index:", err)
		return err
	}

	logger.Debug("Index created with name:", createdIndexName)
	return nil
}
func NewMongo(config *MongoConfig, logger logger.ILogger, ctx context.Context) *mongo.Database {
	var client *mongo.Client

	// Set client options
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().ApplyURI(config.URI).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		panic(err)
	}

	go func() {
		select {
		case <-ctx.Done():
			err := client.Disconnect(ctx)
			if err != nil {
				panic(err)
			}

			logger.Debug("Disconnected from database")
		}
	}()

	err = client.Ping(context.Background(), nil)
	if err != nil {
		panic(err)
	}

	db := client.Database(config.DatabaseName)

	// Set up unique index for product collection
	err = setupUniqueIndex(db, ProductCollection, bson.D{{"product_id", 1}}, "product_id", logger)
	if err != nil {
		panic(err)
	}

	return client.Database(config.DatabaseName)
}
