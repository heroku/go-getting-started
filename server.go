package main

import (
	"golang-fifa-world-cup-web-service/data"
	"golang-fifa-world-cup-web-service/handlers"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {
	data.PrintUsage()
	myEnv, err1 := godotenv.Read()
	if err1 != nil {
		log.Fatal(err1)
	}
	port := myEnv["PORT"]

	http.HandleFunc("/", handlers.RootHandler)
	http.HandleFunc("/winners", handlers.WinnersHandler)
	http.ListenAndServe(port, nil)
}
