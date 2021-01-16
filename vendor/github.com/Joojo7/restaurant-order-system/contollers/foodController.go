package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"time"

	database "github.com/Joojo7/restaurant-order-system/database"
	helpers "github.com/Joojo7/restaurant-order-system/helpers"
	models "github.com/Joojo7/restaurant-order-system/models"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//get foodCollection
var foodCollection *mongo.Collection = database.OpenCollection(database.Client, "food")

//GetFoods is the api used to get a multiple foods
func GetFoods(response http.ResponseWriter, request *http.Request) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	response.Header().Add("Content-Type", "application/json")

	result, err := foodCollection.Find(context.TODO(), bson.M{})
	fmt.Print(result)

	defer cancel()
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(err.Error()))
	}
	var allFoods []bson.M
	if err = result.All(ctx, &allFoods); err != nil {
		log.Fatal(err)
	}

	response.Header().Add("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)

	json.NewEncoder(response).Encode(allFoods)

	// response.Write(jsonBytes)
}

//GetFood is the api used to tget a single food
func GetFood(response http.ResponseWriter, request *http.Request) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	response.Header().Add("Content-Type", "application/json")

	params := mux.Vars(request)

	// id, _ := primitive.ObjectIDFromHex(params["id"])

	var food models.Food

	err := foodCollection.FindOne(ctx, bson.M{"food_id": params["id"]}).Decode(&food)
	defer cancel()
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(err.Error()))
	}

	response.Header().Add("Content-Type", "application/json")

	json.NewEncoder(response).Encode(food)

	// response.Write(jsonBytes)
}

//UpdateFood is used to update foods
func UpdateFood(response http.ResponseWriter, request *http.Request) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	// check for content type existence and check for json validity
	helpers.ContentTypeValidator(response, request)

	// call MaxRequestValidator to enforce a maximum read of 1MB .
	dec := helpers.MaxRequestValidator(response, request)

	var food models.Food
	err := dec.Decode(&food)
	helpers.PostPatchRequestValidator(response, request, err)

	params := mux.Vars(request)
	filter := bson.M{"food_id": params["id"]}

	var updateObj primitive.D

	if food.Name != nil {
		updateObj = append(updateObj, bson.E{"name", food.Name})

	}

	if food.Price != nil {
		updateObj = append(updateObj, bson.E{"price", food.Price})
	}

	food.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{"updated_at", food.Updated_at})

	upsert := true
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	result, err := foodCollection.UpdateOne(
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

//CreateFood for creating foods
func CreateFood(response http.ResponseWriter, request *http.Request) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	//set response format to JSON
	response.Header().Add("Content-Type", "application/json")

	// check for content type existence and check for json validity
	helpers.ContentTypeValidator(response, request)

	// call MaxRequestValidator to enforce a maximum read of 1MB .
	dec := helpers.MaxRequestValidator(response, request)

	var food models.Food
	err1 := dec.Decode(&food)

	//validate body structure
	if !helpers.PostPatchRequestValidator(response, request, err1) {
		return
	}

	var menu models.Menu

	err := menuCollection.FindOne(ctx, bson.M{"menu_id": food.Menu_id}).Decode(&menu)
	defer cancel()
	if err != nil {
		msg := fmt.Sprintf("message: Menu was not found")
		response.WriteHeader(http.StatusNotFound)
		response.Write([]byte(msg))
		return
	}

	food.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	food.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	food.ID = primitive.NewObjectID()
	food.Food_id = food.ID.Hex()
	var num = toFixed(*food.Price, 2)
	food.Price = &num

	foodCollection.InsertOne(ctx, food)
	defer cancel()

	json.NewEncoder(response).Encode(food)
	defer cancel()

}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}
