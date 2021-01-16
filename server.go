package main

import (
	"golang-fifa-world-cup-web-service/handlers"
	"net/http"
	"os"

	"github.com/subosito/gotenv"
)

func init() {
	gotenv.Load()
}

func main() {

	port := os.Getenv("PORT")

	http.HandleFunc("/", handlers.RootHandler)
	http.HandleFunc("/winners", handlers.WinnersHandler)
	http.ListenAndServe(port, nil)
}
