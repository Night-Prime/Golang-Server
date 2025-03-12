package main

import (
	// "fmt"
	"net/http"
	"log"
	"os"
	"context"
	"time"
	"os/signal"

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
	checkError(err)

	err = client.Ping(ctx, readpref.Primary())
	checkError(err)

	db = client.Database(dbName)
}

func checkError(err error){
	if err != nil {
		log.Println("Error Ocurred --------------------------------------- \n")
		log.Fatal(err)
	}
}

func main() {
	godotenv.Load(".env")

    port := os.Getenv("PORT")
	if port == "" {
		log.Println("PORT is not found")
	}
	port = ":" + port

	app := &http.Server{
		Addr: 			port,
		Handler: 		chi.NewRouter(),
		ReadTimeout:	60 * time.Second,
		WriteTimeout:	60 * time.Second,
	}

	// channel to handle graceful shutdown (receive signal)
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt)

	// then set up the server in a different goroutine
	go func() {
		log.Println("--------------------------------------------- \n")
    	log.Printf(" Starting Task Master Server on port%s\n", port)
		log.Println("--------------------------------------------- \n")

   		if err := app.ListenAndServe(); err != nil {
	    	log.Fatal("TaskMaster failed to start: %v", err)
   		}
	} ()

	// listen for signals that could shutdown the server
	sig := <-stopChan
	log.Printf("Signal received %v \n", sig)

	// disconnects the mongo DB instance
	if err := client.Disconnect(context.Background()); err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := app.Shutdown(ctx); err != nil {
		log.Fatalf("TaskMasteer shutdown failed: %v \n", err)
	}
	log.Println("TaskMaster is shutting down now")
}