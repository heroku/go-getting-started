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

type InvoiceViewFormat struct {
	Invoice_id       string
	Payment_method   string
	Order_id         string
	Payment_status   *string
	Payment_due      interface{}
	Table_number     interface{}
	Payment_due_date time.Time
	Order_details    interface{}
}

//get invoiceCollection
var invoiceCollection *mongo.Collection = database.OpenCollection(database.Client, "invoice")

//GetInvoices is the api used to get a multiple invoices
func GetInvoices(response http.ResponseWriter, request *http.Request) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	response.Header().Add("Content-Type", "application/json")

	result, err := invoiceCollection.Find(context.TODO(), bson.M{})
	fmt.Print(result)

	defer cancel()
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(err.Error()))
	}
	var allInvoices []bson.M
	if err = result.All(ctx, &allInvoices); err != nil {
		log.Fatal(err)
	}

	response.Header().Add("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)

	json.NewEncoder(response).Encode(allInvoices)

	// response.Write(jsonBytes)
}

//GetInvoice is the api used to tget a single invoice
func GetInvoice(response http.ResponseWriter, request *http.Request) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	response.Header().Add("Content-Type", "application/json")

	params := mux.Vars(request)

	// id, _ := primitive.ObjectIDFromHex(params["id"])

	var invoice models.Invoice

	err := invoiceCollection.FindOne(ctx, bson.M{"invoice_id": params["id"]}).Decode(&invoice)
	defer cancel()
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(err.Error()))
	}

	allOrderItems, err := ItemsByOrder(invoice.Order_id)

	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(err.Error()))
	}

	//set response and response headers
	response.Header().Add("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)

	var invoiceView InvoiceViewFormat

	invoiceView.Order_id = invoice.Order_id
	invoiceView.Payment_due_date = invoice.Payment_due_date

	invoiceView.Payment_method = "null"
	if invoice.Payment_method != nil {
		invoiceView.Payment_method = *invoice.Payment_method
	}

	invoiceView.Invoice_id = invoice.Invoice_id
	invoiceView.Payment_status = *&invoice.Payment_status
	invoiceView.Payment_due = allOrderItems[0]["payment_due"]
	invoiceView.Table_number = allOrderItems[0]["table_number"]
	invoiceView.Order_details = allOrderItems[0]["order_items"]

	json.NewEncoder(response).Encode(invoiceView)

	// json.NewEncoder(response).Encode(invoice)

	// response.Write(jsonBytes)
}

//UpdateInvoice is used to update invoices
func UpdateInvoice(response http.ResponseWriter, request *http.Request) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	// check for content type existence and check for json validity
	helpers.ContentTypeValidator(response, request)

	// call MaxRequestValidator to enforce a maximum read of 1MB .
	dec := helpers.MaxRequestValidator(response, request)

	var invoice models.Invoice
	err := dec.Decode(&invoice)
	helpers.PostPatchRequestValidator(response, request, err)

	params := mux.Vars(request)
	filter := bson.M{"invoice_id": params["id"]}

	var updateObj primitive.D

	if invoice.Payment_method != nil {
		updateObj = append(updateObj, bson.E{"payment_method", invoice.Payment_method})
	}

	if invoice.Payment_status != nil {
		updateObj = append(updateObj, bson.E{"payment_status", invoice.Payment_status})
	}

	invoice.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{"updated_at", invoice.Updated_at})

	upsert := true
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	status := "PENDING"
	if invoice.Payment_status == nil {
		invoice.Payment_status = &status
	}

	result, err := invoiceCollection.UpdateOne(
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

//CreateInvoice for creating invoices
func CreateInvoice(response http.ResponseWriter, request *http.Request) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	//set response format to JSON
	response.Header().Add("Content-Type", "application/json")

	// check for content type existence and check for json validity
	helpers.ContentTypeValidator(response, request)

	// call MaxRequestValidator to enforce a maximum read of 1MB .
	dec := helpers.MaxRequestValidator(response, request)

	var invoice models.Invoice
	err1 := dec.Decode(&invoice)

	//validate body structure
	if !helpers.PostPatchRequestValidator(response, request, err1) {
		return
	}

	var order models.Order

	err := orderCollection.FindOne(ctx, bson.M{"order_id": invoice.Order_id}).Decode(&order)
	defer cancel()
	if err != nil {
		msg := fmt.Sprintf("message: Order was not found")
		response.WriteHeader(http.StatusNotFound)
		response.Write([]byte(msg))
		return
	}

	status := "PENDING"
	if invoice.Payment_status == nil {
		invoice.Payment_status = &status
	}

	invoice.Payment_due_date, _ = time.Parse(time.RFC3339, time.Now().AddDate(0, 0, 1).Format(time.RFC3339))
	invoice.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	invoice.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	invoice.ID = primitive.NewObjectID()
	invoice.Invoice_id = invoice.ID.Hex()

	// var num = toFixed(*invoice.Payment_due, 2)
	// invoice.Payment_due = &num

	invoiceCollection.InsertOne(ctx, invoice)
	defer cancel()

	json.NewEncoder(response).Encode(invoice)
	defer cancel()

}
