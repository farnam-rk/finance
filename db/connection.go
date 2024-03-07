package db

import (
	"company/finance/config"
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Client Client of MongoDB
var Client *mongo.Client

// DB DB Connection
var DB *mongo.Database

// Disconnect Disconnect MongoDB
func Disconnect() {
	err := Client.Disconnect(context.TODO())
	if err != nil {
		log.Println(err.Error())
	}

}

// DBconnection Connect to MongoDB
func DBconnection() {
	log.Println("Conecting to DB")
	dbData := fmt.Sprintf("%s://%s:%d", config.ConfigValue.DBDriver, config.ConfigValue.DBHost1, config.ConfigValue.DBPort1)
	// dbData := fmt.Sprintf("%s://%s:%s@%s:%d,%s:%d/?replicaSet=%s", config.ConfigValue.DBDriver, config.ConfigValue.DBUsername, config.ConfigValue.DBPassword, config.ConfigValue.DBHost1, config.ConfigValue.DBPort1, config.ConfigValue.DBHost2, config.ConfigValue.DBPort2, config.ConfigValue.Replica)
	clientOption := options.Client().ApplyURI(dbData)

	client, err := mongo.Connect(context.TODO(), clientOption)
	if err != nil {
		log.Println(err.Error())
	}
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Println(err.Error())
	}
	Client = client

	DB = client.Database(config.ConfigValue.DBName)
	log.Println("Connected to DB")
}
