package mongo

//This module is a wrapper around a mongoDb connection
//it allows the connection to be reused in different goroutines
import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DBT struct {
    DB *mongo.Collection
}

var QuizDB DBT
//This function opens a connection to the remote mongoDB hosted on MongoAtlas
func init() {
	godotenv.Load("./.env")
	// Set client options
	clientOptions := options.Client().ApplyURI(os.Getenv("MONGO_URI"))

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}
	QuizDB.DB = client.Database("questions").Collection("question")
}