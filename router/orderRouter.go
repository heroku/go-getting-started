package router

import (
	OrderController "github.com/Joojo7/go-getting-started/contollers"
	"github.com/gorilla/mux"
)

//OrderRoutes function
func OrderRoutes(incomingRoutes *mux.Router) {

	// myRouter := mux.NewRouter().NewRoute().Subrouter().StrictSlash(true)

	incomingRoutes.HandleFunc("/orders", OrderController.GetOrders).Methods("GET")
	incomingRoutes.HandleFunc("/orders/{id}", OrderController.GetOrder).Methods("GET")
	incomingRoutes.HandleFunc("/orders/{id}", OrderController.UpdateOrder).Methods("PATCH")
	incomingRoutes.HandleFunc("/orders", OrderController.CreateOrder).Methods("POST")

}
