package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Menu is the model that governs all notes objects retrived or inserted into the DB
type Menu struct {
	ID         primitive.ObjectID `bson:"_id"`
	Name       string             `json:"name" validate:"required"`
	Category   string             `json:"category" validate:"required"`
	Start_Date *time.Time         `json:"start_date" validate:"required"`
	End_Date   *time.Time         `json:"end_date" validate:"required"`
	Created_at time.Time          `json:"created_at"`
	Updated_at time.Time          `json:"updated_at"`
	Menu_id    string             `json:"food_id"`
}
