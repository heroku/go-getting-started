package controller

import (
	"context"
	"encoding/json"
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

//get tableCollection
var tableCollection *mongo.Collection = database.OpenCollection(database.Client, "table")

//GetTables is the api used to get a multiple tables
func GetTables(response http.ResponseWriter, request *http.Request) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	response.Header().Add("Content-Type", "application/json")

	result, err := tableCollection.Find(context.TODO(), bson.M{})

	defer cancel()
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(err.Error()))
	}
	var allTables []bson.M
	if err = result.All(ctx, &allTables); err != nil {
		log.Fatal(err)
	}

	response.Header().Add("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)

	json.NewEncoder(response).Encode(allTables)

	// response.Write(jsonBytes)
}

//GetTable is the api used to tget a single table
func GetTable(response http.ResponseWriter, request *http.Request) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	response.Header().Add("Content-Type", "application/json")

	params := mux.Vars(request)

	// id, _ := primitive.ObjectIDFromHex(params["id"])

	var table models.Table

	err := tableCollection.FindOne(ctx, bson.M{"table_id": params["id"]}).Decode(&table)
	defer cancel()
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(err.Error()))
	}

	response.Header().Add("Content-Type", "application/json")

	json.NewEncoder(response).Encode(table)

	// response.Write(jsonBytes)
}

//UpdateTable is used to update tables
func UpdateTable(response http.ResponseWriter, request *http.Request) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	// check for content type existence and check for json validity
	helpers.ContentTypeValidator(response, request)

	// call MaxRequestValidator to enforce a maximum read of 1MB .
	dec := helpers.MaxRequestValidator(response, request)

	var table models.Table
	err := dec.Decode(&table)
	helpers.PostPatchRequestValidator(response, request, err)

	params := mux.Vars(request)
	filter := bson.M{"table_id": params["id"]}

	var updateObj primitive.D

	if table.Number_of_guests != nil {
		updateObj = append(updateObj, bson.E{"number_of_guests", table.Number_of_guests})

	}

	if table.Table_number != nil {
		updateObj = append(updateObj, bson.E{"table_number", table.Table_number})
	}

	if table.Order_id != nil {
		updateObj = append(updateObj, bson.E{"order_id", table.Order_id})
	}

	table.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	upsert := true
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	result, err := tableCollection.UpdateOne(
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

//CreateTable for creating tables
func CreateTable(response http.ResponseWriter, request *http.Request) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	//set response format to JSON
	response.Header().Add("Content-Type", "application/json")

	// check for content type existence and check for json validity
	helpers.ContentTypeValidator(response, request)

	// call MaxRequestValidator to enforce a maximum read of 1MB .
	dec := helpers.MaxRequestValidator(response, request)

	var table models.Table
	err1 := dec.Decode(&table)

	//validate body structure
	if !helpers.PostPatchRequestValidator(response, request, err1) {
		return
	}

	table.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	table.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	table.ID = primitive.NewObjectID()
	table.Table_id = table.ID.Hex()

	tableCollection.InsertOne(ctx, table)
	defer cancel()

	json.NewEncoder(response).Encode(table)
	defer cancel()

}
