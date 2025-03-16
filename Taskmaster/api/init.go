package init

import(
	"fmt"
	"os"
	"log"
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

	// Setup variables
	var rnd *renderer.Render
	var client *mongo.Client
	var db *mongo.Database

	const (
		dbName			string = "taskmaster"
		collectionName 	string = "todo"
	)

// init() - configure & setup
func Init() {
	godotenv.Load(".env")
	
	log.Println("Currently Initializing TaskMaster")
	log.Println("--------------------------------------------- \n")

	DB_URI := os.Getenv("DB_URI")
    if DB_URI == "" {
        log.Println("DB_URI environment variable is not set")
    }

	rnd = renderer.New()
	var err error

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err = mongo.Connect(ctx, options.Client().ApplyURI(DB_URI))
	CheckError(err)

	err = client.Ping(ctx, readpref.Primary())
	CheckError(err)

	db = client.Database(dbName)
}