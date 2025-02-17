// First attempt to building a server in Golang
package main

import (
	"fmt"
	"io"
	"net/http"
	"log"
)

func Init(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "This is my first server")
}


func main() {
	portString := ":6060"
	 
	go func() {
		// introducing handleFunc since it's the one responsible for registering functions for http server
		http.HandleFunc("/", Init)

		fmt.Printf("Starting my first server, PORT %s ...", portString)

		if err := http.ListenAndServe(portString, nil); err != nil {
			log.Fatal(err)
		}
	}()

	// trying to print headers
	resp, err := http.Get("http://localhost"+ portString)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	bodyString := string(bodyBytes)
	fmt.Println(bodyString)

	fmt.Println("------------------------------ \n")
	fmt.Println("Response Headers \n")

	for k, v := range resp.Header {
		fmt.Printf("The %s: %v \n", k, v)
	}
}