package models

import(
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	TaskModel struct {
		ID 			primitive.ObjectID 	`bson:"_id,omitempty"`
		Title 		string				`bson:"title"`
		Completed	bool				`bson:"completed"`
		CreatedAt 	time.Time			`bson:"created_at"`
	}

	Task struct {
		ID 			string 				`json:"id"`
		Title 		string				`json:"title"`
		Completed	bool				`json:"completed"`
		CreatedAt 	time.Time			`json:"created_at"`
	}

	GetTaskResponse struct {
		Message	string 	`json:"message"`
		Data	[]Task	`json:"data"`
	}

	CreateTask struct {
		Title	string	`json:"title"`
	}

	UpdateTask struct {
		Title	string	`json:"title"`
		Completed bool	`json:"completed"`
	}
)