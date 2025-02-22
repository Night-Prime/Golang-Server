package main

import (
	"fmt"
	"net/http"
	"log"
	"os"
	"context"
	"time"

	"github.com/joho/godotenv"
	"github.com/go-chi/chi/v5"
	// "github.com/go-chi/chi/v5/middleware"
	"github.com/thedevsaddam/renderer"
	// "go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

)


	// Setup variables
	var rnd *renderer.Render
	var client *mongo.Client
	var db *mongo.Database

	const (
		dbName		string = "taskmaster"
		collectionName string = "todo"
	)

// init() - configure & setup
func init() {
	godotenv.Load(".env")
	
	fmt.Println("Currrently Initializing Server")
	fmt.Println("--------------------------------------------- \n")

	DB_URI := os.Getenv("DB_URI")
    if DB_URI == "" {
        log.Fatal("DB_URI environment variable is not set")
    }

	rnd = renderer.New()
	var err error

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err = mongo.Connect(ctx, options.Client().ApplyURI(DB_URI))
	checkError(err)

	err = client.Ping(ctx, readpref.Primary())
	checkError(err)

	db = client.Database(dbName)
}

func checkError(err error){
	if err != nil {
		fmt.Println("Error Ocurred --------------------------------------- \n")
		log.Fatal(err)
	}
}

func main() {
	godotenv.Load(".env")

    port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT is not found")
	}
	port = ":" + port

	app := &http.Server{
		Addr: 			port,
		Handler: 		chi.NewRouter(),
		ReadTimeout:	60 * time.Second,
		WriteTimeout:	60 * time.Second,
	}

    fmt.Println("--------------------------------------------- \n")
    fmt.Printf(" Starting Task Master Server on port%s\n", port)
	fmt.Println("--------------------------------------------- \n")

   if err := app.ListenAndServe(); err != nil {
	        log.Fatalf("Server failed to start: %v", err)
   }
}