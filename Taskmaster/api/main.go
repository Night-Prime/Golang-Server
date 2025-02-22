package main

import (
	"fmt"
	"net/http"
	"log"
	"os"
	"database/sql"

	"github.com/joho/godotenv"
	"github.com/go-chi/chi"
)

func main() {
	godotenv.Load(".env")

    port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT is not found")
	}
	port = ":" + port

    fmt.Println("--------------------------------------------- \n")
    fmt.Printf(" Starting Task Master Server on port%s\n", port)
	fmt.Println("--------------------------------------------- \n")

    err := http.ListenAndServe(port, nil)
    if err != nil {
        log.Fatalf("Server failed to start: %v", err)
    }
}