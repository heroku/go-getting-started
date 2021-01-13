package controller

import (
	"context"
	"encoding/json"
	"log"
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

type OrderItemPack struct {
	Table_id    *string
	Order_items []models.OrderItem
}

//get orderItemCollection
var orderItemCollection *mongo.Collection = database.OpenCollection(database.Client, "orderItem")

//GetOrderItems is the api used to get a multiple orderItems
func GetOrderItems(response http.ResponseWriter, request *http.Request) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	response.Header().Add("Content-Type", "application/json")

	result, err := orderItemCollection.Find(context.TODO(), bson.M{})

	defer cancel()
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(err.Error()))
	}
	var allOrderItems []bson.M
	if err = result.All(ctx, &allOrderItems); err != nil {
		log.Fatal(err)
	}

	response.Header().Add("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)

	json.NewEncoder(response).Encode(allOrderItems)

	// response.Write(jsonBytes)
}

//GetOrderItems is the api used to get a multiple orderItems
func GetOrderItemsByOrder(response http.ResponseWriter, request *http.Request) {

	response.Header().Add("Content-Type", "application/json")

	params := mux.Vars(request)

	allOrderItems, err := ItemsByOrder(params["id"])

	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(err.Error()))
	}

	response.Header().Add("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)

	json.NewEncoder(response).Encode(allOrderItems)

	// response.Write(jsonBytes)
}

func ItemsByOrder(id string) (OrderItems []primitive.M, err error) {

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	matchStage := bson.D{{"$match", bson.D{{"order_id", id}}}}
	lookupStage := bson.D{{"$lookup", bson.D{{"from", "food"}, {"localField", "food_id"}, {"foreignField", "food_id"}, {"as", "food"}}}}
	unwindStage := bson.D{{"$unwind", bson.D{{"path", "$food"}, {"preserveNullAndEmptyArrays", true}}}}

	lookupOrderStage := bson.D{{"$lookup", bson.D{{"from", "order"}, {"localField", "order_id"}, {"foreignField", "order_id"}, {"as", "order"}}}}
	unwindOrderStage := bson.D{{"$unwind", bson.D{{"path", "$order"}, {"preserveNullAndEmptyArrays", true}}}}

	lookupTableStage := bson.D{{"$lookup", bson.D{{"from", "table"}, {"localField", "order.table_id"}, {"foreignField", "table_id"}, {"as", "table"}}}}
	unwindTableStage := bson.D{{"$unwind", bson.D{{"path", "$table"}, {"preserveNullAndEmptyArrays", true}}}}
	projectStage := bson.D{
		{"$project", bson.D{
			{"_id", 0},
			{"amount", "$food.price"},
			{"total_count", 1},
			{"food_name", "$food.name"},
			{"food_image", "$food.food_image"},
			{"table_number", "$table.table_number"},
			{"table_id", "$table.table_id"},
			{"order_id", "$order.order_id"},
			{"price", "$food.price"},
			{"quantity", 1},
		}}}
	groupStage := bson.D{{"$group", bson.D{{"_id", bson.D{{"order_id", "$order_id"}, {"table_id", "$table_id"}, {"table_number", "$table_number"}}}, {"payment_due", bson.D{{"$sum", "$amount"}}}, {"total_count", bson.D{{"$sum", 1}}}, {"order_items", bson.D{{"$push", "$$ROOT"}}}}}}

	projectStage2 := bson.D{
		{"$project", bson.D{
			{"_id", 0},
			{"payment_due", 1},
			{"total_count", 1},
			{"table_number", "$_id.table_number"},
			{"order_items", 1},
		}}}

	result, err := orderItemCollection.Aggregate(ctx, mongo.Pipeline{
		matchStage,
		lookupStage,
		unwindStage,
		lookupOrderStage,
		unwindOrderStage,
		lookupTableStage,
		unwindTableStage,
		projectStage,
		groupStage,
		projectStage2})

	if err != nil {
		panic(err)
	}

	if err = result.All(ctx, &OrderItems); err != nil {
		panic(err)
	}

	defer cancel()

	return OrderItems, err
}

//GetOrderItem is the api used to tget a single orderItem
func GetOrderItem(response http.ResponseWriter, request *http.Request) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	response.Header().Add("Content-Type", "application/json")

	params := mux.Vars(request)

	// id, _ := primitive.ObjectIDFromHex(params["id"])

	var orderItem models.OrderItem

	err := orderItemCollection.FindOne(ctx, bson.M{"orderItem_id": params["id"]}).Decode(&orderItem)
	defer cancel()
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(err.Error()))
	}

	response.Header().Add("Content-Type", "application/json")

	json.NewEncoder(response).Encode(orderItem)

	// response.Write(jsonBytes)
}

//UpdateOrderItem is used to update orderItems
func UpdateOrderItem(response http.ResponseWriter, request *http.Request) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	// check for content type existence and check for json validity
	helpers.ContentTypeValidator(response, request)

	// call MaxRequestValidator to enforce a maximum read of 1MB .
	dec := helpers.MaxRequestValidator(response, request)

	var orderItem models.OrderItem
	err := dec.Decode(&orderItem)
	helpers.PostPatchRequestValidator(response, request, err)

	params := mux.Vars(request)
	filter := bson.M{"order_item_id": params["id"]}

	var updateObj primitive.D

	if orderItem.Unit_price != nil {

		updateObj = append(updateObj, bson.E{"unit_price", *orderItem.Unit_price})

	}

	if orderItem.Quantity != nil {
		updateObj = append(updateObj, bson.E{"quantity", *orderItem.Quantity})
	}

	if orderItem.Food_id != nil {
		updateObj = append(updateObj, bson.E{"food_id", *orderItem.Food_id})
	}

	orderItem.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{"updated_at", orderItem.Updated_at})

	upsert := true
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	result, err := orderItemCollection.UpdateOne(
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

//CreateOrderItem for creating orderItems
func CreateOrderItem(response http.ResponseWriter, request *http.Request) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	var orderItemPack OrderItemPack
	var order models.Order

	order.Order_Date, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	orderItemsToBeInserted := []interface{}{}
	//set response format to JSON
	response.Header().Add("Content-Type", "application/json")

	// check for content type existence and check for json validity
	helpers.ContentTypeValidator(response, request)

	// call MaxRequestValidator to enforce a maximum read of 1MB .
	dec := helpers.MaxRequestValidator(response, request)

	// var orderItem models.OrderItem
	errOrderItemPack := dec.Decode(&orderItemPack)
	//validate body structure of the order item pack
	if !helpers.PostPatchRequestValidator(response, request, errOrderItemPack) {
		return
	}

	order.Table_id = orderItemPack.Table_id
	order_id := OrderItemOrderCreator(order)

	for _, orderItem := range orderItemPack.Order_items {
		orderItem.Order_id = order_id

		orderItem.ID = primitive.NewObjectID()
		orderItem.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		orderItem.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		orderItem.Order_item_id = orderItem.ID.Hex()
		var num = toFixed(*orderItem.Unit_price, 2)
		orderItem.Unit_price = &num

		orderItemsToBeInserted = append(orderItemsToBeInserted, orderItem)
	}

	//validate existence of request body

	//validate body structure

	insertedOrderItems, err := orderItemCollection.InsertMany(ctx, orderItemsToBeInserted)
	if err != nil {
		log.Fatal(err)
	}
	defer cancel()

	json.NewEncoder(response).Encode(insertedOrderItems.InsertedIDs)

	defer cancel()

}
