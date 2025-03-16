package common

	// Setup variables
	var Rnd *renderer.Render
	var Client *mongo.Client
	var DB *mongo.Database

	const (
		DBName			string = "taskmaster"
		CollectionName 	string = "todo"
	)