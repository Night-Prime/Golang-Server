package shared

import(
	"log"

	"github.com/thedevsaddam/renderer"
	"go.mongodb.org/mongo-driver/mongo"

)

	// Setup variables
	var Rnd *renderer.Render
	var Client *mongo.Client
	var DB *mongo.Database

	const (
		DBName			string = "taskmaster"
		CollectionName 	string = "todo"
	)


	func CheckError(err error){
	if err != nil {
		log.Println("Error Ocurred --------------------------------------- \n")
		log.Fatal(err)
	}
}