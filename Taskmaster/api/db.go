package main

import(
	"log"
	"os"
	"context"
	"time"

	"github.com/joho/godotenv"
	"github.com/thedevsaddam/renderer"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

)

	// // Setup variables
	// var Rnd *renderer.Render
	// var Client *mongo.Client
	// var DB *mongo.Database

	// const (
	// 	DBName			string = "taskmaster"
	// 	CollectionName 	string = "todo"
	// )

// init() - configure & setup
func Init() {
	godotenv.Load(".env")
	
	log.Println("Currently Initializing TaskMaster")
	log.Println("--------------------------------------------- \n")

	DB_URI := os.Getenv("DB_URI")
    if DB_URI == "" {
        log.Println("DB_URI environment variable is not set")
    }

	Rnd = renderer.New()
	var err error

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	Client, err = mongo.Connect(ctx, options.Client().ApplyURI(DB_URI))
	CheckError(err)

	err = Client.Ping(ctx, readpref.Primary())
	CheckError(err)

	DB = Client.Database(DBName)
}