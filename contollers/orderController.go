package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	database "github.com/Joojo7/go-getting-started/database"
	helpers "github.com/Joojo7/go-getting-started/helpers"
	models "github.com/Joojo7/go-getting-started/models"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

//get orderCollection
var orderCollection *mongo.Collection = database.OpenCollection(database.Client, "order")

//GetOrders is the api used to get a multiple orders
func GetOrders(response http.ResponseWriter, request *http.Request) {

	response.Header().Add("Content-Type", "application/json")

	result, err := orderCollection.Find(context.TODO(), bson.M{})
	fmt.Print(result)

	defer cancel()
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(err.Error()))
	}
	var allOrders []bson.M
	if err = result.All(ctx, &allOrders); err != nil {
		log.Fatal(err)
	}

	response.Header().Add("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)

	json.NewEncoder(response).Encode(allOrders)

	// response.Write(jsonBytes)
}

//GetOrder is the api used to tget a single order
func GetOrder(response http.ResponseWriter, request *http.Request) {

	response.Header().Add("Content-Type", "application/json")

	params := mux.Vars(request)

	// id, _ := primitive.ObjectIDFromHex(params["id"])

	var order models.Order

	err := orderCollection.FindOne(ctx, bson.M{"order_id": params["id"]}).Decode(&order)
	defer cancel()
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(err.Error()))
	}

	response.Header().Add("Content-Type", "application/json")

	json.NewEncoder(response).Encode(order)

	// response.Write(jsonBytes)
}

//UpdateOrder is used to update orders
func UpdateOrder(response http.ResponseWriter, request *http.Request) {

	// check for content type existence and check for json validity
	helpers.ContentTypeValidator(response, request)

	// call MaxRequestValidator to enforce a maximum read of 1MB .
	dec := helpers.MaxRequestValidator(response, request)

	var order models.Order
	err := dec.Decode(&order)
	helpers.PostPatchRequestValidator(response, request, err)

	params := mux.Vars(request)
	filter := bson.M{"order_id": params["id"]}

	var updateObj primitive.D

	if order.Table_id != nil {
		updateObj = append(updateObj, bson.E{"name", order.Table_id})
	}

	order.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{"updated_at", order.Updated_at})

	upsert := true
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	result, err := orderCollection.UpdateOne(
		ctx,
		filter,
		bson.D{
			{"$set", updateObj},
		},
		&opt,
	)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(err.Error()))
	}

	defer cancel()

	response.Header().Add("Content-Type", "application/json")
	json.NewEncoder(response).Encode(result)

}

// var validate *validator.Validate

//CreateOrder for creating orders
func CreateOrder(response http.ResponseWriter, request *http.Request) {

	//set response format to JSON
	response.Header().Add("Content-Type", "application/json")

	// check for content type existence and check for json validity
	helpers.ContentTypeValidator(response, request)

	// call MaxRequestValidator to enforce a maximum read of 1MB .
	dec := helpers.MaxRequestValidator(response, request)

	var order models.Order
	err1 := dec.Decode(&order)

	//validate body structure
	if !helpers.PostPatchRequestValidator(response, request, err1) {
		return
	}

	order.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	order.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	order.ID = primitive.NewObjectID()
	order.Order_id = order.ID.Hex()

	orderCollection.InsertOne(ctx, order)

	defer cancel()

	json.NewEncoder(response).Encode(order)
	defer cancel()
}

//OrderItemOrderCreator is for creating orders for the order items
func OrderItemOrderCreator(order models.Order) string {

	order.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	order.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	order.ID = primitive.NewObjectID()
	order.Order_id = order.ID.Hex()

	orderCollection.InsertOne(ctx, order)
	defer cancel()

	return order.Order_id

}
