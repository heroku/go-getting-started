package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Invoice is the model that governs all notes objects retrived or inserted into the DB
type Invoice struct {
	ID               primitive.ObjectID `bson:"_id"`
	Invoice_id       string             `json:"invoice_id"`
	Order_id         string             `json:"order_id" `
	Payment_method   *string            `json:"payment_method" validate:"eq=CARD|eq=CASH"`
	Payment_status   *string            `json:"payment_status" validate:"required,eq=PENDING|eq=PAID"`
	Payment_due_date time.Time          `json:"Payment_due_date"`
	Created_at       time.Time          `json:"created_at"`
	Updated_at       time.Time          `json:"updated_at"`
}
