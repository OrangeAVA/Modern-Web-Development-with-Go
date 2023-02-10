package server

import (
	"context"
	"log"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitDatabase(config *viper.Viper) *mongo.Client {
	connectionString := config.GetString("database.connection_string")

	if connectionString == "" {
		log.Fatalf("Database connectin string is missing")
	}

	clientOptions := options.Client().ApplyURI(connectionString)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatalf("Error while initializing database: %v", err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		client.Disconnect(context.TODO())
		log.Fatalf("Error while validating database: %v", err)
	}

	return client
}
