package router

import (
	OrderItemController "github.com/Joojo7/restaurant-order-system/contollers"
	"github.com/gorilla/mux"
)

//OrderItemRoutes function
func OrderItemRoutes(incomingRoutes *mux.Router) {

	// myRouter := mux.NewRouter().NewRoute().Subrouter().StrictSlash(true)

	incomingRoutes.HandleFunc("/orderItems", OrderItemController.GetOrderItems).Methods("GET")
	incomingRoutes.HandleFunc("/orderItems/{id}", OrderItemController.GetOrderItem).Methods("GET")
	incomingRoutes.HandleFunc("/orderItems-order/{id}", OrderItemController.GetOrderItemsByOrder).Methods("GET")
	incomingRoutes.HandleFunc("/orderItems/{id}", OrderItemController.UpdateOrderItem).Methods("PATCH")
	incomingRoutes.HandleFunc("/orderItems", OrderItemController.CreateOrderItem).Methods("POST")

}
