package main

import (
	"golang-fifa-world-cup-web-service/data"
	"golang-fifa-world-cup-web-service/handlers"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/subosito/gotenv"
)

func main() {
	gotenv.Load()
	data.Init()
	data.PrintUsage()
	_, err1 := godotenv.Read()
	if err1 != nil {
		log.Fatal(err1)
	}

	port := os.Getenv("PORT")

	http.HandleFunc("/", handlers.RootHandler)
	http.HandleFunc("/winners", handlers.WinnersHandler)
	http.ListenAndServe(port, nil)
}
