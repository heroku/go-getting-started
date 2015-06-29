package main

import (
	"log"
	"net/http"
	"os"

	"github.com/heroku/go-getting-started/Godeps/_workspace/src/github.com/codegangsta/negroni"
	"github.com/heroku/go-getting-started/Godeps/_workspace/src/github.com/unrolled/render"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello From Go!"))
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	r := render.New()
	mux := http.NewServeMux()

	mux.HandleFunc("/html", helloHandler)

	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		// Assumes you have templates in ./templates
		r.HTML(w, http.StatusOK, "index", nil)
	})

	n := negroni.Classic()
	n.UseHandler(mux)
	n.Run(":" + port)
}
