package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Joojo7/restaurant-order-system/router"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var MyEnv map[string]string

func main() {
	myEnv, err1 := godotenv.Read()
	if err1 != nil {
		log.Fatal(err1)
	}
	port := myEnv["PORT"]

	myRouter := mux.NewRouter().StrictSlash(true)

	// ROuter files
	router.Routes(myRouter)
	router.FoodRoutes(myRouter)
	router.OrderItemRoutes(myRouter)
	router.TableRoutes(myRouter)
	router.InvoiceRoutes(myRouter)
	router.OrderRoutes(myRouter)

	fmt.Printf("listening on %v \n", port)
	error1 := http.ListenAndServe(port, myRouter)
	if error1 != nil {
		panic(error1)
	}
}
