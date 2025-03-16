package main

import (
	"log"
)


func CheckError(err error){
	if err != nil {
		log.Println("Error Ocurred --------------------------------------- \n")
		log.Fatal(err)
	}
}