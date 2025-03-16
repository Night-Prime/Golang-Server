package todoHandler

import (
	"net/http"

	"github.com/thedevsaddam/renderer"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/bson"
)

// Getting all Task

func GetTasks(rw http.ResponseWriter r *http.Request) {
	var taskListFromDB = []TaskModel{}

	filter := bson.D{}
	cursor, err := db.Collection(collectionName).Find(context.Background(), filter)

	if err != nil {
		log.Printf("failed to fetch Task records from the db: %v\n", err.Error())
		rnd.JSON(rw, http.StatusBadRequest, renderer.M{
			"message": "Could not fetch the Task collection",
			"error":   err.Error(),
		})
		return
	}

	taskList := []Task{}
	if err = cursor.All(context.Background(), &taskListFromDB); err != nil {
		checkError(err)
	}

	// conversion from model to JSON, append to array(taskList)

	for _, td := range taskListFromDB {
		taskList = append(taskList, Task{
			ID:			td.ID.Hex(),
			Title: 		td.Title,
			Completed:	td.Completed,
			CreatedAt:	td.CreatedAt
		})
	}

	rnd.JSON(rw, http.StatusOk, GetTaskResponse{
		Message: `All Tasks Received`,
		Data: taskList
	})
}

// Creating Task
func CreateTask(rw http.ResponseWriter r *http.Request) {
	var taskReq CreateTask

	if err := json.NewDecoder(r.Body).Decode(&taskRequestBody); err != nil {
		log.Printf("Failed to Decode the JSON data %v \n", err.Error())
		rnd.JSON(rw, &http.StatusBadRequest, renderer.M{
			"message": "Could'nt decode data"
		})

		return
	}

	if taskReq.title == "" {
		log.Println("No title provided")
		rnd.JSON(rw, http.StatusBadRequest, renderer.M{
			"message": "No title provided, Please Add a Title"
		})
	}

	// create model
	taskModel := TaskModel{
		ID:        primitive.NewObjectID(),
		Title:     todoReq.Title,
		Completed: false,
		CreatedAt: time.Now(),
	}

	// append to the DB
	data, err := db.Collection(collectionName).InsertOne(r.context(), taskModel)
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
		"task": data
	})
}

// Updating Task
func UpdateTask(rw http.ResponseWriter r *http.Rquest) {
	// grabs the ID from the url params
	id := strings.TrimSpace(chi.UrlParam(r, "id"))
	// converts the ID from Hex value into Mongo ID Value
	res, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("Parameter is not a valid Hex Value: %v \n", err.Error())
		rnd.JSON(rw, http.StatusBadRequest, render.M{
			"message": "Paramter is not Valid",
			"error": err.Error()
		})

		return
	}

	// Variable to handle the update request body
	var updateTaskReq UpdateTask

	if err := json.NewDecoder(r.Body).Decode(&updateTaskReq); err != nil {
		log.Printf("Failed to Decode the JSON data: %v \n", err.Error())
		rnd.JSON(rw, &http.StatusBadRequest, err.Error())
	}

	if updateTaskReq.Title == "" {
		rnd.JSON(rw, http.StatusBadRequest, renderer.M{
			"message": "Title can not be empty"
		})

		return
	}

	// update task in the DB
	filter := bson.M{"id": res}
	update := bson.M{"$set": bson.M{"title": updateTaskReq.Title, "completed": updateTaskReq.Completed}}

	data, err := db.Collection(collectionName).UpdateOne(r.Context(), filter, update)
	if err != nil {
		log.Printf("Failed to Update the Task in DB: %v \n", err.Error())
		rnd.JSON(rw, http.StatusInternalServerError, renderer.M{
			"message": "Failed to Update Task in the DB",
			"error": err.Error()
		})

		return
	}

	rnd.JSON(rw, http.StatusOk, renderer.M{
		"message": "Task Updated Sucessfully",
		"data": data.ModifiedCount
	})
}

func DeleteTask(rw http.ResponseWriter, r *http.Request) {
	id := chi.UrlParam(r, "id")
	res, err := primitive.ObjectIDFromHex(id)
	if err != nil{
		log.Printf("Parameter is not a valid Hex Value: %v \n", err.Error())
		rnd.JSON(rw, http.StatusBadRequest, render.M{
			"message": "Invalid ID",
			"error": err.Error()
		})

		return
	}

	filter := bson.M{"id": res}
	if data, err := db.Collection(collectionName).DeleteOne(r.context(), filter); err != nil{
		log.Printf("Failed to Delete the Task in DB: %v \n", err.Error())
		rnd.JSON(rw, http.StatusInternalServerError, renderer.M{
			"message": "Failed to Delete Task in the DB",
			"error": err.Error()
		})
	} else {
		rnd.JSON(rw, http.StatusOK, renderer.M{
			"message": "item deleted successfully",
			"data":    data,
		})
	}
}

