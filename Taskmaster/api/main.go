package main

import (
	// "fmt"
	"net/http"
	"log"
	"os"
	"context"
	"time"
	"os/signal"
	"strings"
	"encoding/json"

	"github.com/joho/godotenv"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/thedevsaddam/renderer"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/Night-Prime/Golang-Server.git/taskmaster/api/models"

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

func homeHandler(rw http.ResponseWriter, r *http.Request) {
	filePath := "./ReadMe.md"
	err := rnd.FileView(rw, http.StatusOK, filePath, "readme.md")
	checkError(err)
}


func taskHandler() http.Handler {
	router := chi.NewRouter()
	router.Group(func(r chi.Router) {
		r.Get("/", getTasks)
		r.Post("/", createTask)
		r.Put("/{id}", updateTask)
		r.Delete("/{id}", deleteTask)
	})

	return router
}

func getTasks(rw http.ResponseWriter, r *http.Request) {
	var taskListFromDB = []models.TaskModel{}

	filter := bson.D{}
	cursor, err := db.Collection(collectionName).Find(context.Background(), filter)

	if err != nil {
		log.Printf("Failed to fetch Task records from the db: %v\n", err.Error())
		rnd.JSON(rw, http.StatusBadRequest, renderer.M{
			"message": "Could not fetch the Task collection",
			"error":   err.Error(),
		})
		return
	}

	taskList := []models.Task{}
	if err = cursor.All(context.Background(), &taskListFromDB); err != nil {
		checkError(err)
	}

	// conversion from model to JSON, append to array(taskList)
	for _, td := range taskListFromDB {
		taskList = append(taskList, models.Task{
			ID:			td.ID.Hex(),
			Title: 		td.Title,
			Completed:	td.Completed,
			CreatedAt:	td.CreatedAt,
		})
	}

	rnd.JSON(rw, http.StatusOK, models.GetTaskResponse{
		Message: `All Tasks Received`,
		Data: taskList,
	})
}

func createTask(rw http.ResponseWriter, r *http.Request) {
	var taskReq models.CreateTask

	if err := json.NewDecoder(r.Body).Decode(&taskReq); err != nil {
		log.Printf("Failed to Decode the JSON data: %v \n", err.Error())
		rnd.JSON(rw, http.StatusBadRequest, renderer.M{
			"message": "Could'nt decode data",
		})

		return
	}

	if taskReq.Title == "" {
		log.Println("No title provided")
		rnd.JSON(rw, http.StatusBadRequest, renderer.M{
			"message": "No title provided, Please Add a Title",
		})
	}

	// create model
	taskModel := models.TaskModel{
		ID:        primitive.NewObjectID(),
		Title:     taskReq.Title,
		Completed: false,
		CreatedAt: time.Now(),
	}

	// append to the DB
	data, err := db.Collection(collectionName).InsertOne(r.Context(), taskModel)
	if err != nil{
		log.Printf("Failed to Insert Task into the DB %v \n", err.Error())
		rnd.JSON(rw, http.StatusInternalServerError, renderer.M{
			"message": "Failed to Add Task to the DB",
			"error": err.Error(),
		})

		return
	}

	rnd.JSON(rw, http.StatusCreated, renderer.M{
		"message": "Task Sucessfully Created",
		"task": data,
	})
}

func updateTask(rw http.ResponseWriter, r *http.Request) {
	// grabs the ID from the url params
	id := strings.TrimSpace(chi.URLParam(r, "id"))
	log.Printf("The ID: %v", id)
	// converts the ID from Hex value into Mongo ID Value
	res, err := primitive.ObjectIDFromHex(id)
	log.Printf("The result is : %v", res)
	if err != nil {
		log.Printf("Parameter is not a valid Hex Value: %v \n", err.Error())
		rnd.JSON(rw, http.StatusBadRequest, renderer.M{
			"message": "Paramter is not Valid",
			"error": err.Error(),
		})

		return
	}

	// Variable to handle the update request body
	var updateTaskReq models.UpdateTask

	if err := json.NewDecoder(r.Body).Decode(&updateTaskReq); err != nil {
		log.Printf("Failed to Decode the JSON data: %v \n", err.Error())
		rnd.JSON(rw, http.StatusBadRequest, err.Error())
	}

	log.Println("The Data coming in : %v", updateTaskReq)

	if updateTaskReq.Title == "" {
		rnd.JSON(rw, http.StatusBadRequest, renderer.M{
			"message": "Title can not be empty",
		})

		return
	}

	// update task in the DB
	filter := bson.M{"_id": res}
	update := bson.M{"$set": bson.M{"title": updateTaskReq.Title, "completed": updateTaskReq.Completed}}

	data, err := db.Collection(collectionName).UpdateOne(r.Context(), filter, update)
	log.Println("The updated data: %v", data)
	if err != nil {
		log.Printf("Failed to Update the Task in DB: %v \n", err.Error())
		rnd.JSON(rw, http.StatusInternalServerError, renderer.M{
			"message": "Failed to Update Task in the DB",
			"error": err.Error(),
		})

		return
	}

	rnd.JSON(rw, http.StatusOK, renderer.M{
		"message": "Task Updated Sucessfully",
		"data": data.ModifiedCount,
	})
}

func deleteTask(rw http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	res, err := primitive.ObjectIDFromHex(id)
	if err != nil{
		log.Printf("Parameter is not a valid Hex Value: %v \n", err.Error())
		rnd.JSON(rw, http.StatusBadRequest, renderer.M{
			"message": "Invalid ID",
			"error": err.Error(),
		})

		return
	}

	filter := bson.M{"_id": res}
	if data, err := db.Collection(collectionName).DeleteOne(r.Context(), filter); err != nil{
		log.Printf("Failed to Delete the Task in DB: %v \n", err.Error())
		rnd.JSON(rw, http.StatusInternalServerError, renderer.M{
			"message": "Failed to Delete Task in the DB",
			"error": err.Error(),
		})
	} else {
		rnd.JSON(rw, http.StatusOK, renderer.M{
			"message": "item deleted successfully",
			"data":    data,
		})
	}
}

func main() {
	godotenv.Load(".env")

    port := os.Getenv("PORT")
	if port == "" {
		log.Println("PORT is not found")
	}
	port = ":" + port

	// set up routing
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Get("/", homeHandler)
	router.Mount("/task", taskHandler())

	app := &http.Server{
		Addr: 			port,
		Handler: 		router,
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