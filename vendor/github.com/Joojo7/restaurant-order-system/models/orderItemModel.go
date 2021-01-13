package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Food is the model that governs all notes objects retrived or inserted into the DB
type OrderItem struct {
	ID            primitive.ObjectID `bson:"_id"`
	Quantity      *string            `json:"quantity" validate:"required,eq=S|eq=M|eq=L"`
	Unit_price    *float64           `json:"unit_price" validate:"required"`
	Created_at    time.Time          `json:"created_at"`
	Updated_at    time.Time          `json:"updated_at"`
	Food_id       *string            `json:"food_id" validate:"required"`
	Order_item_id string             `json:"order_item_id"`
	Order_id      string             `json:"order_id" validate:"required"`
}
